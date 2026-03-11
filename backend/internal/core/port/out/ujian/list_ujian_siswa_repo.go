package ujian_repo

import (
	"context"
	// "time"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
)

type ListUjianSiswaRepository interface {
	GetUjiansSiswa(ctx context.Context, filter query.ListUjianFilter)([]ujian.ListUjian,error)
}
