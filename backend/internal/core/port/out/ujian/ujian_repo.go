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

type AttemptUjianRepository interface {
	GetAttemptUjianById(ctx context.Context, idAttempt ujian.ID) (ujian.AttemptUjian, error)
	CreateAttemptUjian(ctx context.Context, data ujian.AttemptUjian) error
	UpdateAttemptUjian(ctx context.Context, idAttempt ujian.ID, data updatepatch.UpdateAttemptUjianPatch) error
	DeleteAttemptUjian(ctx context.Context, idAttempt ujian.ID) error
}

type HasilUjianRepository interface {
	GetHasilUjianByAttempId(ctx context.Context, idAttempt ujian.ID) (ujian.HasilUjian, error)
	CreateHasilUjian(ctx context.Context, data ujian.HasilUjian) error
	UpdateHasilUjian(ctx context.Context, idAttempt ujian.ID, data updatepatch.UpdateHasilUjianPatch) error
	DeleteHasilUjian(ctx context.Context, idAttempt ujian.ID) error
}

type JawabanUjianRepository interface {
	GetJawabanUjianByAttemptId(ctx context.Context, idAttempt ujian.ID) ([]ujian.JawabanUjian, error)
	CreateJawabanUjian(ctx context.Context, data ujian.JawabanUjian) error
	UpdateJawabanUjian(ctx context.Context, idAttempt ujian.ID, data updatepatch.UpdateJawabanUjianPatch) error
	DeleteJawabanUjian(ctx context.Context, idAttempt ujian.ID) error
}