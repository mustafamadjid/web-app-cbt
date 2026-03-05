package ujian_service

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

type deleteUjianRepository interface {
	DeleteUjian(ctx context.Context, id ujian.ID) error
}

type deletePesertaUjianRepository interface {
	DeletePesertaUjian(ctx context.Context, id ujian.ID) error
}

type deleteJawabanUjianRepository interface {
	DeleteJawabanUjianSiswa(ctx context.Context, id ujian.ID) error
}

type UjianRepository struct {
	repo deleteUjianRepository
}

func NewUjianRepository(repo deleteUjianRepository) UjianRepository {
	return UjianRepository{repo: repo}
}

func (r UjianRepository) DeleteUjian(ctx context.Context, id ujian.ID) error {
	return r.repo.DeleteUjian(ctx, id)
}

type PesertaUjianRepository struct {
	repo deletePesertaUjianRepository
}

func NewPesertaUjianRepository(repo deletePesertaUjianRepository) PesertaUjianRepository {
	return PesertaUjianRepository{repo: repo}
}

func (r PesertaUjianRepository) DeletePesertaUjian(ctx context.Context, id ujian.ID) error {
	return r.repo.DeletePesertaUjian(ctx, id)
}

type JawabanUjianRepository struct {
	repo deleteJawabanUjianRepository
}

func NewJawabanUjianRepository(repo deleteJawabanUjianRepository) JawabanUjianRepository {
	return JawabanUjianRepository{repo: repo}
}

func (r JawabanUjianRepository) DeleteJawabanUjianSiswa(ctx context.Context, id ujian.ID) error {
	return r.repo.DeleteJawabanUjianSiswa(ctx, id)
}
