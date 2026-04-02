package gradingrepo

import (
	"context"

	"github.com/jackc/pgx/v5"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

type gradingJobScanner interface {
	Scan(dest ...any) error
}

func scanGradingJobRow(row gradingJobScanner) (ujian.GradingJob, error) {
	var item ujian.GradingJob
	if err := row.Scan(
		&item.IDgradingJob,
		&item.IDAttempt,
		&item.Status,
		&item.RetryCount,
		&item.MaxRetries,
		&item.AvailableAt,
		&item.LockedAt,
		&item.ErrorCode,
		&item.ErrorMessage,
	); err != nil {
		return ujian.GradingJob{}, err
	}

	return item, nil
}

func (r *GradingRepo) scanGradingJobRows(ctx context.Context, op string, rows pgx.Rows) ([]ujian.GradingJob, error) {
	var items []ujian.GradingJob
	for rows.Next() {
		item, err := scanGradingJobRow(rows)
		if err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning grading job", "layer", "repo.db", "op", op, "err", err)
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating grading job", "layer", "repo.db", "op", op, "err", err)
		return nil, err
	}

	return items, nil
}

func buildStatistikSoalPayload(items []ujian.StatistikSoal, isBenar bool) []statistikSoalUpsertItem {
	if len(items) == 0 {
		return nil
	}

	payload := make([]statistikSoalUpsertItem, 0, len(items))

	for _, item := range items {
		entry := statistikSoalUpsertItem{
			IDSoal:  item.IDSoal,
			IDUjian: item.IDUjian,
		}

		if isBenar {
			entry.JumlahJawabanBenar++
		} else {
			entry.JumlahJawabanSalah++
		}

		payload = append(payload, entry)
	}

	return payload
}
