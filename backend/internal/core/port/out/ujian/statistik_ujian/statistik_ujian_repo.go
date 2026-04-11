package statistikujian_repo

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

type StatistikUjianRepository interface {
	GetStatistikUjianByIdJadwal(ctx context.Context, idJadwal ujian.ID)(ujian.StatistikUjian,error)
}