package kelas_service

import (
	"context"

	kelas_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/kelas"
)

type DeleteKelasService struct {
	deleteRepo kelas_repo.KelasRepository
}

func NewDeleteKelasService(kelasRepo kelas_repo.KelasRepository) *DeleteKelasService {
	return &DeleteKelasService{deleteRepo: kelasRepo}
}

func(s *DeleteKelasService)DeleteNamaKelas(ctx context.Context, idNamaKelas int) error{
	return s.deleteRepo.DeleteNamaKelas(ctx, idNamaKelas)
}