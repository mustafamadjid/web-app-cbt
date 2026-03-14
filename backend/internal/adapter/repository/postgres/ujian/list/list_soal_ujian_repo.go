package ujianlistrepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres"
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

	var (
		itemSoalResult []ujian.SoalUjianSiswa
		orderedSoalIDs []ujian.ID
		gambar         sql.NullString
		idPilgan       sql.NullInt64
		isiPilgan      sql.NullString
		isBenar        sql.NullBool
	)

	rows, err := r.q.Query(ctx, query, idBankSoal)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed listing soal ujian", "op", "soal_ujian.list", "err", err)
		return nil, err
	}
	defer rows.Close()

	itemsBySoalID := make(map[ujian.ID]*ujian.SoalUjianSiswa)
	for rows.Next() {
		var (
			idSoal            ujian.ID
			idBankSoalVersion ujian.ID
			tipeSoal          string
			pertanyaan        string
			bobotSoal         int
			noUrutSoal        int
		)

		if err := rows.Scan(
			&idSoal,
			&idBankSoalVersion,
			&tipeSoal,
			&pertanyaan,
			&gambar,
			&bobotSoal,
			&noUrutSoal,
			&idPilgan,
			&isiPilgan,
			&isBenar,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning soal ujian", "op", "soal_ujian.list", "err", err)
			return nil, err
		}

		item, exists := itemsBySoalID[idSoal]
		if !exists {
			item = &ujian.SoalUjianSiswa{
				IdSoal:            idSoal,
				IdBankSoalVersion: idBankSoalVersion,
				TipeSoal:          tipeSoal,
				Pertanyaan:        pertanyaan,
				BobotSoal:         bobotSoal,
				NoUrutSoal:        noUrutSoal,
			}

			if gambar.Valid {
				item.Gambar = gambar.String
			}

			itemsBySoalID[idSoal] = item
			orderedSoalIDs = append(orderedSoalIDs, idSoal)
		}

		if idPilgan.Valid {
			opsi := ujian.OpsiPilganUjian{
				IdPilihanGanda: ujian.ID(idPilgan.Int64),
				IdSoal:         item.IdSoal,
			}

			if isiPilgan.Valid {
				opsi.IsiPilihan = isiPilgan.String
			}

			if isBenar.Valid {
				opsi.IsBenar = isBenar.Bool
			}

			item.OpsiJawaban = append(item.OpsiJawaban, opsi)
		}
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating soal ujian rows", "op", "soal_ujian.list", "err", err)
		return nil, err
	}

	for _, idSoal := range orderedSoalIDs {
		itemSoalResult = append(itemSoalResult, *itemsBySoalID[idSoal])
	}

	return itemSoalResult, nil
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
	var (
		itemSoalResult []ujian.SoalUjianSiswa
		orderedSoalIDs []ujian.ID
		gambar         sql.NullString
		idPilgan       sql.NullInt64
		isiPilgan      sql.NullString
	)
	rows, err := tx.Query(ctx, query, idBankSoal)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed listing soal ujian", "op", "soal_ujian.list", "err", err)
		return nil, false, err
	}
	defer rows.Close()

	itemsBySoalID := make(map[ujian.ID]*ujian.SoalUjianSiswa)
	for rows.Next() {
		var (
			idSoal     ujian.ID
			tipeSoal   string
			pertanyaan string
			bobotSoal  int
			noUrutSoal int
		)

		if err := rows.Scan(
			&idSoal,
			&tipeSoal,
			&pertanyaan,
			&gambar,
			&bobotSoal,
			&noUrutSoal,
			&idPilgan,
			&isiPilgan,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning soal ujian", "op", "soal_ujian.list", "err", err)
			return nil, false, err
		}

		item, exists := itemsBySoalID[idSoal]
		if !exists {
			item = &ujian.SoalUjianSiswa{
				IdSoal:     idSoal,
				TipeSoal:   tipeSoal,
				Pertanyaan: pertanyaan,
				BobotSoal:  bobotSoal,
				NoUrutSoal: noUrutSoal,
			}

			if gambar.Valid {
				item.Gambar = gambar.String
			}

			itemsBySoalID[idSoal] = item
			orderedSoalIDs = append(orderedSoalIDs, idSoal)
		}

		if idPilgan.Valid {
			opsi := ujian.OpsiPilganUjian{
				IdPilihanGanda: ujian.ID(idPilgan.Int64),
			}

			if isiPilgan.Valid {
				opsi.IsiPilihan = isiPilgan.String
			}

			item.OpsiJawaban = append(item.OpsiJawaban, opsi)
		}
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating soal ujian rows", "op", "soal_ujian.list", "err", err)
		return nil, false, err
	}

	for _, idSoal := range orderedSoalIDs {
		itemSoalResult = append(itemSoalResult, *itemsBySoalID[idSoal])
	}

	if err := tx.Commit(ctx); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed commit list soal ujian siswa", "layer", "repo.db", "op", "ujian.list_soal.commit", "err", err)
		return nil, false, err
	}

	committed = true
	return itemSoalResult, acakSoal, nil
}
