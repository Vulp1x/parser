package queues

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/domain"
	"github.com/inst-api/parser/internal/instagrapi"
	"github.com/inst-api/parser/internal/mw"
	"github.com/inst-api/parser/pkg/logger"
	"go.uber.org/multierr"
)

// ErrBotIsBlocked бот заблокирован, необходимо воспользоваться другим
var ErrBotIsBlocked = errors.New("provided bot is blocked, need to choose another")

func (s Service) AddParseTargetsTask(ctx context.Context, datasetID uuid.UUID, bloggers []dbmodel.Blogger) (err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf(ctx, "recovered panic: %s, stack: %s", r, string(debug.Stack()))
			err = fmt.Errorf("recovered panic: %s", r)
		}
	}()

	err = s.parseTargetsQueue.Add(s.parseTargetsTask.WithArgs(ctx, datasetID, bloggers[0]))
	if err != nil {
		return err
	}

	if len(bloggers) == 1 {
		return nil
	}

	bloggers = bloggers[1:]
	bloggersPerBot := len(bloggers) / 5
	allBloggersProcessed := false

	for i := 0; i < 5; i++ {
		rightBorderOfBatch := (i + 1) * bloggersPerBot
		processorCtx := logger.WithFields(ctx, logger.Fields{"processor_index": i})

		var bloggersBatch []dbmodel.Blogger

		if rightBorderOfBatch >= len(bloggers) {
			bloggersBatch = bloggers[i*bloggersPerBot:]
			allBloggersProcessed = true
		} else {
			bloggersBatch = bloggers[i*bloggersPerBot : rightBorderOfBatch]
		}

		_ = s.parseTargetsQueue.Add(s.parseTargetsTask.WithArgs(processorCtx, datasetID, bloggersBatch))

		if allBloggersProcessed {
			logger.Warnf(ctx, "all bloggers were sent to bots in %d batches", i+1)
			break
		}
	}

	return nil

}

func (s Service) processParseTargetsFailedTask(ctx context.Context, datasetID uuid.UUID, _ []dbmodel.Blogger) error {
	logger.Info(ctx, "ParseTargets task failed, changing dataset status to ready for parsing")

	q := dbmodel.New(s.dbf(ctx))

	return q.UpdateDatasetStatus(ctx, dbmodel.UpdateDatasetStatusParams{Status: dbmodel.ReadyForParsingDatasetStatus, ID: datasetID})
}

func (s Service) parseTargetUsers(ctx context.Context, datasetID uuid.UUID, bloggers []dbmodel.Blogger) error {
	startedAt := time.Now()
	q := dbmodel.New(s.dbf(ctx))

	bots, err := q.LockAvailableBots(ctx, 1)
	if err != nil {
		return fmt.Errorf("failed to find available bots: %v", err)
	}

	if len(bots) == 0 {
		return ErrNoReadyBots
	}

	bot := bots[0]

	processorCtx := logger.WithFields(ctx, logger.Fields{"bot_username": bot.Username})
	err = s.parseAndSaveTargets(processorCtx, datasetID, bloggers, bot)
	if err != nil {
		return fmt.Errorf("")
	}
	err2 := q.UpdateDatasetStatus(ctx, dbmodel.UpdateDatasetStatusParams{Status: dbmodel.ReadyForParsingDatasetStatus, ID: datasetID})
	if err2 != nil {
		err2 = fmt.Errorf("failed to update dataset status to ReadyForParsingDatasetStatus (%d): %v", dbmodel.ReadyForParsingDatasetStatus, err2)
		logger.Error(ctx, err2)

		err = multierr.Append(err, err2)
	}

	logger.Infof(ctx, "all goroutines completed in %s", time.Since(startedAt))

	return err
}

