package user_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
)

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
