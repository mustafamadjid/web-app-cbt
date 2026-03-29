package ujian_repo

import (
	"context"
	"time"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
)

type SiswaUjianChecker interface {
	CheckAttemptOwnershipBySiswa(ctx context.Context, idSiswa int, idAttempt ujian.ID) (bool, error)

	CheckValidSiswaInPesertaUjianById(ctx context.Context, idSiswa int, idJadwalUjian int) (bool, int, error)
	CheckTokenUjian(ctx context.Context, token string, idJadwalUjian int) (bool, error)

	GetDeadlineUjian(ctx context.Context, idJadwalUjian int) (time.Time, error)
}

type UjianSiswaRepository interface {
	ListUjianSiswa(ctx context.Context, idSiswa int, filter query.ListUjianFilter) ([]ujian.ListUjian, error)
	GetWaktuSelesaiUjian(ctx context.Context, idJadwalUjian int) (time.Time, error)
	GetActiveUjianAttemptBySiswa(ctx context.Context, idSiswa int, idJadwalUjian int) (ujian.AttemptUjian, error)
}

