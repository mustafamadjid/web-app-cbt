package analisissoal_repo

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

type AnalisisSoalInterface interface {
	GetListAnalisisSoal(ctx context.Context, idJadwalUjian ujian.ID) ([]ujian.AnalisisSoal, error)
}