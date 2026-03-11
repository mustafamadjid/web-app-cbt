package siswaujian_service

import ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"

type ListUjianSiswaService struct {
	repo ujian_repo.ListUjianSiswaRepository
}

func NewListUjianSiswaService(repo ujian_repo.ListUjianSiswaRepository) *ListUjianSiswaService {
	return &ListUjianSiswaService{
		repo: repo,
	}
}

func(s *ListUjianSiswaService) ListUjianSiswa(idSiswa int) ([]ujian_repo.ListUjianSiswa, error) {
	
}