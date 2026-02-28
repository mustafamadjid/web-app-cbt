package user_service

import (
	"context"

	"strings"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
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

// -----------------------
// Sanitizer and validator
// -----------------------

func sanitizeListSiswaFilter(filter query.ListSiswaFilter) query.ListSiswaFilter {
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

	return filter
}

func validateListSiswaFilter(filter query.ListSiswaFilter, nowYear int) (query.ListSiswaFilter, error) {
	if _, ok := allowedSort[filter.SortBy]; !ok {
		return filter, coreerror.ErrInvalidInput
	}

	if filter.Status == nil {
		defaultStatus := user.AKTIF
		filter.Status = &defaultStatus
	} else if *filter.Status != user.AKTIF && *filter.Status != user.NONAKTIF {
		return filter, coreerror.ErrInvalidInput
	}

	if filter.Angkatan != nil && (*filter.Angkatan > nowYear || *filter.Angkatan < 2019) {
		return filter, coreerror.ErrInvalidInput
	}

	if filter.TingkatKelas != nil && *filter.TingkatKelas < 0 {
		return filter, coreerror.ErrInvalidInput
	}

	if filter.JenisKelamin != nil && (*filter.JenisKelamin <= 0 || *filter.JenisKelamin > 2) {
		return filter, coreerror.ErrInvalidInput
	}

	return filter, nil
}
