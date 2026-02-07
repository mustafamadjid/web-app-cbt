package user_service

import (
	"context"
	"strings"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"

	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
)

type GetGuruService struct {
	guruSvc      out.GetGuruListRepo
	profilGuruSv out.ProfilGuruRepository
}

func NewGetListGuruService(guruSvc out.GetGuruListRepo, profilGuruSv out.ProfilGuruRepository) *GetGuruService {
	return &GetGuruService{guruSvc: guruSvc, profilGuruSv: profilGuruSv}
}

var allowedSortGuru = map[string]struct{}{
	"nama_lengkap": {},
	"created_at":   {},
	"username":     {},
	"nip":          {},
}

func (s *GetGuruService) ListGuru(ctx context.Context, filter query.ListGuruFilter) ([]query.GuruListItem, error) {
	filter.Search = strings.TrimSpace(filter.Search)

	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	if filter.Limit > 50 {
		filter.Limit = 50
	}

	if filter.Offset < 0 {
		filter.Offset = 0
	}

	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}

	if _, ok := allowedSortGuru[filter.SortBy]; !ok {
		return nil, coreerror.ErrInvalidInput
	}

	if filter.Status == nil {
		s := user.AKTIF
		filter.Status = &s
	} else {
		if *filter.Status != user.AKTIF && *filter.Status != user.NONAKTIF {
			return nil, coreerror.ErrInvalidInput
		}
	}

	if filter.Bidang != nil {
		trimmed := strings.TrimSpace(*filter.Bidang)
		if trimmed == "" {
			return nil, coreerror.ErrInvalidInput
		}
		filter.Bidang = &trimmed
	}

	items, err := s.guruSvc.GetListGuru(ctx, filter)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *GetGuruService) FindProfilGuruByID(ctx context.Context, id user.ID) (user.DataGuru, error) {
	return s.profilGuruSv.FindProfilGuruByID(ctx, id)
}
