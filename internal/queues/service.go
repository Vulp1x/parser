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
	similarQueue taskq.Queue
	dbf          dbmodel.DBTXFunc
	cli          *instagrapi.Client
}

func NewService(instagrapiHost string, dbf dbmodel.DBTXFunc) Service {
	q := memqueue.NewQueue(&taskq.QueueOptions{
		Name:            "find_similar",
		ReservationSize: 1,
		Storage:         taskq.NewLocalStorage(),
	})

	return Service{
		similarQueue: q,
		cli:          instagrapi.NewClient(instagrapiHost),
		dbf:          dbf,
	}
}
