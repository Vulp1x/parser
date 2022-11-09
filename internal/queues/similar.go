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
	"github.com/vmihailenco/taskq/v3"
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

	task := taskq.RegisterTask(&taskq.TaskOptions{
		Name:       fmt.Sprintf("dataset#%s", datasetID),
		Handler:    s.findSimilarBloggers,
		RetryLimit: 5,
		MinBackoff: 5 * time.Second,
	})

	err = s.similarQueue.Add(task.WithArgs(ctx, datasetID, initialBloggers, botsPerDataset))
	if err != nil {
		return err
	}

	return nil

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

	bloggersLen := len(initialBloggers)
	bloggersPerBot := bloggersLen / len(bots)
	allBloggersProcessed := false

	wg := &sync.WaitGroup{}

	for i := range bots {
		wg.Add(1)
		rightBorderOfBatch := (i + 1) * botsPerDataset
		processorCtx := logger.WithFields(ctx, logger.Fields{"processor_index": i})

		if rightBorderOfBatch >= bloggersLen {
			rightBorderOfBatch = bloggersLen - 1
			allBloggersProcessed = true
		}

		bloggersBatch := initialBloggers[i*bloggersPerBot : rightBorderOfBatch]
		go s.findAndSaveSimilarBloggers(processorCtx, datasetID, bloggersBatch, bots[i], wg)

		if allBloggersProcessed && i != len(bots)-1 {
			logger.Warnf(ctx, "all bloggers processed with %d/%d bots", i+1, len(bots))
			break
		}
	}

	logger.Info(ctx, "started all goroutines, waiting for them")
	wg.Wait()

	logger.Infof(ctx, "all goroutines completed in %s", time.Since(startedAt))

	return nil
}

func (s Service) findAndSaveSimilarBloggers(
	ctx context.Context,
	datasetID uuid.UUID,
	initialBloggers []dbmodel.Blogger,
	bot dbmodel.Bot,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	startedAt := time.Now()
	ctx = logger.WithFields(ctx, logger.Fields{"bot_username": bot.Username})

	q := dbmodel.New(s.dbf(ctx))
	err := s.cli.CheckBot(ctx, bot.SessionID)
	if err != nil {
		logger.Errorf(ctx, "failed to check bot '%s': %v", bot.SessionID, err)
		if !errors.Is(err, instagrapi.ErrBotIsBlocked) {
			return
		}

		err2 := q.BlockBot(ctx, bot.ID)
		if err2 != nil {
			logger.Errorf(ctx, "failed to block bot (%s): %v", bot.ID, err)
			return
		}

		return
	}

	var users domain.InstUsers

	var count, totalCount int64

	for i, blogger := range initialBloggers {
		users, err = s.cli.FindSimilarBloggers(ctx, bot.SessionID, blogger.Username)
		if err != nil {
			logger.Errorf(ctx, "failed to find similar accounts: %v", err)
			continue
		}

		count, err = q.SaveBloggers(ctx, users.ToSaveBloggersParmas(datasetID))
		if err != nil {
			logger.Errorf(ctx, "failed to save parsed bloggers: %v", err)
			continue
		}

		totalCount += count

		logger.Infof(ctx, "saved %d bloggers from initial blogger '%s' (parsed %d/%d)", count, blogger.Username, i, len(initialBloggers))
	}

	logger.Infof(ctx, "saved %d new bloggers in %s", time.Since(startedAt))
}
