package ujian_service

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type updateUjianRepository interface {
	UpdateUjian(ctx context.Context, id ujian.ID, payload updatepatch.UpdatePenjadwalanUjian) error
}

type updatePesertaUjianRepository interface {
	UpdatePesertaUjian(ctx context.Context, id ujian.ID, payload updatepatch.UpdatePesertaUjianPatch) error
}

type updateJawabanUjianRepository interface {
	UpdateJawabanUjianSiswa(ctx context.Context, id ujian.ID, payload updatepatch.UpdateJawabanUjianSiswaPatch) error
}

type UjianRepository struct {
	repo updateUjianRepository
}

func NewUjianRepository(repo updateUjianRepository) UjianRepository {
	return UjianRepository{repo: repo}
}

func (r UjianRepository) UpdateUjian(ctx context.Context, id ujian.ID, payload updatepatch.UpdatePenjadwalanUjian) error {
	return r.repo.UpdateUjian(ctx, id, payload)
}

type PesertaUjianRepository struct {
	repo updatePesertaUjianRepository
}

func NewPesertaUjianRepository(repo updatePesertaUjianRepository) PesertaUjianRepository {
	return PesertaUjianRepository{repo: repo}
}

func (r PesertaUjianRepository) UpdatePesertaUjian(ctx context.Context, id ujian.ID, payload updatepatch.UpdatePesertaUjianPatch) error {
	return r.repo.UpdatePesertaUjian(ctx, id, payload)
}

type JawabanUjianRepository struct {
	repo updateJawabanUjianRepository
}

func NewJawabanUjianRepository(repo updateJawabanUjianRepository) JawabanUjianRepository {
	return JawabanUjianRepository{repo: repo}
}

func (r JawabanUjianRepository) UpdateJawabanUjianSiswa(ctx context.Context, id ujian.ID, payload updatepatch.UpdateJawabanUjianSiswaPatch) error {
	return r.repo.UpdateJawabanUjianSiswa(ctx, id, payload)
}
