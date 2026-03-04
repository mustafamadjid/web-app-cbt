package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type UjianRepo struct {
	q Executor
	logger corelog.Logger
	pool *pgxpool.Pool
}

func NewUjianRepo(q Executor, logger corelog.Logger, pool *pgxpool.Pool) *UjianRepo {
	return &UjianRepo{q:q, logger: logger, pool: pool}
}

func (r *UjianRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}


func(r *UjianRepo)CreateUjian(ctx context.Context, ujian ujian.PenjadwalanUjian) error{
	
	tx,err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
	r.loggerFor(ctx).Error(ctx, "failed begin tx create ujian", "layer", "repo.db", "op", "ujian.create.begin_tx", "err", err)
		return err
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	var idUjian int64

	err = tx.QueryRow(ctx,
		`INSERT INTO ujian (
			id_bank_soal,
			id_kelas,
			id_nama_kelas,
			id_guru,
			nama_ujian,
			deskripsi_ujian,
			acak_soal
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id_ujian`,
		ujian.Ujian.IdBankSoal,
		ujian.Ujian.IdKelas,
		ujian.Ujian.IdNamaKelas,
		ujian.Ujian.IdGuru,
		ujian.Ujian.NamaUjian,
		ujian.Ujian.DeskripsiUjian,
		ujian.Ujian.AcakSoal,
	).Scan(&idUjian)

	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed to insert ujian", "layer", "repo.db", "op", "ujian.create.insert", "err", err)
		return err
	}
	
	_,err = tx.Exec(ctx,
		`INSERT INTO jadwal_ujian (
		id_ujian,
		id_sesi,
		id_ruangan,
		id_pengawas,
		tanggal_ujian,
		waktu_mulai,
		waktu_selesai,
		token,
		status_ujian)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		idUjian,
		ujian.JadwalUjian.IdSesi,
		ujian.JadwalUjian.IdRuangan,
		ujian.JadwalUjian.IdPengawas,
		ujian.JadwalUjian.TanggalUjian,
		ujian.JadwalUjian.WaktuMulai,
		ujian.JadwalUjian.WaktuSelesai,
		ujian.JadwalUjian.Token,
		ujian.JadwalUjian.StatusUjian,
	)

	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed to insert jadwal ujian", "layer", "repo.db", "op", "ujian.create.insert", "err", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed commit create ujian", "layer", "repo.db", "op", "ujian.create.commit", "err", err)
		return err
	}

	committed = true
	return nil
}

func(r *UjianRepo)CreatePesertaUjian(ctx context.Context, peserta ujian.PesertaUjian) (ujian.ID, error) {
	query := `
	INSERT INTO peserta_ujian (
		id_jadwal_ujian,
		id_siswa,
		waktu_mulai,
		waktu_submit,
		nilai_ujian
	)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id_peserta_ujian
	`

	var idPesertaUjian ujian.ID
	err := r.q.QueryRow(ctx, 
		query, 
		peserta.IdJadwalUjian,
		 peserta.IdSiswa, 
		 peserta.WaktuMulai, 
		 peserta.WaktuSubmit, 
		 peserta.NilaiUjian).Scan(&idPesertaUjian)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed creating peserta ujian", "layer", "repo.db", "op", "peserta_ujian.create", "err", err)
		return 0, err
	}

	return idPesertaUjian, nil
}

func(r *UjianRepo)CreateJawabanUjianSiswa(ctx context.Context, jawaban ujian.JawabanUjianSiswa) (ujian.ID, error) {
	query := `
	INSERT INTO jawaban_ujian_siswa (
		id_peserta_ujian,
		id_soal,
		id_pilihan,
		jawaban_essay,
		is_benar,
		waktu_jawab
	)
	VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING id_jawaban
	`

	var idJawaban ujian.ID
	err := r.q.QueryRow(ctx, 
		query, 
		jawaban.IdPesertaUjian,
		 jawaban.IdSoal, 
		 jawaban.IdPilihan, 
		 jawaban.JawabanEssay, 
		 jawaban.IsBenar, 
		 jawaban.WaktuJawab).Scan(&idJawaban)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed creating jawaban ujian siswa", "layer", "repo.db", "op", "jawaban_ujian_siswa.create", "err", err)
		return 0, err
	}

	return idJawaban, nil
}


