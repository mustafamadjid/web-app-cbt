package bank_soal_service

import (
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	"strings"
)

func sanitized(bankSoal bank_soal.BankSoal) bank_soal.BankSoal {
	bankSoal.NamaBankSoal = strings.TrimSpace(bankSoal.NamaBankSoal)
	bankSoal.Deskripsi = strings.TrimSpace(bankSoal.Deskripsi)
	bankSoal.Materi = strings.TrimSpace(bankSoal.Materi)
	return bankSoal
}
