package user_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
)

func validateListGuruFilter(filter query.ListGuruFilter) (query.ListGuruFilter, error) {
	if _, ok := allowedSortGuru[filter.SortBy]; !ok {
		return filter, coreerror.ErrInvalidInput
	}
	if filter.Status == nil {
		defaultStatus := user.AKTIF
		filter.Status = &defaultStatus
	} else if *filter.Status != user.AKTIF && *filter.Status != user.NONAKTIF {
		return filter, coreerror.ErrInvalidInput
	}
	if filter.Bidang != nil && *filter.Bidang == "" {
		return filter, coreerror.ErrInvalidInput
	}
	return filter, nil
}
