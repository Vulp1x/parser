package queues

import (
	"time"

	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/instagrapi"
	"github.com/vmihailenco/taskq/v3"
	"github.com/vmihailenco/taskq/v3/memqueue"
)

const botLockDuration = 10 * time.Minute

type Service struct {
	similarQueue      taskq.Queue
	dbf               dbmodel.DBTXFunc
	cli               *instagrapi.Client
	findSimilarTask   *taskq.Task
	parseTargetsTask  *taskq.Task
	parseTargetsQueue *memqueue.Queue
}

func NewService(instagrapiHost string, dbf dbmodel.DBTXFunc) Service {
	simQueue := memqueue.NewQueue(&taskq.QueueOptions{
		Name:            "find_similar",
		ReservationSize: 1,
		Storage:         taskq.NewLocalStorage(),
	})

	parseTargetsQueue := memqueue.NewQueue(&taskq.QueueOptions{
		Name:            "parse_targets",
		ReservationSize: 1,
		Storage:         taskq.NewLocalStorage(),
	})

	service := Service{
		similarQueue:      simQueue,
		parseTargetsQueue: parseTargetsQueue,
		cli:               instagrapi.NewClient(instagrapiHost),
		dbf:               dbf,
	}

	task := taskq.RegisterTask(&taskq.TaskOptions{
		Name:            "find similar bloggers",
		Handler:         service.findSimilarBloggers,
		FallbackHandler: service.processFailedTask,
		RetryLimit:      5,
		MinBackoff:      5 * time.Second,
	})

	parseTargetsTask := taskq.RegisterTask(&taskq.TaskOptions{
		Name:            "parse targets from bloggers",
		Handler:         service.findSimilarBloggers,
		FallbackHandler: service.processFailedTask,
		RetryLimit:      5,
		MinBackoff:      5 * time.Second,
	})

	service.findSimilarTask = task
	service.parseTargetsTask = parseTargetsTask

	return service
}
