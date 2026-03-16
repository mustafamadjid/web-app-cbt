package siswaujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type jawabanAttemptOwnershipChecker interface {
	CheckAttemptOwnershipBySiswa(ctx context.Context, idSiswa int, idAttempt ujian.ID) (bool, error)
}

type SiswaGetJawabanUjianService struct {
	checker jawabanAttemptOwnershipChecker
	getter  *GetJawabanUjianService
}

func NewSiswaGetJawabanUjianService(checker jawabanAttemptOwnershipChecker, getter *GetJawabanUjianService) *SiswaGetJawabanUjianService {
	return &SiswaGetJawabanUjianService{
		checker: checker,
		getter:  getter,
	}
}

func (s *SiswaGetJawabanUjianService) GetJawabanUjianByAttemptId(ctx context.Context, idSiswa int, idAttempt ujian.ID) ([]ujian.JawabanUjian, error) {
	logger := corelog.FromContext(ctx)

	if idSiswa <= 0 {
		logger.Error(ctx, "failed get jawaban ujian by attempt id for siswa", "layer", "core.service", "op", "ujian.jawaban.get_by_attempt_id_siswa", "err", coreerror.ErrMissingId)
		return []ujian.JawabanUjian{}, coreerror.ErrMissingId
	}

	if idAttempt <= 0 {
		logger.Error(ctx, "failed get jawaban ujian by attempt id for siswa", "layer", "core.service", "op", "ujian.jawaban.get_by_attempt_id_siswa", "err", coreerror.ErrMissingId)
		return []ujian.JawabanUjian{}, coreerror.ErrMissingId
	}

	owned, err := s.checker.CheckAttemptOwnershipBySiswa(ctx, idSiswa, idAttempt)
	if err != nil {
		logger.Error(ctx, "failed get jawaban ujian by attempt id for siswa", "layer", "core.service", "op", "ujian.jawaban.get_by_attempt_id_siswa", "err", err)
		return []ujian.JawabanUjian{}, err
	}
	if !owned {
		logger.Error(ctx, "failed get jawaban ujian by attempt id for siswa", "layer", "core.service", "op", "ujian.jawaban.get_by_attempt_id_siswa", "err", coreerror.ErrNotFound)
		return []ujian.JawabanUjian{}, coreerror.ErrNotFound
	}

	items, err := s.getter.GetJawabanUjianByAttemptId(ctx, idAttempt)
	if err != nil {
		logger.Error(ctx, "failed get jawaban ujian by attempt id for siswa", "layer", "core.service", "op", "ujian.jawaban.get_by_attempt_id_siswa", "err", err)
		return []ujian.JawabanUjian{}, err
	}

	return items, nil
}
