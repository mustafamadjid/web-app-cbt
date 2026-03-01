package fake_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
)

type FakeBankSoalRepo struct {
	GetBankSoalData []bank_soal.BankSoal
	GetBankSoalErr  error
	GetCalled       bool
	GotFilter       query.BankSoalFilter

	GetByIDData   bank_soal.BankSoal
	GetByIDErr    error
	GetByIDCalled bool
	GotID         bank_soal.ID

	GetByGuruData   []bank_soal.BankSoal
	GetByGuruErr    error
	GetByGuruCalled bool
	GotGuruID       bank_soal.ID
}

func (f *FakeBankSoalRepo) GetBankSoal(_ context.Context, filter query.BankSoalFilter) ([]bank_soal.BankSoal, error) {
	f.GetCalled = true
	f.GotFilter = filter
	return f.GetBankSoalData, f.GetBankSoalErr
}

func (f *FakeBankSoalRepo) GetBankSoalByGuru(_ context.Context, idPengguna bank_soal.ID) ([]bank_soal.BankSoal, error) {
	f.GetByGuruCalled = true
	f.GotGuruID = idPengguna
	return f.GetByGuruData, f.GetByGuruErr
}

func (f *FakeBankSoalRepo) GetBankSoalById(_ context.Context, idBankSoal bank_soal.ID) (bank_soal.BankSoal, error) {
	f.GetByIDCalled = true
	f.GotID = idBankSoal
	return f.GetByIDData, f.GetByIDErr
}

func (f *FakeBankSoalRepo) CreateBankSoal(_ context.Context, _ bank_soal.BankSoal) error {
	panic("not used in this test")
}

func (f *FakeBankSoalRepo) UpdateBankSoal(_ context.Context, _ bank_soal.ID, _ updatepatch.UpdateBankSoalPatch) error {
	panic("not used in this test")
}

func (f *FakeBankSoalRepo) DeleteBankSoal(_ context.Context, _ bank_soal.ID) error {
	panic("not used in this test")
}
