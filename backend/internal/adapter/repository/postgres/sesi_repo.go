package postgres

import (
	"context"

	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
)

type SesiRepo struct {
	q Executor
	logger corelog.Logger
}	

func NewSesirepo(q Executor, logger corelog.Logger) *SesiRepo {
	return &SesiRepo{q: q, logger: logger}
}

func (r *SesiRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}


func (r *SesiRepo)GetSesi(ctx context.Context,filter query.ListSesiFilter) ([]sesi.Sesi, error) {
	query := `
		SELECT
			
	`
}