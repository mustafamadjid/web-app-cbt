package ujian_repo

import (
	"context"
	"time"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
)

type SiswaUjianChecker interface {
	CheckValidSiswaInPesertaUjianById(ctx context.Context, idSiswa int, idJadwalUjian int) (bool, int, error)
	CheckTokenUjian(ctx context.Context,token string, idJadwalUjian int) (bool, error)

	GetDeadlineUjian(ctx context.Context, idJadwalUjian int) (time.Time, error)
}

type ListUjianSiswaRepository interface {
	ListUjianSiswa(ctx context.Context, idSiswa int, filter query.ListUjianFilter) ([]ujian.ListUjian, error)
}

type WaktuSelesaiUjianRepository interface {
	GetWaktuSelesaiUjian(ctx context.Context, idJadwalUjian int) (time.Time, error)
}
