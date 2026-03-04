package ujian_service

import (
	"context"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type UpdateUjianService struct {
	repo ujian_repo.UjianRepository
}

func NewUpdateUjianService(repo ujian_repo.UjianRepository) *UpdateUjianService {
	return &UpdateUjianService{
		repo: repo,
	}
}

func (r *UpdateUjianService) UpdateUjianService(ctx context.Context, id ujian.ID, payload updatepatch.UpdateUjianPatch) error {
	logger := corelog.FromContext(ctx)

	// Nil check dulu
	if err := validateUpdateUjianPatch(payload); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	if err := validateUpdateUjianID(id); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", coreerror.ErrMissingId)
		return err
	}

	if err := validateUpdateUjianPatchID(payload); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", coreerror.ErrMissingId)
		return err
	}

	if err := sanitizeNamaUjianPatch(&payload); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	if err := sanitizeDeskripsiUjianPatch(&payload); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	if err := r.repo.UpdateUjian(ctx, id, payload); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	return nil
}

func (r *UpdateUjianService) UpdateJadwalUjianService(ctx context.Context, id ujian.ID, payload updatepatch.UpdateJadwalUjianPatch) error {
	logger := corelog.FromContext(ctx)

	// Nil check dulu
	if err := validateUpdateJadwalUjianPatch(payload); err != nil {
		logger.Error(ctx, "failed update jadwal ujian", "layer", "core.service", "op", "ujian.update_jadwal", "err", err)
		return err
	}

	if err := validateUpdateUjianID(id); err != nil {
		logger.Error(ctx, "failed update jadwal ujian", "layer", "core.service", "op", "ujian.update_jadwal", "err", coreerror.ErrMissingId)
		return err
	}

	if err := validateUpdateJadwalUjianPatchID(payload); err != nil {
		logger.Error(ctx, "failed update jadwal ujian", "layer", "core.service", "op", "ujian.update_jadwal", "err", coreerror.ErrMissingId)
		return err
	}

	if err := sanitizeStatusUjianPatch(&payload); err != nil {
		logger.Error(ctx, "failed update jadwal ujian", "layer", "core.service", "op", "ujian.update_jadwal", "err", err)
		return err
	}

	if err := validateUpdateJadwalUjianStatus(payload); err != nil {
		logger.Error(ctx, "failed update jadwal ujian", "layer", "core.service", "op", "ujian.update_jadwal", "err", err)
		return err
	}

	if err := validateUpdateJadwalUjianTime(payload); err != nil {
		logger.Error(ctx, "failed update jadwal ujian", "layer", "core.service", "op", "ujian.update_jadwal", "err", err)
		return err
	}

	if err := r.repo.UpdateJadwalUjian(ctx, id, payload); err != nil {
		logger.Error(ctx, "failed update jadwal ujian", "layer", "core.service", "op", "ujian.update_jadwal", "err", err)
		return err
	}

	return nil
}

func (r *UpdateUjianService) UpdatePesertaUjianService(ctx context.Context, id ujian.ID, payload updatepatch.UpdatePesertaUjianPatch) error {
	logger := corelog.FromContext(ctx)

	// Nil check dulu
	if err := validateUpdatePesertaUjianPatch(payload); err != nil {
		logger.Error(ctx, "failed update peserta ujian", "layer", "core.service", "op", "ujian.update_peserta", "err", err)
		return err
	}

	if err := validateUpdateUjianID(id); err != nil {
		logger.Error(ctx, "failed update peserta ujian", "layer", "core.service", "op", "ujian.update_peserta", "err", coreerror.ErrMissingId)
		return err
	}

	if err := validateUpdatePesertaUjianPatchID(payload); err != nil {
		logger.Error(ctx, "failed update peserta ujian", "layer", "core.service", "op", "ujian.update_peserta", "err", coreerror.ErrMissingId)
		return err
	}

	if err := validateUpdatePesertaUjianTime(payload); err != nil {
		logger.Error(ctx, "failed update peserta ujian", "layer", "core.service", "op", "ujian.update_peserta", "err", err)
		return err
	}

	if err := validateUpdatePesertaUjianNilai(payload); err != nil {
		logger.Error(ctx, "failed update peserta ujian", "layer", "core.service", "op", "ujian.update_peserta", "err", err)
		return err
	}

	if err := r.repo.UpdatePesertaUjian(ctx, id, payload); err != nil {
		logger.Error(ctx, "failed update peserta ujian", "layer", "core.service", "op", "ujian.update_peserta", "err", err)
		return err
	}

	return nil
}

func (r *UpdateUjianService) UpdateJawabanUjianSiswaService(ctx context.Context, id ujian.ID, payload updatepatch.UpdateJawabanUjianSiswaPatch) error {
	logger := corelog.FromContext(ctx)

	// Nil check dulu
	if err := validateUpdateJawabanUjianPatch(payload); err != nil {
		logger.Error(ctx, "failed update jawaban ujian siswa", "layer", "core.service", "op", "ujian.update_jawaban", "err", err)
		return err
	}

	if err := validateUpdateUjianID(id); err != nil {
		logger.Error(ctx, "failed update jawaban ujian siswa", "layer", "core.service", "op", "ujian.update_jawaban", "err", coreerror.ErrMissingId)
		return err
	}

	if err := validateUpdateJawabanUjianPatchID(payload); err != nil {
		logger.Error(ctx, "failed update jawaban ujian siswa", "layer", "core.service", "op", "ujian.update_jawaban", "err", coreerror.ErrMissingId)
		return err
	}

	if err := sanitizeJawabanEssayPatch(&payload); err != nil {
		logger.Error(ctx, "failed update jawaban ujian siswa", "layer", "core.service", "op", "ujian.update_jawaban", "err", err)
		return err
	}

	if err := validateUpdateJawabanUjianTime(payload); err != nil {
		logger.Error(ctx, "failed update jawaban ujian siswa", "layer", "core.service", "op", "ujian.update_jawaban", "err", err)
		return err
	}

	if err := validateUpdateJawabanUjianPair(payload); err != nil {
		logger.Error(ctx, "failed update jawaban ujian siswa", "layer", "core.service", "op", "ujian.update_jawaban", "err", err)
		return err
	}

	if err := r.repo.UpdateJawabanUjianSiswa(ctx, id, payload); err != nil {
		logger.Error(ctx, "failed update jawaban ujian siswa", "layer", "core.service", "op", "ujian.update_jawaban", "err", err)
		return err
	}

	return nil
}
