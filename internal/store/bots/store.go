package bots

import (
	"context"

	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/domain"
)

func NewStore(dbtxFunc dbmodel.DBTXFunc, txFunc dbmodel.TxFunc) *Store {

	return &Store{
		dbtxf: dbtxFunc,
		txf:   txFunc,
	}
}

type Store struct {
	dbtxf dbmodel.DBTXFunc
	txf   dbmodel.TxFunc
}

func (s Store) SaveBots(ctx context.Context, bots domain.Bots) (int64, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	return q.SaveBots(ctx, bots.ToSaveParams())
}
