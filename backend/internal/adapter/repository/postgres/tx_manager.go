package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out/tx"
)

type TxManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

func (t *TxManager) Begin(ctx context.Context) (out.Tx, error) {
	tx, error := t.pool.BeginTx(ctx, pgx.TxOptions{})
	if error != nil {
		return nil,error
	}

	return &pgTx{ctx: ctx, tx: tx}, nil
}