package bank_soal

import "time"

type ID int

type BankSoal struct {
	IdBankSoal    ID
	IdMapel       ID
	IdKelas       ID
	IdPengguna    ID
	NamaBankSoal  string
	Deskripsi     string
	Materi        string
	CreatedAt     time.Time
	TanggalDibuat string
	SoalUploaded  bool
}
