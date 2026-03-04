package ujian_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func validateCreateUjian(data ujian.PenjadwalanUjian) error {
	if data.Ujian.IdBankSoal <= 0 || data.Ujian.IdKelas <= 0 || data.Ujian.IdGuru <= 0 {
		return coreerror.ErrMissingId
	}
	if data.Ujian.IdNamaKelas != nil && *data.Ujian.IdNamaKelas <= 0 {
		return coreerror.ErrMissingId
	}
	if data.Ujian.NamaUjian == "" {
		return coreerror.ErrMissingField
	}
	if len(data.Ujian.NamaUjian) > 100 {
		return coreerror.ErrInvalidInput
	}
	if data.JadwalUjian.IdSesi <= 0 || data.JadwalUjian.IdRuangan <= 0 || data.JadwalUjian.IdPengawas <= 0 {
		return coreerror.ErrMissingId
	}
	if data.JadwalUjian.TanggalUjian.IsZero() || data.JadwalUjian.WaktuMulai.IsZero() || data.JadwalUjian.WaktuSelesai.IsZero() {
		return coreerror.ErrInvalidInput
	}
	if !data.JadwalUjian.WaktuMulai.Before(data.JadwalUjian.WaktuSelesai) {
		return coreerror.ErrInvalidInput
	}
	if data.JadwalUjian.Token == "" {
		return coreerror.ErrMissingField
	}
	if len(data.JadwalUjian.Token) > 100 {
		return coreerror.ErrInvalidInput
	}
	if data.JadwalUjian.StatusUjian == "" {
		return coreerror.ErrMissingField
	}
	if !data.JadwalUjian.StatusUjian.ValidStatus() {
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
