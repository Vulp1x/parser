package workers

import (
	"context"

	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/pkg/pgqueue"
)

type ParseUsersFromMediaHandler struct {
	dbTxF dbmodel.DBTXFunc
}

func (p *ParseUsersFromMediaHandler) HandleTask(ctx context.Context, task pgqueue.Task) error {
	// TODO implement me
	panic("implement me")
}
