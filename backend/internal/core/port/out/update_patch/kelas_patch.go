package updatepatch

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"

type TingkatKelasPatch struct {
	TingkatKelas *int
}

type NamaKelasPatch struct {
	IdTingkatKelas *kelas.ID
	NamaKelas      *string
}