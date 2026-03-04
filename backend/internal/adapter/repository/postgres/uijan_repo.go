package postgres

import (
	"context"

	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type UjianRepo struct {
	q Executor
	logger corelog.Logger
}

func NewUjianRepo(q Executor, logger corelog.Logger) *UjianRepo {
	return &UjianRepo{q:q, logger: logger}
}

func (r *UjianRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}


