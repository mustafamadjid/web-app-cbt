package bank_soal_repo

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
)

type BankSoalRepository interface {
	GetBankSoal(ctx context.Context)

	CreateBankSoal(ctx context.Context, bankSoal bank_soal.BankSoal) error
 	
}