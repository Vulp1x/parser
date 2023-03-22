package workers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/dbtx"
	"github.com/inst-api/parser/internal/domain"
	"github.com/inst-api/parser/internal/pb/instaproxy"
	"github.com/inst-api/parser/pkg/logger"
	"github.com/inst-api/parser/pkg/pgqueue"
)

// размер одного батча для парсинга подписчиков
const defaultFollowersLimit int32 = 500

type ParseFollowersHandler struct {
	dbTxF dbmodel.DBTXFunc
	cli   instaproxy.InstaProxyClient
	queue *pgqueue.Queue
}

func (h *ParseFollowersHandler) HandleTask(ctx context.Context, task pgqueue.Task) error {
	taskKey, err := h.parseTaskKey(ctx, task)
	if err != nil {
		return err
	}

	ctx = logger.WithFields(ctx, logger.Fields{"blogger": taskKey.bloggerPk, "dataset_id": taskKey.datasetID})

	db := h.dbTxF(ctx)
	q := dbmodel.New(db)
	dataset, err := q.GetDatasetByID(ctx, taskKey.datasetID)
	if err != nil {
		return fmt.Errorf("failed to find datatset with id '%s': %v", taskKey.datasetID, err)
	}

	if dataset.FollowersCount <= 0 {
		return fmt.Errorf("%w: got %d followers count for dataset6 expected at least 0", pgqueue.ErrMustCancelTask, dataset.FollowersCount)
	}

	var followersLimit = defaultFollowersLimit
	if followersLimit > dataset.FollowersCount {
		followersLimit = dataset.FollowersCount
	}

	targetsResp, err := h.cli.ParseFollowers(ctx, &instaproxy.ParseFollowersRequest{
		BloggerPk: taskKey.bloggerPk, FollowersCount: followersLimit,
	})
	if err != nil {
		return fmt.Errorf("failed to parse targets from media '%d': %v", taskKey.bloggerPk, err)
	}

	return h.savedParsedFollowers(ctx, dataset, targetsResp, taskKey)
}

func (h *ParseFollowersHandler) savedParsedFollowers(ctx context.Context, dataset dbmodel.Dataset, targetsResp *instaproxy.ParsingResponse, taskKey parseFollowersTaskKey) error {
	tx, err := h.dbTxF(ctx).Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer dbtx.RollbackUnlessCommitted(ctx, tx)
	q := dbmodel.New(tx)

	targets := domain.ShortUsersFromProto(targetsResp.GetTargets())

	count, err := q.SaveTargetUsers(ctx, targets.ToSaveTargetsParams(targetsResp.BloggerPk, taskKey.datasetID))
	if err != nil {
		return fmt.Errorf("failed to save targets: %v", err)
	}

	parsedFollowers, err := q.CountParsedTargets(ctx, dbmodel.CountParsedTargetsParams{
		DatasetID: taskKey.datasetID,
		MediaPk:   targetsResp.BloggerPk,
	})
	if err != nil {
		return fmt.Errorf("failed to count parsed followers: %v", err)
	}
	if int32(parsedFollowers) < dataset.FollowersCount {
		taskKey.orderNumber++
		logger.Infof(ctx, "going to add new %d task for parsing followers", taskKey.orderNumber)

		if err = h.queue.PushTaskTx(ctx, tx, taskKey.newParsingTask(), pgqueue.WithDelay(3*time.Second)); err != nil {
			return fmt.Errorf("failed to push next task to queue: %v", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	logger.Infof(ctx, "saved %d/%d targets", count, len(targets))

	return nil
}

func (h *ParseFollowersHandler) parseTaskKey(ctx context.Context, task pgqueue.Task) (parseFollowersTaskKey, error) {
	logger.Infof(ctx, "starting processing task %s", task.ExternalKey)

	externalKeyParts := strings.Split(task.ExternalKey, "::")
	if len(externalKeyParts) != 3 {
		return parseFollowersTaskKey{}, fmt.Errorf(
			`%w: expected external key in format "<dataset id>::<blogger username>::<order number>", got '%s'`,
			pgqueue.ErrMustCancelTask, task.ExternalKey)
	}

	datasetID, err := uuid.Parse(externalKeyParts[0])
	if err != nil {
		return parseFollowersTaskKey{}, fmt.Errorf("%w: failed to parse datasaet id from '%s': %v", pgqueue.ErrMustCancelTask, externalKeyParts[0], err)
	}

	bloggerPk, err := strconv.ParseInt(externalKeyParts[1], 10, 64)
	if err != nil {
		return parseFollowersTaskKey{}, fmt.Errorf("%w: failed to parse blogger pk from '%s': %v",
			pgqueue.ErrMustCancelTask, externalKeyParts[1], err)
	}

	orderNumber, err := strconv.ParseInt(externalKeyParts[2], 10, 32)
	if err != nil {
		return parseFollowersTaskKey{}, fmt.Errorf("%w: failed to parse order number from '%s': %v",
			pgqueue.ErrMustCancelTask, externalKeyParts[2], err)
	}

	return parseFollowersTaskKey{datasetID: datasetID, bloggerPk: bloggerPk, orderNumber: int(orderNumber)}, nil
}

type parseFollowersTaskKey struct {
	datasetID   uuid.UUID
	bloggerPk   int64
	orderNumber int
}

func (p parseFollowersTaskKey) String() string {
	return fmt.Sprintf("%s::%d::%d", p.datasetID.String(), p.bloggerPk, p.orderNumber)
}

func (p parseFollowersTaskKey) newParsingTask() pgqueue.Task {
	return pgqueue.Task{Kind: ParseFollowersTaskKind, Payload: EmptyPayload, ExternalKey: p.String()}
}
