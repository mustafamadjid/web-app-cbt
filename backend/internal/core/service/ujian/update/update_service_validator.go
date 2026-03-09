package ujian_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

func validateUpdateUjianID(id ujian.ID) error {
	if id <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}
func validateUpdatePenjadwalanUjianPatch(payload updatepatch.UpdatePenjadwalanUjian) error {
	if err := validateUpdateUjianPatch(payload.Ujian); err == nil {
		return nil
	}
	if err := validateUpdateJadwalUjianPatch(payload.JadwalUjian); err == nil {
		return nil
	}
	return coreerror.ErrNoFieldToUpdate
}
func validateUpdateUjianPatch(payload updatepatch.UpdateUjianPatch) error {
	if payload.IdBankSoal == nil && payload.IdKelas == nil && payload.IdNamaKelas == nil && payload.IdGuru == nil && payload.NamaUjian == nil && payload.DeskripsiUjian == nil && payload.AcakSoal == nil {
		return coreerror.ErrNoFieldToUpdate
	}
	return nil
}
func validateUpdateUjianPatchID(payload updatepatch.UpdateUjianPatch) error {
	if payload.IdBankSoal != nil && *payload.IdBankSoal <= 0 {
		return coreerror.ErrMissingId
	}
	if payload.IdKelas != nil && *payload.IdKelas <= 0 {
		return coreerror.ErrMissingId
	}
	if payload.IdNamaKelas != nil && *payload.IdNamaKelas < 0 {
		return coreerror.ErrMissingId
	}
	if payload.IdGuru != nil && *payload.IdGuru <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}
func validateUpdateJadwalUjianPatch(payload updatepatch.UpdateJadwalUjianPatch) error {
	if payload.IdUjian == nil &&
		payload.IdSesi == nil &&
		payload.IdRuangan == nil &&
		payload.IdPengawas == nil &&
		payload.TanggalUjian == nil &&
		payload.Token == nil &&
		payload.WaktuMulai == nil &&
		payload.WaktuSelesai == nil &&
		payload.StatusUjian == nil {
		return coreerror.ErrNoFieldToUpdate
	}
	return nil
}
func validateUpdateJadwalUjianPatchID(payload updatepatch.UpdateJadwalUjianPatch) error {
	if payload.IdUjian != nil && *payload.IdUjian <= 0 {
		return coreerror.ErrMissingId
	}
	if payload.IdSesi != nil && *payload.IdSesi <= 0 {
		return coreerror.ErrMissingId
	}
	if payload.IdRuangan != nil && *payload.IdRuangan <= 0 {
		return coreerror.ErrMissingId
	}
	if payload.IdPengawas != nil && *payload.IdPengawas <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}
func validateUpdateJadwalUjianStatus(payload updatepatch.UpdateJadwalUjianPatch) error {
	if payload.StatusUjian == nil {
		return nil
	}
	if !payload.StatusUjian.ValidStatus() {
		return coreerror.ErrInvalidInput
	}
	return nil
}
func validateUpdateJadwalUjianTime(payload updatepatch.UpdateJadwalUjianPatch) error {
	if payload.TanggalUjian != nil && payload.TanggalUjian.IsZero() {
		return coreerror.ErrInvalidInput
	}
	if payload.WaktuMulai != nil && payload.WaktuMulai.IsZero() {
		return coreerror.ErrInvalidInput
	}
	if payload.WaktuSelesai != nil && payload.WaktuSelesai.IsZero() {
		return coreerror.ErrInvalidInput
	}
	if payload.WaktuMulai != nil && payload.WaktuSelesai != nil && !payload.WaktuMulai.Before(*payload.WaktuSelesai) {
		return coreerror.ErrInvalidInput
	}
	return nil
}
func validateUpdateJadwalUjianToken(payload updatepatch.UpdateJadwalUjianPatch) error {
	if payload.Token == nil {
		return nil
	}
	if len(*payload.Token) > 100 {
		return coreerror.ErrInvalidInput
	}
	return nil
}
