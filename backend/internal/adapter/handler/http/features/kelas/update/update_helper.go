package httpx

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"

func toKelasIDPointer(v *int) *kelas.ID {
	if v == nil {
		return nil
	}
	id := kelas.ID(*v)
	return &id
}