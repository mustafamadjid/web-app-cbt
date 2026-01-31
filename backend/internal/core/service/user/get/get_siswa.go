package user_service

import (
	"context"
	
	"strings"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
)

type GetSiswaService struct {
	siswaSvc out.GetListSiswaRepo
}

func NewGetListSiswaService(siswaSvc out.GetListSiswaRepo) *GetSiswaService {
	return &GetSiswaService{siswaSvc: siswaSvc}
}
var allowedSort = map[string]struct{}{
	"nama_lengkap":       {},
	"created_at": {},
	"username":   {},
	"nisn":        {},
}
func (s *GetSiswaService) ListSiswa(ctx context.Context, filter query.ListSiswaFilter) ([]query.SiswaListItem,error) {
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

	

	if _,ok := allowedSort[filter.SortBy]; !ok {
		return nil, coreerror.ErrInvalidInput
	}

	
	if filter.Status == nil {
		s := user.AKTIF
		filter.Status = &s
	}else {
		if *filter.Status != user.AKTIF && *filter.Status != user.NONAKTIF {
			return nil, coreerror.ErrInvalidInput
		}
	}

	if filter.Angkatan != nil {
		nowYear := time.Now().Year()
		if *filter.Angkatan > nowYear || *filter.Angkatan < 2019 {
			return nil, coreerror.ErrInvalidInput
		}
	}

	if filter.TingkatKelas != nil {
		if *filter.TingkatKelas < 0{
			return nil, coreerror.ErrInvalidInput
		}
	}

	if filter.JenisKelamin != nil {
		if *filter.JenisKelamin <= 0 || *filter.JenisKelamin > 2{
			return nil, coreerror.ErrInvalidInput
		}
	}

	items, err := s.siswaSvc.GetListSiswa(ctx,filter)
	if err != nil {
		return  nil,err
	}

	return items, nil
}