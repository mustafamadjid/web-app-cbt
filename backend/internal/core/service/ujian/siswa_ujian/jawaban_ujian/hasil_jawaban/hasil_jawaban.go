package siswaujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
)

type HasilJawabanUjianService struct {
	hasilRepo ujian_repo.JawabanUjianRepository
}

func NewHasilJawabanUjianService(repo ujian_repo.JawabanUjianRepository) *HasilJawabanUjianService {
	return &HasilJawabanUjianService{
		hasilRepo: repo,
	}
}

func(r *HasilJawabanUjianService)ListHasilJawabanUjianByAttempt(ctx context.Context, idAttempt int)([]ujian.HasilJawabanUjian,error){
	logger := corelog.FromContext(ctx)

	if idAttempt <= 0 {
		logger.Error(ctx, "failed get jawaban ujian by attempt id", "layer", "core.service", "op", "ujian.jawaban.get_by_attempt_id", "err", coreerror.ErrMissingId)
		return []ujian.HasilJawabanUjian{}, coreerror.ErrMissingId
	}

	items, err := r.hasilRepo.ListHasilJawabanUjianByIdAttempt(ctx,ujian.ID(idAttempt))
	if err != nil {
		logger.Error(ctx, "failed get jawaban ujian by attempt id", "layer", "core.service", "op", "ujian.jawaban.get_by_attempt_id", "err", err)
		return []ujian.HasilJawabanUjian{}, err
	}

	return items, nil
}