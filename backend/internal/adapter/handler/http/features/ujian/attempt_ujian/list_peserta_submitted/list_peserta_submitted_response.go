package httpx

type ListPesertaUjianSubmittedResponse struct {
	IDPesertaUjian int      `json:"id_peserta_ujian"`
	IDAttempt      int      `json:"id_attempt"`
	IDSiswa        int      `json:"id_siswa"`
	TingkatKelas   int      `json:"tingkat_kelas"`
	NamaKelas      string   `json:"nama_kelas"`
	NamaLengkap    string   `json:"nama_lengkap"`
	NoAbsen        int      `json:"no_absen"`
	NilaiAkhir     *float64 `json:"nilai_akhir"`
	WaktuMulai     *string  `json:"waktu_mulai"`
	WaktuSubmit    *string  `json:"waktu_submit"`
}
