package jawabanujian_repo

import (
	"context"
	"encoding/json"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type JawabanUjianRepo struct {
	q pg.Executor
	logger corelog.Logger
}

func NewJawabanUjianRepo(q pg.Executor, logger corelog.Logger) *JawabanUjianRepo {
	return &JawabanUjianRepo{q: q, logger: logger}
}

func (r *JawabanUjianRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func(r *JawabanUjianRepo) SaveJawabanUjian(ctx context.Context,idAttempt ujian.ID, jawaban []ujian.JawabanUjian) error{
	query := `
		INSERT INTO jawaban_ujian_siswa (
			id_attempt,
			id_soal,
			id_pilihan,
			jawaban_essay,
			waktu_jawab
			)
			SELECT
			$1,
			x.id_soal,
			x.id_pilihan,
			x.jawaban_essay,
			x.waktu_jawab
		FROM jsonb_to_recordset($2::jsonb) AS x(
			id_soal int,
			id_pilihan int,
			jawaban_essay text
			)
		ON CONFLICT (attempt_id, soal_id)
		DO UPDATE SET
			id_pilihan = EXCLUDED.id_pilihan,
			jawaban_essay = EXCLUDED.jawaban_essay
		`
	// Marshalling
	payloadJson, err := json.Marshal(jawaban)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed marshal jawaban ujian", "layer", "adapter.repository", "op", "ujian.jawaban.save", "err", err)
		return err
	}

	_,err = r.q.Exec(ctx,query,idAttempt,payloadJson)
	if err != nil{
		r.loggerFor(ctx).Error(ctx, "failed save jawaban ujian", "layer", "adapter.repository", "op", "ujian.jawaban.save", "err", err)
		return err
	}
	return nil
}
