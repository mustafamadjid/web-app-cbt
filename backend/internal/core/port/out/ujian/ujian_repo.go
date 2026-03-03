package ujian_repo

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
)

type ListUjianRepository interface {
	GetAllUjian(ctx context.Context, filter query.ListUjianFilter) ([]ujian.Ujian, error)
	GetUjianById(ctx context.Context, id ujian.ID) (ujian.Ujian, error)

	GetJadwalUjianById(ctx context.Context, id ujian.ID) (ujian.JadwalUjian, error)

	GetAllPesertaUjian(ctx context.Context, peserta ujian.PesertaUjian) ([]ujian.PesertaUjian, error)
	GetPesertaUjianBySiswa(ctx context.Context,idSiswa ujian.ID, peserta ujian.PesertaUjian) (ujian.PesertaUjian, error)

	GetAllJawabanUjianSiswa(ctx context.Context, jawaban ujian.JawabanUjianSiswa) ([]ujian.JawabanUjianSiswa, error)
	GetJawabanBySiswa(ctx context.Context,idSiswa ujian.ID, jawaban ujian.JawabanUjianSiswa) (ujian.JawabanUjianSiswa, error)
}

type UjianRepository interface {
	CreateUjian(ctx context.Context, ujian ujian.Ujian) (ujian.ID, error)
	CreateJadwalUjian(ctx context.Context, jadwal ujian.JadwalUjian) (ujian.ID, error)
	CreatePesertaUjian(ctx context.Context, peserta ujian.PesertaUjian) (ujian.ID, error)
	CreateJawabanUjianSiswa(ctx context.Context, jawaban ujian.JawabanUjianSiswa) (ujian.ID, error)

	UpdateUjian(ctx context.Context, id ujian.ID, payload updatepatch.UpdateUjianPatch) error
	UpdateJadwalUjian(ctx context.Context, id ujian.ID, payload updatepatch.UpdateJadwalUjianPatch) error
	UpdatePesertaUjian(ctx context.Context, id ujian.ID, payload updatepatch.UpdatePesertaUjianPatch) error
	UpdateJawabanUjianSiswa(ctx context.Context, id ujian.ID, payload updatepatch.UpdateJawabanUjianSiswaPatch) error

	DeleteUjian(ctx context.Context, id ujian.ID) error
	DeleteJadwalUjian(ctx context.Context, id ujian.ID) error
	DeletePesertaUjian(ctx context.Context, id ujian.ID) error
	DeleteJawabanUjianSiswa(ctx context.Context, id ujian.ID) error
}
