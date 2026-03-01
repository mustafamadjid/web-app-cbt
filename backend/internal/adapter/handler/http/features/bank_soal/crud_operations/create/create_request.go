package httpx

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"

type CreateBankSoalRequest struct {
	IdPengguna   bank_soal.ID `json:"id_pengguna"`
	IdMapel      bank_soal.ID `json:"id_mapel"`
	IdKelas      bank_soal.ID `json:"id_kelas"`
	NamaBankSoal string       `json:"nama_bank_soal"`
	Deskripsi    string       `json:"deskripsi"`
	Materi       string       `json:"materi"`
}
