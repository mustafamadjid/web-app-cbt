package in

import (
	"context"

	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
)

type ListProfilSiswaUseCase interface {
	ListSiswa(ctx context.Context, filter query.ListSiswaFilter) ([]query.SiswaListItem,error)
}