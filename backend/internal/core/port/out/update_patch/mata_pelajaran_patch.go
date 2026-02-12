package updatepatch

import matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"

type UpdateMapelPatch struct {
	IdKelas   *matapelajaran.ID
	KodeMapel *string
	NamaMapel *string
	Deskripsi *string
}
