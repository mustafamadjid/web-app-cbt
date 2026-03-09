package ujian_repo

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
)

type ListUjianRepository interface {
	GetAllUjian(ctx context.Context, filter query.ListUjianFilter) ([]ujian.ListUjian, error)
	GetUjianById(ctx context.Context, id ujian.ID) (ujian.ListUjian, error)
}

type UjianRepository interface {
	CreateUjian(ctx context.Context, ujian ujian.PenjadwalanUjian) error

	UpdateUjian(ctx context.Context, id ujian.ID, payload updatepatch.UpdatePenjadwalanUjian) error

	DeleteUjian(ctx context.Context, id ujian.ID) error
}

type SoalUjianRepository interface {
	GetSoalUjianByBankSoal(ctx context.Context, idBankSoal ujian.ID) ([]ujian.SoalUjianSiswa, error)
}
