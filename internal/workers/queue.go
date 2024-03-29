package workers

import (
	"context"
	"fmt"
	"time"

	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/pb/instaproxy"
	"github.com/inst-api/parser/pkg/pgqueue"
	"github.com/inst-api/parser/pkg/pgqueue/pkg/delayer"
	"github.com/inst-api/parser/pkg/pgqueue/pkg/executor"
	"google.golang.org/grpc"
)

const (
	FindSimilarBloggersTaskKind   = 1
	ParseBloggersMediaTaskKind    = 2
	ParseUsersFromMediaTaskKind   = 3
	TransitToSimilarFoundTaskKind = 4
	ParseFollowersTaskKind        = 5
	ParseFullUsersTaskKind        = 6
	// PrepareParseFollowersTaskKind готоваит блогера для парсинга
	PrepareParseFollowersTaskKind = 7
	// TransitToCompletedTaskKind выставляем конечный статус для датасета
	TransitToCompletedTaskKind = 8
)

var EmptyPayload = []byte(`{"empty":true}`)

func NewQueuue(ctx context.Context, executor executor.Executor, txFunc dbmodel.DBTXFunc, conn *grpc.ClientConn) *pgqueue.Queue {
	queue := pgqueue.NewQueue(ctx, executor)

	// ищем похожих блогеров на начальных блогеров
	queue.RegisterKind(FindSimilarBloggersTaskKind, &SimilarBloggersHandler{dbTxF: txFunc, cli: instaproxy.NewInstaProxyClient(conn)}, pgqueue.KindOptions{
		Name:                 "similar-bloggers",
		WorkerCount:          pgqueue.NewConstProvider(int16(5)),
		MaxAttempts:          10,
		AttemptTimeout:       40 * time.Second,
		MaxTaskErrorMessages: 10,
		Delayer:              delayer.NewJitterDelayer(delayer.EqualJitter, 5*time.Second),
		TerminalTasksTTL:     pgqueue.NewConstProvider(1000 * time.Hour),
		Loop: pgqueue.LoopOptions{
			JanitorPeriod: pgqueue.NewConstProvider(15 * time.Hour),
			FetcherPeriod: pgqueue.NewConstProvider(5 * time.Second),
		},
	})

	// ищем посты для дальнейшего парсинга
	queue.RegisterKind(ParseBloggersMediaTaskKind, &ParseMediasHandler{dbTxF: txFunc, cli: instaproxy.NewInstaProxyClient(conn), queue: queue}, pgqueue.KindOptions{
		Name:                 "find-medias",
		WorkerCount:          pgqueue.NewConstProvider(int16(10)),
		MaxAttempts:          10,
		AttemptTimeout:       30 * time.Second,
		MaxTaskErrorMessages: 10,
		Delayer:              delayer.NewJitterDelayer(delayer.EqualJitter, 15*time.Second),
		TerminalTasksTTL:     pgqueue.NewConstProvider(1000 * time.Hour),
		Loop: pgqueue.LoopOptions{
			JanitorPeriod: pgqueue.NewConstProvider(15 * time.Hour),
			FetcherPeriod: pgqueue.NewConstProvider(5 * time.Second),
		},
	})

	// парсим комментаторов из конкретного поста у блоггера
	queue.RegisterKind(ParseUsersFromMediaTaskKind, &ParseUsersFromMediaHandler{dbTxF: txFunc, cli: instaproxy.NewInstaProxyClient(conn)}, pgqueue.KindOptions{
		Name:                 "parse-targets",
		WorkerCount:          pgqueue.NewConstProvider(int16(40)),
		MaxAttempts:          10,
		AttemptTimeout:       30 * time.Second,
		MaxTaskErrorMessages: 10,
		Delayer:              delayer.NewJitterDelayer(delayer.EqualJitter, 15*time.Second),
		TerminalTasksTTL:     pgqueue.NewConstProvider(1000 * time.Hour),
		Loop: pgqueue.LoopOptions{
			JanitorPeriod: pgqueue.NewConstProvider(15 * time.Hour),
			FetcherPeriod: pgqueue.NewConstProvider(3 * time.Second),
		},
	})

	// готовим блоггера для парсинга, нужно получать полную информацию по нему и сохранить
	queue.RegisterKind(PrepareParseFollowersTaskKind, &PrepareParseFollowersHandler{dbTxF: txFunc, queue: queue, cli: instaproxy.NewInstaProxyClient(conn)}, pgqueue.KindOptions{
		Name:                 "prepare-parsing-followers",
		WorkerCount:          pgqueue.NewConstProvider(int16(10)),
		MaxAttempts:          10,
		AttemptTimeout:       10 * time.Second,
		MaxTaskErrorMessages: 20,
		Delayer:              delayer.NewJitterDelayer(delayer.EqualJitter, 3*time.Second),
		TerminalTasksTTL:     pgqueue.NewConstProvider(1000 * time.Hour),
		Loop: pgqueue.LoopOptions{
			JanitorPeriod: pgqueue.NewConstProvider(15 * time.Hour),
			FetcherPeriod: pgqueue.NewConstProvider(3 * time.Second),
		},
	})

	// парсим подписчиков у блоггера
	queue.RegisterKind(ParseFollowersTaskKind, &ParseFollowersHandler{dbTxF: txFunc, queue: queue, cli: instaproxy.NewInstaProxyClient(conn)}, pgqueue.KindOptions{
		Name:                 "parse-followers",
		WorkerCount:          pgqueue.NewConstProvider(int16(100)),
		MaxAttempts:          100,
		AttemptTimeout:       40 * time.Second,
		MaxTaskErrorMessages: 20,
		Delayer:              delayer.NewJitterDelayer(delayer.EqualJitter, 3*time.Second),
		TerminalTasksTTL:     pgqueue.NewConstProvider(1000 * time.Hour),
		Loop: pgqueue.LoopOptions{
			JanitorPeriod: pgqueue.NewConstProvider(15 * time.Hour),
			FetcherPeriod: pgqueue.NewConstProvider(3 * time.Second),
		},
	})

	// парсим комментаторов из конкретного поста у блоггера
	queue.RegisterKind(ParseFullUsersTaskKind, &ParseUsersFromMediaHandler{dbTxF: txFunc, cli: instaproxy.NewInstaProxyClient(conn)}, pgqueue.KindOptions{
		Name:                 "parse-full-users",
		WorkerCount:          pgqueue.NewConstProvider(int16(1)),
		MaxAttempts:          10,
		AttemptTimeout:       40 * time.Second,
		MaxTaskErrorMessages: 10,
		Delayer:              delayer.NewJitterDelayer(delayer.EqualJitter, 15*time.Second),
		TerminalTasksTTL:     pgqueue.NewConstProvider(1000 * time.Hour),
		Loop: pgqueue.LoopOptions{
			JanitorPeriod: pgqueue.NewConstProvider(15 * time.Hour),
			FetcherPeriod: pgqueue.NewConstProvider(3 * time.Second),
		},
	})

	// переводим датасет в статус парсинг закончен после того, как все блогеры распаршены
	queue.RegisterKind(TransitToSimilarFoundTaskKind, &TransitToSimilarFoundHandler{dbTxF: txFunc, queue: queue}, pgqueue.KindOptions{
		Name:                 "similar-found",
		WorkerCount:          pgqueue.NewConstProvider(int16(5)),
		MaxAttempts:          100,
		AttemptTimeout:       20 * time.Second,
		MaxTaskErrorMessages: 10,
		Delayer:              delayer.NewJitterDelayer(delayer.EqualJitter, 50*time.Second),
		TerminalTasksTTL:     pgqueue.NewConstProvider(1000 * time.Hour),
		Loop: pgqueue.LoopOptions{
			JanitorPeriod: pgqueue.NewConstProvider(15 * time.Hour),
			FetcherPeriod: pgqueue.NewConstProvider(10 * time.Second),
		},
	})

	// переводим датасет в статус парсинг закончен после того, как все блогеры распаршены
	queue.RegisterKind(TransitToCompletedTaskKind, &TransitToTaskCompleteddHandler{dbTxF: txFunc, queue: queue}, pgqueue.KindOptions{
		Name:                 "all-parsed",
		WorkerCount:          pgqueue.NewConstProvider(int16(5)),
		MaxAttempts:          100,
		AttemptTimeout:       20 * time.Second,
		MaxTaskErrorMessages: 10,
		Delayer:              delayer.NewJitterDelayer(delayer.EqualJitter, 50*time.Second),
		TerminalTasksTTL:     pgqueue.NewConstProvider(1000 * time.Hour),
		Loop: pgqueue.LoopOptions{
			JanitorPeriod: pgqueue.NewConstProvider(15 * time.Hour),
			FetcherPeriod: pgqueue.NewConstProvider(10 * time.Second),
		},
	})

	queue.Start()

	return queue
}

// TaskKindFromDataset по датасету получаем тип задачи для парсинга
func TaskKindFromDataset(dataset dbmodel.Dataset) int16 {
	var taskKind int16
	switch dataset.Type {
	case dbmodel.DatasetTypeLikesAndComments:
		taskKind = ParseBloggersMediaTaskKind
	case dbmodel.DatasetTypeFollowers:
		taskKind = PrepareParseFollowersTaskKind
	case dbmodel.DatasetTypePhoneNumbers:
		taskKind = ParseFullUsersTaskKind
	default:
		panic(fmt.Sprintf("unexpected dataset type '%s'", dataset.Type))
	}

	return taskKind
}
