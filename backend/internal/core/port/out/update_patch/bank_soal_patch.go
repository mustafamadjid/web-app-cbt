package updatepatch

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"

type UpdateBankSoalPatch struct {
	IdMapel   *bank_soal.ID
	IdKelas   *bank_soal.ID
	IdPengguna *bank_soal.ID
	NamaBankSoal  *string
	Deskripsi *string
	Materi *string
}