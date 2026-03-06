package httpx

type BankSoalResponse struct {
	IDBankSoal    int    `json:"id_bank_soal"`
	IDMapel       int    `json:"id_mapel"`
	IDKelas       int    `json:"id_kelas"`
	IDPengguna    int    `json:"id_pengguna"`

	Mapel         string `json:"mapel"`
	GuruPembuat   string `json:"guru_pembuat"`
	Kelas         string `json:"kelas"`
	NamaBankSoal  string `json:"nama_bank_soal"`
	Deskripsi     string `json:"deskripsi"`
	Materi        string `json:"materi"`
	TanggalDibuat string `json:"tanggal_dibuat"`
	SoalUploaded  bool   `json:"soal_uploaded"`
}
