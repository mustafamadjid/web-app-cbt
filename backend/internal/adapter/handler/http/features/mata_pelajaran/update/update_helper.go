package httpx

import matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"

func toMapelIDPointer(v *int) *matapelajaran.ID {
	if v == nil {
		return nil
	}

	id := matapelajaran.ID(*v)
	return &id
}
