package ujianrepo

import (
	"context"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type JawabanUjianRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewJawabanUjianRepo(q pg.Executor, logger corelog.Logger) *JawabanUjianRepo {
	return &JawabanUjianRepo{q: q, logger: logger}
}

func (r *JawabanUjianRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *JawabanUjianRepo) CreateJawabanUjianSiswa(ctx context.Context, jawaban ujian.JawabanUjianSiswa) (ujian.ID, error) {
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
	err := r.q.QueryRow(
		ctx,
		query,
		jawaban.IdPesertaUjian,
		jawaban.IdSoal,
		jawaban.IdPilihan,
		jawaban.JawabanEssay,
		jawaban.IsBenar,
		jawaban.WaktuJawab,
	).Scan(&idJawaban)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed creating jawaban ujian siswa", "layer", "repo.db", "op", "jawaban_ujian_siswa.create", "err", err)
		return 0, err
	}

	return idJawaban, nil
}

func (r *JawabanUjianRepo) UpdateJawabanUjianSiswa(ctx context.Context, id ujian.ID, payload updatepatch.UpdateJawabanUjianSiswaPatch) error {
	tag, err := r.q.Exec(
		ctx,
		`UPDATE jawaban_ujian_siswa
		SET
			id_peserta_ujian = COALESCE($1, id_peserta_ujian),
			id_soal = COALESCE($2, id_soal),
			id_pilihan = COALESCE($3, id_pilihan),
			jawaban_essay = COALESCE($4, jawaban_essay),
			is_benar = COALESCE($5, is_benar),
			waktu_jawab = COALESCE($6, waktu_jawab)
		WHERE id_jawaban = $7`,
		payload.IdPesertaUjian,
		payload.IdSoal,
		payload.IdPilihan,
		payload.JawabanEssay,
		payload.IsBenar,
		payload.WaktuJawab,
		id,
	)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed updating jawaban ujian siswa", "layer", "repo.db", "op", "jawaban_ujian_siswa.update", "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}

	return nil
}

func (r *JawabanUjianRepo) DeleteJawabanUjianSiswa(ctx context.Context, id ujian.ID) error {
	tag, err := r.q.Exec(ctx, `DELETE FROM jawaban_ujian_siswa WHERE id_jawaban = $1`, id)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed deleting jawaban ujian siswa", "layer", "repo.db", "op", "jawaban_ujian_siswa.delete", "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}

	return nil
}
