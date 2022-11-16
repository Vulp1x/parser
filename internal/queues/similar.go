package queues

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/domain"
	"github.com/inst-api/parser/internal/instagrapi"
	"github.com/inst-api/parser/pkg/logger"
	"go.uber.org/multierr"
)

// ErrNoReadyBots нет ботов доступных для работы
var ErrNoReadyBots = errors.New("все боты заняты, попробуйте позже")

func (s Service) AddFindSimilarTask(
	ctx context.Context,
	datasetID uuid.UUID,
	initialBloggers []dbmodel.Blogger,
	botsPerDataset int,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf(ctx, "recovered panic: %s, stack: %s", r, string(debug.Stack()))
			err = fmt.Errorf("recovered panic: %s", r)
		}
	}()

	err = s.similarQueue.Add(s.findSimilarTask.WithArgs(ctx, datasetID, initialBloggers, botsPerDataset))
	if err != nil {
		return err
	}

	return nil

}

func (s Service) processFailedTask(ctx context.Context, datasetID uuid.UUID, _ []dbmodel.Blogger, _ int) error {
	logger.Info(ctx, "findSimilarTask failed, changing dataset status to draft")

	q := dbmodel.New(s.dbf(ctx))

	return q.UpdateDatasetStatus(ctx, dbmodel.UpdateDatasetStatusParams{Status: dbmodel.DraftDatasetStatus, ID: datasetID})
}

func (s Service) findSimilarBloggers(ctx context.Context, datasetID uuid.UUID, initialBloggers []dbmodel.Blogger, botsPerDataset int) error {
	startedAt := time.Now()
	q := dbmodel.New(s.dbf(ctx))

	bots, err := q.LockAvailableBots(ctx, int32(botsPerDataset))
	if err != nil {
		return fmt.Errorf("failed to find available bots: %v", err)
	}

	if len(bots) == 0 {
		return ErrNoReadyBots
	}

	if len(bots) < botsPerDataset {
		logger.Warnf(ctx, "got %d available bots, expected %d", len(bots), botsPerDataset)
	}

	if len(bots) > len(initialBloggers) {
		var botsUnlocked int
		for _, bot := range bots[len(initialBloggers):] {
			err = q.UnlockBot(ctx, bot.ID)
			if err != nil {
				logger.Errorf(ctx, "failed to unblock bot %s: %v", bot.Username, err)
			}

			botsUnlocked++
		}

		logger.Infof(ctx, "got %d initial bloggers and %d bots, unlocked %d bots", len(initialBloggers), len(bots), botsUnlocked)

		bots = bots[:len(initialBloggers)]
	}

	bloggersLen := len(initialBloggers)
	bloggersPerBot := bloggersLen / len(bots)
	allBloggersProcessed := false

	wg := &sync.WaitGroup{}

	errc := make(chan error, 100)

	for i := range bots {
		wg.Add(1)
		rightBorderOfBatch := (i + 1) * bloggersPerBot
		processorCtx := logger.WithFields(ctx, logger.Fields{"processor_index": i})

		var bloggersBatch []dbmodel.Blogger

		if rightBorderOfBatch >= bloggersLen {
			bloggersBatch = initialBloggers[i*bloggersPerBot:]
			allBloggersProcessed = true
		} else {
			bloggersBatch = initialBloggers[i*bloggersPerBot : rightBorderOfBatch]
		}

		go s.findAndSaveSimilarBloggers(processorCtx, datasetID, bloggersBatch, bots[i], wg, errc)

		if allBloggersProcessed && i != len(bots)-1 {
			logger.Warnf(ctx, "all bloggers processed with %d/%d bots", i+1, len(bots))
			break
		}
	}

	logger.Infof(ctx, "started %d bots for %d bloggers, waiting for them", len(bots), len(initialBloggers))
	wg.Wait()

	close(errc)

	for err2 := range errc {
		err = multierr.Append(err, err2)
	}
	if err != nil {
		return fmt.Errorf("got err from bots: %v", err)
	}

	err = q.UpdateDatasetStatus(ctx, dbmodel.UpdateDatasetStatusParams{Status: dbmodel.ReadyForParsingDatasetStatus, ID: datasetID})
	if err != nil {
		return fmt.Errorf("failed to update dataset status to ReadyForParsingDatasetStatus (%d): %v", dbmodel.ReadyForParsingDatasetStatus, err)
	}

	logger.Infof(ctx, "all goroutines completed in %s", time.Since(startedAt))

	return err
}

