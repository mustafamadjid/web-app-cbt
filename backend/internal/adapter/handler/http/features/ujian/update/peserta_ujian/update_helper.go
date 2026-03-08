package httpx

import ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"

func toIDPesertaUjianPointer(value *int) *ujian.ID {
	if value == nil {
		return nil
	}

	id := ujian.ID(*value)
	return &id
}
