package ujian_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func validateCreateUjian(data ujian.Ujian) error {
	if data.IdBankSoal <= 0 || data.IdKelas <= 0 || data.IdGuru <= 0 {
		return coreerror.ErrMissingId
	}
	if data.IdNamaKelas != nil && *data.IdNamaKelas <= 0 {
		return coreerror.ErrMissingId
	}
	if data.NamaUjian == "" {
		return coreerror.ErrMissingField
	}
	if len(data.NamaUjian) > 100 {
		return coreerror.ErrInvalidInput
	}
	return nil
}
func validateCreateJadwalUjian(data ujian.JadwalUjian) error {
	if data.IdUjian <= 0 || data.IdSesi <= 0 || data.IdRuangan <= 0 {
		return coreerror.ErrMissingId
	}
	if data.TanggalUjian.IsZero() || data.WaktuMulai.IsZero() || data.WaktuSelesai.IsZero() {
		return coreerror.ErrMissingField
	}
	if !data.WaktuMulai.Before(data.WaktuSelesai) {
		return coreerror.ErrInvalidInput
	}
	if !data.StatusUjian.ValidStatus() {
		return coreerror.ErrInvalidInput
	}
	return nil
}
func validateCreatePesertaUjian(data ujian.PesertaUjian) error {
	if data.IdJadwalUjian <= 0 || data.IdSiswa <= 0 {
		return coreerror.ErrMissingId
	}
	if data.WaktuMulai != nil && data.WaktuMulai.IsZero() {
		return coreerror.ErrInvalidInput
	}
	if data.WaktuSubmit != nil && data.WaktuSubmit.IsZero() {
		return coreerror.ErrInvalidInput
	}
	if data.WaktuMulai != nil && data.WaktuSubmit != nil && data.WaktuSubmit.Before(*data.WaktuMulai) {
		return coreerror.ErrInvalidInput
	}
	if data.NilaiUjian != nil && (*data.NilaiUjian < 0 || *data.NilaiUjian > 100) {
		return coreerror.ErrInvalidInput
	}
	return nil
}
func validateCreateJawabanUjianSiswa(data ujian.JawabanUjianSiswa) error {
	if data.IdPesertaUjian <= 0 || data.IdSoal <= 0 {
		return coreerror.ErrMissingId
	}
	if data.IdPilihan != nil && *data.IdPilihan <= 0 {
		return coreerror.ErrMissingId
	}
	if data.IdPilihan != nil && data.JawabanEssay != nil {
		return coreerror.ErrInvalidInput
	}
	if data.WaktuJawab != nil && data.WaktuJawab.IsZero() {
		return coreerror.ErrInvalidInput
	}
	return nil
}
