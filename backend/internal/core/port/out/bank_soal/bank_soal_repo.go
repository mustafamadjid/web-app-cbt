package bank_soal_repo

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
)

type BankSoalRepository interface {
	GetBankSoal(ctx context.Context, filter query.BankSoalFilter) ([]bank_soal.BankSoal, error)
	GetBankSoalUploaded(ctx context.Context, filter query.BankSoalFilter) ([]bank_soal.BankSoal, error)
	GetBankSoalByGuru(ctx context.Context, idPengguna bank_soal.ID) ([]bank_soal.BankSoal, error)
	GetBankSoalById(ctx context.Context, idBankSoal bank_soal.ID) (bank_soal.BankSoal, error)

	CreateBankSoal(ctx context.Context, bankSoal bank_soal.BankSoal) error
	UpdateBankSoal(ctx context.Context, idBankSoal bank_soal.ID, bankSoal updatepatch.UpdateBankSoalPatch) error
	DeleteBankSoal(ctx context.Context, idBankSoal bank_soal.ID) error

	GetIdBankSoalByAttemptId(ctx context.Context,idAttempt int) (int,error)
}
