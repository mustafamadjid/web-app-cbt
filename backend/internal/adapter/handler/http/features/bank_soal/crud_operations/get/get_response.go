package httpx

type BankSoalResponse struct {
	IDBankSoal   int    `json:"id_bank_soal"`
	IDMapel      int    `json:"id_mapel"`
	IDKelas      int    `json:"id_kelas"`
	IDPengguna   int    `json:"id_pengguna"`
	NamaBankSoal string `json:"nama_bank_soal"`
	Deskripsi    string `json:"deskripsi"`
	Materi       string `json:"materi"`
}
