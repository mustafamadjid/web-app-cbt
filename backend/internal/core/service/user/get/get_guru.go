package user_service

import (
	"context"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
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
	filter = sanitizeListGuruFilter(filter)

	var err error
	filter, err = validateListGuruFilter(filter)
	if err != nil {
		return nil, err
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
