package gradingujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	grading_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian/grading"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
)

type ListUjianEssayUngradedService struct {
	listRepo grading_repo.ListGradingUjianRepository
}

func NewListUjianEssayUngradedService(listRepo grading_repo.ListGradingUjianRepository) *ListUjianEssayUngradedService {
	return &ListUjianEssayUngradedService{listRepo: listRepo}
}

func (r *ListUjianEssayUngradedService) ListUjianEssayUngraded(ctx context.Context, filter query.ListUjianEssayUngradedFilter) ([]ujian.ListUjian, error) {
	logger := corelog.FromContext(ctx)

	var err error
	filter, err = sanitizeAndValidateListUjianEssayUngradedFilter(filter)
	if err != nil {
		logger.Error(ctx, "failed list ungraded essay ujian", "layer", "core.service", "op", "ujian.grading.list_essay_ungraded.filter", "err", err)
		return nil, coreerror.ErrInvalidInput
	}

	items, err := r.listRepo.ListUjianEssayUngraded(ctx, filter)
	if err != nil {
		logger.Error(ctx, "failed list ungraded essay ujian", "layer", "core.service", "op", "ujian.grading.list_essay_ungraded", "err", err)
		return nil, err
	}

	return items, nil
}
