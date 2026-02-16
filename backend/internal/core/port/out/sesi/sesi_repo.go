package sesi_repo

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
)

type SesiRepository interface {
	GetSesi(ctx context.Context,filter query.ListSesiFilter) ([]sesi.Sesi, error)
	GetSesiById(ctx context.Context, idSesi int) (sesi.Sesi, error)
	GetSesiByKode(ctx context.Context, kodeSesi string) (sesi.Sesi, error)

	CreateSesi(ctx context.Context, sesi sesi.Sesi) error
	UpdateSesi(ctx context.Context, idSesi int, sesi updatepatch.UpdateSesiPatch) error
	DeleteSesi(ctx context.Context, idSesi int) error
}