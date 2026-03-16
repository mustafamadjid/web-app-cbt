package httpx

import ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"

func toIDUjianPointer(value *int) *ujian.ID {
	if value == nil {
		return nil
	}

	id := ujian.ID(*value)
	return &id
}

func toStatusUjianPointer(value *string) *ujian.StatusUjian {
	if value == nil {
		return nil
	}

	status := ujian.StatusUjian(*value)
	return &status
}