func (s Service) findAndSaveSimilarBloggers(ctx context.Context, datasetID uuid.UUID, initialBloggers []dbmodel.Blogger, bot dbmodel.Bot, wg *sync.WaitGroup, errc chan error) {
	defer wg.Done()

	startedAt := time.Now()
	ctx = logger.WithFields(ctx, logger.Fields{"bot_username": bot.Username})

	q := dbmodel.New(s.dbf(ctx))
	err := s.cli.CheckBot(ctx, bot.SessionID)
	if err != nil {
		nonBlockingWriteError(ctx, errc, fmt.Errorf("failed to check bot '%s': %v", bot.SessionID, err))

		if !errors.Is(err, instagrapi.ErrBotIsBlocked) {
			return
		}

		err2 := q.BlockBot(ctx, bot.ID)
		if err2 != nil {
			nonBlockingWriteError(ctx, errc, fmt.Errorf("failed to block bot (%s): %v", bot.ID, err))
		}

		return
	}

	defer func() {
		err = q.UnlockBot(ctx, bot.ID)
		if err != nil {
			nonBlockingWriteError(ctx, errc, fmt.Errorf("failed to unlock bot after processing (%s): %v", bot.ID, err))
			return
		}

		logger.Info(ctx, "unlocked bot after processing")
	}()

	var users domain.ShortInstUsers

	var count, totalCount int64

	for i, blogger := range initialBloggers {
		users, err = s.cli.FindSimilarBloggersShort(ctx, bot.SessionID, blogger.Username)
		if err != nil {
			logger.Errorf(ctx, "failed to find similar accounts: %v", err)

			if errors.Is(err, instagrapi.ErrBotIsBlocked) {
				err2 := q.BlockBot(ctx, bot.ID)
				if err2 != nil {
					nonBlockingWriteError(ctx, errc, fmt.Errorf("failed to block bot (%s): %v", bot.ID, err))
				}
				continue
			}

			if errors.Is(err, instagrapi.ErrBloggerNotFound) {
				err = q.SetBloggerIsParsed(ctx, dbmodel.SetBloggerIsParsedParams{IsCorrect: false, ID: blogger.ID})
				if err != nil {
					nonBlockingWriteError(ctx, errc, fmt.Errorf("failed to set bot  (%s): %v", bot.ID, err))
				}
			}

			continue
		}

		if len(users) == 0 {
			logger.Warnf(ctx, "found 0 bloggers for blogger %s", blogger.Username)
			err = q.SetBloggerIsParsed(ctx, dbmodel.SetBloggerIsParsedParams{IsCorrect: false, ID: blogger.ID})
			if err != nil {
				nonBlockingWriteError(ctx, errc, fmt.Errorf("failed to set bot  (%s): %v", bot.ID, err))
			}
			continue
		}

		if users[0].Username == blogger.Username {
			err = q.UpdateBlogger(ctx, domain.InstUserShort(users[0]).ToUpdateParams(blogger.ID, true))
			if err != nil {
				nonBlockingWriteError(ctx, errc, fmt.Errorf("failed to update initial bot (%s): %v", bot.ID, err))
			}

			if len(users) == 1 {
				users = []domain.InstUserShort{}
			} else {
				users = users[1:]
			}
		}

		count, err = q.SaveBloggers(ctx, users.ToSaveBloggersParmas(datasetID))
		if err != nil {
			nonBlockingWriteError(ctx, errc, fmt.Errorf("failed to save parsed bloggers from blogger  (%s): %v", bot.ID, err))
			continue
		}

		totalCount += count

		logger.Infof(ctx, "saved %d bloggers from initial blogger '%s' (parsed %d/%d)", count, blogger.Username, i+1, len(initialBloggers))

		err = q.MarkBloggerAsSimilarAccountsFound(ctx, blogger.ID)
		if err != nil {
			nonBlockingWriteError(ctx, errc, fmt.Errorf("failed to mark blogger (%s) as ready for target's parsing : %v", blogger.ID, err))
		}
	}

	logger.Infof(ctx, "saved %d new bloggers in %s", totalCount, time.Since(startedAt))
}

func nonBlockingWriteError(ctx context.Context, errc chan error, err error) {
	logger.Error(ctx, err)

	select {
	case errc <- err:
	default:
		logger.Errorf(ctx, "failed to write error (%v) to err chan", err)
	}
}
