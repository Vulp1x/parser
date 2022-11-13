package datasets

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/queues"
)

const workersPerTask = 1

// ErrDatasetInvalidStatus переход по статусам не возможен
var ErrDatasetInvalidStatus = errors.New("invalid task status")

func NewStore(timeout time.Duration, dbtxFunc dbmodel.DBTXFunc, txFunc dbmodel.TxFunc, instagrapiHost string) *Store {

	return &Store{
		// tasksChan:   make(chan domain.Task, 10),
		taskCancels:  make(map[uuid.UUID]func()),
		pushTimeout:  timeout,
		dbtxf:        dbtxFunc,
		txf:          txFunc,
		taskMu:       &sync.Mutex{},
		queueService: queues.NewService(instagrapiHost, dbtxFunc),
		// instaClient: instagrapi.NewClient(instagrapiHost),
	}
}

type Store struct {
	// tasksChan   chan domain.Task
	taskCancels  map[uuid.UUID]func()
	taskMu       *sync.Mutex
	pushTimeout  time.Duration
	dbtxf        dbmodel.DBTXFunc
	txf          dbmodel.TxFunc
	queueService queues.Service

	// instaClient instagrapiClient
}
