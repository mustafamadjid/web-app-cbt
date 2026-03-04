package ujian_service

import (
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"strings"
)

func sanitizeCreateUjian(data ujian.PenjadwalanUjian) ujian.PenjadwalanUjian {
	data.Ujian.NamaUjian = strings.TrimSpace(data.Ujian.NamaUjian)
	if data.Ujian.DeskripsiUjian != nil {
		deskripsi := strings.TrimSpace(*data.Ujian.DeskripsiUjian)
		if deskripsi == "" {
			data.Ujian.DeskripsiUjian = nil
		} else {
			data.Ujian.DeskripsiUjian = &deskripsi
		}
	}
	data.JadwalUjian.Token = strings.ToUpper(strings.TrimSpace(data.JadwalUjian.Token))
	data.JadwalUjian.StatusUjian = ujian.StatusUjian(strings.ToUpper(strings.TrimSpace(string(data.JadwalUjian.StatusUjian))))
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
