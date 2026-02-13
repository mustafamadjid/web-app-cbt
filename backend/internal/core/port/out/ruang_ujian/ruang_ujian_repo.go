package ruangujian_repo

import (
	"context"

	ruangujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ruang_ujian"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ruang_ujian"
)

type RuangUjianRepo interface {
	GetRuangUjian(ctx context.Context, filter query.ListRuangUjianFilter)([]ruangujian.RuangUjian, error)
	GetRuangUjianById(ctx context.Context, idRuangan int)(ruangujian.RuangUjian, error)
	GetRuangUjianByKode(ctx context.Context, kodeRuang string)(ruangujian.RuangUjian, error)

	CreateRuangUjian(ctx context.Context, ruangUjian ruangujian.RuangUjian) error
	UpdateRuangUjian(ctx context.Context, idRuangan int, ruangUjian updatepatch.UpdateRuangUjianPatch) error
	DeleteRuangUjian(ctx context.Context, idRuangan int) error
}