func (s Service) parseAndSaveTargets(
	ctx context.Context,
	datasetID uuid.UUID,
	bloggersToParse []dbmodel.Blogger,
	bot dbmodel.Bot,
) (err error) {
	startedAt := time.Now()
	q := dbmodel.New(s.dbf(ctx))

	ctx = logger.WithKV(ctx, "uniq_task_processor", mw.ShortID())
	dataset, err := q.GetDatasetByID(ctx, datasetID)
	if err != nil {
		return fmt.Errorf("failed to find dataset: %v", err)
	}

	err = s.cli.CheckBot(ctx, bot.SessionID)
	if err != nil {
		logger.Errorf(ctx, "failed to check bot '%s': %v", bot.SessionID, err)

		if !errors.Is(err, instagrapi.ErrBotIsBlocked) {
			return err
		}

		err2 := q.BlockBot(ctx, bot.ID)
		if err2 != nil {
			return fmt.Errorf("failed to block bot (%s): %v", bot.ID, err)
		}

		return ErrBotIsBlocked
	}

	defer func() {
		err = q.UnlockBot(ctx, bot.ID)
		if err != nil {
			err = fmt.Errorf("failed to unlock bot after processing (%s): %v", bot.ID, err)
			return
		}

		logger.Info(ctx, "unlocked bot after processing")
	}()

	var users domain.ShortInstUsers

	var count, totalCount int64

	initialCtx := ctx
	var sleepDuration = time.Duration(25+rand.Intn(20)) * time.Second

	for i, blogger := range bloggersToParse {
		ctx = logger.WithKV(initialCtx, "blogger_username", blogger.Username)

		logger.Infof(ctx, "going to sleep for %s", sleepDuration)

		time.Sleep(sleepDuration)

		users, err = s.cli.ParseUsers(ctx, bot.SessionID, blogger.UserID, dataset)
		if err != nil {
			logger.Errorf(ctx, "failed to parse users: %v", err)

			if errors.Is(err, instagrapi.ErrBotIsBlocked) {
				err = q.BlockBot(ctx, bot.ID)
				if err != nil {
					logger.Errorf(ctx, "failed to block bot (%s): %v", bot.ID, err)
				}

				return ErrBotIsBlocked
			}

			if errors.Is(err, instagrapi.ErrBloggerNotFound) {
				err = q.SetBloggerIsParsed(ctx, dbmodel.SetBloggerIsParsedParams{IsCorrect: false, ID: blogger.ID})
				if err != nil {
					logger.Errorf(ctx, "failed to mark blogger as parsed  (%s): %v", bot.ID, err)
				}
			}

			if errors.Is(err, instagrapi.ErrToManyRequests) {
				logger.Warn(ctx, "going to sleep for 2 minutes, after 429 resp code")
				time.Sleep(2 * time.Minute)
				continue
			}
		}

		if len(users) == 0 {
			logger.Warn(ctx, "found 0 targets for blogger")
			continue
		}

		uniqueUsers := getUniqueusers(users)
		logger.Infof(ctx, "got %d unique users from %d parsed", len(uniqueUsers), len(users))

		count, err = q.SaveTargetUsers(ctx, uniqueUsers.ToSaveTargetsParams(datasetID))
		if err != nil {
			logger.Errorf(ctx, "failed to save targets from blogger: %v", err)
			continue
		}

		err = q.MarkBloggerAsParsed(ctx, blogger.ID)
		if err != nil {
			logger.Errorf(ctx, "failed to mark blogger as parsed: %v", err)
		}

		totalCount += count

		logger.Infof(ctx, "saved %d new target users (parsed %d/%d bloggers)", count, i+1, len(bloggersToParse))
	}

	logger.Infof(ctx, "saved %d new target in %s", totalCount, time.Since(startedAt))

	return nil
}

func getUniqueusers(users domain.ShortInstUsers) domain.ShortInstUsers {
	var uniqueUsers []domain.InstUserShort
	var m = make(map[int64]domain.InstUserShort)

	for _, user := range users {
		m[user.Pk] = user
	}

	for _, short := range m {
		uniqueUsers = append(uniqueUsers, short)
	}

	return uniqueUsers
}
