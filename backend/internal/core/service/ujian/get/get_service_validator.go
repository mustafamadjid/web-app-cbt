package ujian_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"strings"
)

func validateUjianID(id ujian.ID) error {
	if id <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}
func validatePesertaFilter(filter ujian.PesertaUjian) error {
	if filter.IdPesertaUjian < 0 || filter.IdJadwalUjian < 0 || filter.IdSiswa < 0 {
		return errInvalidPeserta
	}
	return nil
}

func validateJadwalUjian(item ujian.JadwalUjian) error {
	if item.IdJadwalUjian <= 0 ||
		item.IdUjian <= 0 ||
		item.IdSesi <= 0 ||
		item.IdRuangan <= 0 ||
		item.IdPengawas <= 0 {
		return errInvalidJadwalUjian
	}
	if strings.TrimSpace(item.Token) == "" {
		return errInvalidTokenUjian
	}
	if len(item.Token) > 100 {
		return errInvalidTokenUjian
	}
	if item.TanggalUjian.IsZero() || item.WaktuMulai.IsZero() || item.WaktuSelesai.IsZero() {
		return errInvalidTanggalUjian
	}
	if !item.WaktuMulai.Before(item.WaktuSelesai) {
		return errInvalidWaktuUjian
	}
	if !item.StatusUjian.ValidStatus() {
		return errInvalidStatusUjian
	}
	return nil
}

func validateListUjian(item ujian.ListUjian) error {
	if item.IdUjian <= 0 ||
		item.IdBankSoal <= 0 ||
		item.IdGuru <= 0 ||
		item.IdKelas <= 0 ||
		item.IdJadwalUjian <= 0 ||
		item.IdPengawas <= 0 ||
		item.IdSesi <= 0 ||
		item.IdRuangan <= 0 {
		return errInvalidListUjian
	}

	if item.IdNamaKelas != nil && *item.IdNamaKelas <= 0 {
		return errInvalidListUjian
	}

	if item.TingkatKelas <= 0 {
		return errInvalidTingkatKelas
	}

	if strings.TrimSpace(item.NamaUjian) == "" ||
		strings.TrimSpace(item.PembuatUsername) == "" ||
		strings.TrimSpace(item.NamaPengawas) == "" ||
		strings.TrimSpace(item.PengawasUsername) == "" ||
		strings.TrimSpace(item.NamaSesi) == "" ||
		strings.TrimSpace(item.NamaRuangan) == "" {
		return errInvalidNamaUjian
	}

	if item.TanggalUjian.IsZero() || item.WaktuMulai.IsZero() || item.WaktuSelesai.IsZero() {
		return errInvalidTanggalUjian
	}

	if item.WaktuSelesai.Before(item.WaktuMulai) {
		return errInvalidWaktuUjian
	}

	if !item.StatusUjian.ValidStatus() {
		return errInvalidStatusUjian
	}

	return nil
}
