package ujianlistrepo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type ListSoalUjianRepo struct {
	q      pg.Executor
	logger corelog.Logger
	pool   *pgxpool.Pool
}

func NewListSoalUjianRepo(pool *pgxpool.Pool, logger corelog.Logger) *ListSoalUjianRepo {
	return &ListSoalUjianRepo{q: pool, logger: logger, pool: pool}
}

func (r *ListSoalUjianRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *ListSoalUjianRepo) GetSoalUjianByBankSoal(ctx context.Context, idBankSoal ujian.ID) ([]ujian.SoalUjianSiswa, error) {
	query := `
		SELECT
			s.id_soal,
			bs.id_bank_soal_version_aktif,
			s.tipe_soal,
			s.pertanyaan,
			s.gambar,
			s.bobot_soal,
			s.no_urut_soal,
			op.id_pilihan_ganda,
			op.isi_pilihan,
			op.is_benar
		FROM bank_soal bs
		JOIN isi_soal s
			ON s.id_bank_soal_version = bs.id_bank_soal_version_aktif
		LEFT JOIN opsi_pilihan_ganda op
			ON s.id_soal = op.id_soal
		WHERE bs.id_bank_soal = $1
		ORDER BY s.no_urut_soal, op.id_pilihan_ganda;
	`

	rows, err := r.q.Query(ctx, query, idBankSoal)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed listing soal ujian", "op", "soal_ujian.list", "err", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanSoalUjianRows(ctx, "soal_ujian.list", rows)
}

func (r *ListSoalUjianRepo) GetSoalUjianByBankSoalForSiswa(ctx context.Context, idJadwalUjian ujian.ID) ([]ujian.SoalUjianSiswa, bool, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed begin tx", "layer", "repo.db", "op", "ujian.list_soal.begin_tx", "err", err)
		return nil, false, err
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	var (
		idBankSoal ujian.ID
		acakSoal   bool
	)

	err = tx.QueryRow(
		ctx,
		`
			SELECT
				u.id_bank_soal,
				u.acak_soal
			FROM jadwal_ujian ju
			JOIN ujian u
				ON u.id_ujian = ju.id_ujian
			WHERE ju.id_jadwal_ujian = $1
		`,
		idJadwalUjian,
	).Scan(&idBankSoal, &acakSoal)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, false, coreerror.ErrNotFound
	}
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get bank soal for jadwal ujian", "layer", "repo.db", "op", "ujian.list_soal.get_bank_soal", "err", err)
		return nil, false, err
	}

	query := `
		SELECT
			s.id_soal,
			s.tipe_soal,
			s.pertanyaan,
			s.gambar,
			s.bobot_soal,
			s.no_urut_soal,
			op.id_pilihan_ganda,
			op.isi_pilihan
		FROM bank_soal bs
		JOIN isi_soal s
			ON s.id_bank_soal_version = bs.id_bank_soal_version_aktif
		LEFT JOIN opsi_pilihan_ganda op
			ON s.id_soal = op.id_soal
		WHERE bs.id_bank_soal = $1
		ORDER BY s.no_urut_soal, op.id_pilihan_ganda;
	`
	rows, err := tx.Query(ctx, query, idBankSoal)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed listing soal ujian", "op", "soal_ujian.list", "err", err)
		return nil, false, err
	}
	defer rows.Close()

	itemSoalResult, err := r.scanSoalUjianSiswaRows(ctx, "soal_ujian.list", rows)
	if err != nil {
		return nil, false, err
	}

	if err := tx.Commit(ctx); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed commit list soal ujian siswa", "layer", "repo.db", "op", "ujian.list_soal.commit", "err", err)
		return nil, false, err
	}

	committed = true
	return itemSoalResult, acakSoal, nil
}
