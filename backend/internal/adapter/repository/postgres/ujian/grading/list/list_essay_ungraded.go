package gradingrepo

import (
	"context"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
)

type ListGradingRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewListGradingRepo(q pg.Executor, logger corelog.Logger) *ListGradingRepo {
	return &ListGradingRepo{q: q, logger: logger}
}

func (r *ListGradingRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *ListGradingRepo) ListUjianEssayUngraded(ctx context.Context, filter query.ListUjianEssayUngradedFilter) ([]ujian.ListUjian, error) {
	queryText, args := r.buildListUjianEssayUngradedQuery(filter)

	rows, err := r.q.Query(ctx, queryText, args...)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get list essay ungraded ujian", "layer", "repo.db", "op", "ujian.grading.list_essay_ungraded", "err", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanListUjianEssayUngradedRows(ctx, "ujian.grading.list_essay_ungraded", rows)
}
