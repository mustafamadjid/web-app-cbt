package gradingrepo

import ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"

func buildStatistikSoalPayload(items []ujian.StatistikSoal, isBenar bool) []statistikSoalUpsertItem {
	if len(items) == 0 {
		return nil
	}

	payload := make([]statistikSoalUpsertItem, 0, len(items))

	for _, item := range items {
		entry := statistikSoalUpsertItem{
			IDSoal:  item.IDSoal,
			IDUjian: item.IDUjian,
		}

		if isBenar {
			entry.JumlahJawabanBenar++
		} else {
			entry.JumlahJawabanSalah++
		}

		payload = append(payload, entry)
	}

	return payload
}
