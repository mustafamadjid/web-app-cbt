package siswaujian_service

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
)

type ListUjianSiswaService struct {
	repo ujian_repo.ListUjianSiswaRepository
}

func NewListUjianSiswaService(repo ujian_repo.ListUjianSiswaRepository) *ListUjianSiswaService {
	return &ListUjianSiswaService{
		repo: repo,
	}
}

func (s *ListUjianSiswaService) ListUjianSiswa(ctx context.Context,idSiswa int, filter query.ListUjianFilter) ([]ujian.ListUjian, error) {
	
}
