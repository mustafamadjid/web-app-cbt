package bank_soal_service_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
)

type FakeBankSoalRepo struct {
	UpdateErr    error
	UpdateCalled bool
	GotID        bank_soal.ID
	GotPatch     updatepatch.UpdateBankSoalPatch
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

func (f *FakeBankSoalRepo) CreateBankSoal(_ context.Context, _ bank_soal.BankSoal) error {
	panic("not used in this test")
}

func (f *FakeBankSoalRepo) UpdateBankSoal(_ context.Context, idBankSoal bank_soal.ID, patch updatepatch.UpdateBankSoalPatch) error {
	f.UpdateCalled = true
	f.GotID = idBankSoal
	f.GotPatch = patch
	return f.UpdateErr
}

func (f *FakeBankSoalRepo) DeleteBankSoal(_ context.Context, _ bank_soal.ID) error {
	panic("not used in this test")
}

func (f *FakeBankSoalRepo) GetIdBankSoalByAttemptId(_ context.Context, _ ujian.ID) (ujian.ID, error) {
	panic("not used in this test")
}
