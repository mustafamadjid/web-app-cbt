package updatepatch

import (
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type ProfilSiswa struct {
	IdTingkatKelas *user.ID
	IdNamaKelas    *user.ID
	Nisn           *string
	NoAbsen        *int
	Angkatan       *int
	TempatLahir    *string
	TanggalLahir   *time.Time
}
