package ujian_service

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

type createUjianRepository interface {
	CreateUjian(ctx context.Context, ujian ujian.PenjadwalanUjian) error
}

type createPesertaUjianRepository interface {
	CreatePesertaUjian(ctx context.Context, peserta ujian.PesertaUjian) (ujian.ID, error)
}

type createJawabanUjianRepository interface {
	CreateJawabanUjianSiswa(ctx context.Context, jawaban ujian.JawabanUjianSiswa) (ujian.ID, error)
}

type UjianRepository struct {
	repo createUjianRepository
}

func NewUjianRepository(repo createUjianRepository) UjianRepository {
	return UjianRepository{repo: repo}
}

func (r UjianRepository) CreateUjian(ctx context.Context, data ujian.PenjadwalanUjian) error {
	return r.repo.CreateUjian(ctx, data)
}

type PesertaUjianRepository struct {
	repo createPesertaUjianRepository
}

func NewPesertaUjianRepository(repo createPesertaUjianRepository) PesertaUjianRepository {
	return PesertaUjianRepository{repo: repo}
}

func (r PesertaUjianRepository) CreatePesertaUjian(ctx context.Context, data ujian.PesertaUjian) (ujian.ID, error) {
	return r.repo.CreatePesertaUjian(ctx, data)
}

type JawabanUjianRepository struct {
	repo createJawabanUjianRepository
}

func NewJawabanUjianRepository(repo createJawabanUjianRepository) JawabanUjianRepository {
	return JawabanUjianRepository{repo: repo}
}

func (r JawabanUjianRepository) CreateJawabanUjianSiswa(ctx context.Context, data ujian.JawabanUjianSiswa) (ujian.ID, error) {
	return r.repo.CreateJawabanUjianSiswa(ctx, data)
}
