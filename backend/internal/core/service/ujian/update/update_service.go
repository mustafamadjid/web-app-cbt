package ujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type UpdateUjianService struct {
	ujianRepo UjianRepository
}

func NewUpdateUjianService(ujianRepo UjianRepository) *UpdateUjianService {
	return &UpdateUjianService{
		ujianRepo: ujianRepo,
	}
}

func (r *UpdateUjianService) UpdateUjianService(ctx context.Context, id ujian.ID, payload updatepatch.UpdatePenjadwalanUjian) error {
	logger := corelog.FromContext(ctx)

	if err := validateUpdatePenjadwalanUjianPatch(payload); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	if err := validateUpdateUjianID(id); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", coreerror.ErrMissingId)
		return err
	}

	if err := validateUpdateUjianPatchID(payload.Ujian); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", coreerror.ErrMissingId)
		return err
	}

	if err := validateUpdateJadwalUjianPatchID(payload.JadwalUjian); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", coreerror.ErrMissingId)
		return err
	}

	if err := sanitizeNamaUjianPatch(&payload.Ujian); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	if err := sanitizeDeskripsiUjianPatch(&payload.Ujian); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	if err := sanitizeStatusUjianPatch(&payload.JadwalUjian); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	if err := sanitizeTokenJadwalUjianPatch(&payload.JadwalUjian); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	if err := validateUpdateJadwalUjianStatus(payload.JadwalUjian); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	if err := validateUpdateJadwalUjianToken(payload.JadwalUjian); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	if err := validateUpdateJadwalUjianTime(payload.JadwalUjian); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	if err := r.ujianRepo.UpdateUjian(ctx, id, payload); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	return nil
}

type UpdatePesertaUjianService struct {
	pesertaRepo PesertaUjianRepository
}

func NewUpdatePesertaUjianService(pesertaRepo PesertaUjianRepository) *UpdatePesertaUjianService {
	return &UpdatePesertaUjianService{
		pesertaRepo: pesertaRepo,
	}
}

func (r *UpdatePesertaUjianService) UpdatePesertaUjianService(ctx context.Context, id ujian.ID, payload updatepatch.UpdatePesertaUjianPatch) error {
	logger := corelog.FromContext(ctx)

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

	if err := r.pesertaRepo.UpdatePesertaUjian(ctx, id, payload); err != nil {
		logger.Error(ctx, "failed update peserta ujian", "layer", "core.service", "op", "ujian.update_peserta", "err", err)
		return err
	}

	return nil
}

type UpdateJawabanUjianSiswaService struct {
	jawabanRepo JawabanUjianRepository
}

func NewUpdateJawabanUjianSiswaService(jawabanRepo JawabanUjianRepository) *UpdateJawabanUjianSiswaService {
	return &UpdateJawabanUjianSiswaService{
		jawabanRepo: jawabanRepo,
	}
}

func (r *UpdateJawabanUjianSiswaService) UpdateJawabanUjianSiswaService(ctx context.Context, id ujian.ID, payload updatepatch.UpdateJawabanUjianSiswaPatch) error {
	logger := corelog.FromContext(ctx)

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

	if err := r.jawabanRepo.UpdateJawabanUjianSiswa(ctx, id, payload); err != nil {
		logger.Error(ctx, "failed update jawaban ujian siswa", "layer", "core.service", "op", "ujian.update_jawaban", "err", err)
		return err
	}

	return nil
}
