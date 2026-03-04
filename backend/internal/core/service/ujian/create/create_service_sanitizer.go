package ujian_service

import (
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"strings"
)

func sanitizeCreateUjian(data ujian.Ujian) ujian.Ujian {
	data.NamaUjian = strings.TrimSpace(data.NamaUjian)
	if data.DeskripsiUjian != nil {
		deskripsi := strings.TrimSpace(*data.DeskripsiUjian)
		if deskripsi == "" {
			data.DeskripsiUjian = nil
		} else {
			data.DeskripsiUjian = &deskripsi
		}
	}
	return data
}
func sanitizeCreateJadwalUjian(data ujian.JadwalUjian) ujian.JadwalUjian {
	if data.StatusUjian == "" {
		data.StatusUjian = ujian.BELUM_MULAI
	}
	return data
}
func sanitizeCreateJawabanUjianSiswa(data ujian.JawabanUjianSiswa) ujian.JawabanUjianSiswa {
	if data.JawabanEssay != nil {
		jawabanEssay := strings.TrimSpace(*data.JawabanEssay)
		if jawabanEssay == "" {
			data.JawabanEssay = nil
		} else {
			data.JawabanEssay = &jawabanEssay
		}
	}
	return data
}
