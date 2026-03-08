package httpx

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"

func toBankSoalIDPointer(v *int) *bank_soal.ID {
	if v == nil {
		return nil
	}

	id := bank_soal.ID(*v)
	return &id
}