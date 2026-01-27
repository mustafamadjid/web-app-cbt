package in

import (
	"context"

	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
)

type ListProfilUseCase interface {
	ListGuru(ctx context.Context, filter query.ListGuruFilter) ([]query.GuruListItem,error)
}