package bank_soal_service_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
)

type FakeBankSoalRepo struct {
	CreateErr    error
	CreateCalled bool
	GotCreate    bank_soal.BankSoal
}

func (f *FakeBankSoalRepo) GetBankSoal(_ context.Context, _ query.BankSoalFilter) ([]bank_soal.BankSoal, error) {
	panic("not used in this test")
}

func (f *FakeBankSoalRepo) GetBankSoalUploaded(_ context.Context, _ query.BankSoalFilter) ([]bank_soal.BankSoal, error) {
	panic("not used in this test")
}

func (f *FakeBankSoalRepo) GetBankSoalByGuru(_ context.Context, _ bank_soal.ID) ([]bank_soal.BankSoal, error) {
	panic("not used in this test")
}

func (f *FakeBankSoalRepo) GetBankSoalById(_ context.Context, _ bank_soal.ID) (bank_soal.BankSoal, error) {
	panic("not used in this test")
}

func (f *FakeBankSoalRepo) CreateBankSoal(_ context.Context, data bank_soal.BankSoal) error {
	f.CreateCalled = true
	f.GotCreate = data
	return f.CreateErr
}

func (f *FakeBankSoalRepo) UpdateBankSoal(_ context.Context, _ bank_soal.ID, _ updatepatch.UpdateBankSoalPatch) error {
	panic("not used in this test")
}

func (f *FakeBankSoalRepo) DeleteBankSoal(_ context.Context, _ bank_soal.ID) error {
	panic("not used in this test")
}

func (f *FakeBankSoalRepo) GetIdBankSoalByAttemptId(_ context.Context, _ ujian.ID) (ujian.ID, error) {
	panic("not used in this test")
}
