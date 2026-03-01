package fake_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
)

type FakeBankSoalRepo struct {
	DeleteErr    error
	DeleteCalled bool
	GotDeleteID  bank_soal.ID
}

func (f *FakeBankSoalRepo) GetBankSoal(_ context.Context, _ query.BankSoalFilter) ([]bank_soal.BankSoal, error) {
	panic("not used in this test")
}

func (f *FakeBankSoalRepo) GetBankSoalByGuru(_ context.Context, _ bank_soal.ID) ([]bank_soal.BankSoal, error) {
	panic("not used in this test")
}

func (f *FakeBankSoalRepo) GetBankSoalById(_ context.Context, _ bank_soal.ID) (bank_soal.BankSoal, error) {
	panic("not used in this test")
}

func (f *FakeBankSoalRepo) CreateBankSoal(_ context.Context, _ bank_soal.BankSoal) error {
	panic("not used in this test")
}

func (f *FakeBankSoalRepo) UpdateBankSoal(_ context.Context, _ bank_soal.ID, _ updatepatch.UpdateBankSoalPatch) error {
	panic("not used in this test")
}

func (f *FakeBankSoalRepo) DeleteBankSoal(_ context.Context, idBankSoal bank_soal.ID) error {
	f.DeleteCalled = true
	f.GotDeleteID = idBankSoal
	return f.DeleteErr
}
