package user_service

import (
	"context"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
	"time"
)

type GetSiswaService struct {
	siswaSvc      out.GetListSiswaRepo
	profilSiswaSv out.ProfilSiswaRepository
}

func NewGetListSiswaService(siswaSvc out.GetListSiswaRepo, profilSiswaSv out.ProfilSiswaRepository) *GetSiswaService {
	return &GetSiswaService{siswaSvc: siswaSvc, profilSiswaSv: profilSiswaSv}
}

var allowedSort = map[string]struct{}{
	"nama_lengkap": {},
	"created_at":   {},
	"username":     {},
	"nisn":         {},
}

func (s *GetSiswaService) ListSiswa(ctx context.Context, filter query.ListSiswaFilter) ([]query.SiswaListItem, error) {
	logger := corelog.FromContext(ctx)

	filter = sanitizeListSiswaFilter(filter)

	var err error
	filter, err = validateListSiswaFilter(filter, time.Now().Year())
	if err != nil {
		return nil, err
	}

	items, err := s.siswaSvc.GetListSiswa(ctx, filter)
	if err != nil {
		logger.Error(ctx, "failed get list siswa", "layer", "core.service", "op", "user.get", "err", err)
		return nil, err
	}

	return items, nil
}

func (s *GetSiswaService) FindProfilSiswaByID(ctx context.Context, id user.ID) (user.DataSiswa, error) {
	return s.profilSiswaSv.FindProfilSiswaByID(ctx, id)
}
