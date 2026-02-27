package bank_soal_service

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	bank_soal_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/bank_soal"
)

type UpdateBankSoalService struct {
	repo bank_soal_repo.BankSoalRepository
}

func NewUpdateBankSoalService(repo bank_soal_repo.BankSoalRepository) *UpdateBankSoalService {
	return &UpdateBankSoalService{
		repo: repo,
	}
}

func(r *UpdateBankSoalService) UpdateBankSoalService(ctx context.Context, idBankSoal bank_soal.ID, payload bank_soal.BankSoal) (bank_soal.BankSoal, error){
	
}