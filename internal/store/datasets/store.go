package datasets

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
)

const workersPerTask = 1

// ErrTaskNotFound не смогли найти таску
var ErrTaskNotFound = errors.New("task not found")

// ErrTaskInvalidStatus переход по статусам не возможен
var ErrTaskInvalidStatus = errors.New("invalid task status")

func NewStore(timeout time.Duration, dbtxFunc dbmodel.DBTXFunc, txFunc dbmodel.TxFunc, instagrapiHost string) *Store {
	return &Store{
		// tasksChan:   make(chan domain.Task, 10),
		taskCancels: make(map[uuid.UUID]func()),
		pushTimeout: timeout,
		dbtxf:       dbtxFunc,
		txf:         txFunc,
		taskMu:      &sync.Mutex{},
		// instaClient: instagrapi.NewClient(instagrapiHost),
	}
}

type Store struct {
	// tasksChan   chan domain.Task
	taskCancels map[uuid.UUID]func()
	taskMu      *sync.Mutex
	pushTimeout time.Duration
	dbtxf       dbmodel.DBTXFunc
	txf         dbmodel.TxFunc
	// instaClient instagrapiClient
}
