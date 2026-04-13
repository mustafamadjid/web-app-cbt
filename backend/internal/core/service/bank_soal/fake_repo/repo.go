package fakerepo

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	bank_soal_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/bank_soal"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
)

var _ bank_soal_repo.BankSoalRepository = (*FakeBankSoalRepo)(nil)

type FakeBankSoalRepo struct {
	GetBankSoalData []bank_soal.BankSoal
	GetBankSoalErr  error
	GetCalled       bool
	GotFilter       query.BankSoalFilter

	GetBankSoalUploadedData []bank_soal.BankSoal
	GetBankSoalUploadedErr  error
	GetUploadedCalled       bool
	GotUploadedFilter       query.BankSoalFilter

	GetByIDData   bank_soal.BankSoal
	GetByIDErr    error
	GetByIDCalled bool
	GotID         bank_soal.ID

	GetByGuruData   []bank_soal.BankSoal
	GetByGuruErr    error
	GetByGuruCalled bool
	GotGuruID       bank_soal.ID

	CreateErr    error
	CreateCalled bool
	GotCreate    bank_soal.BankSoal

	UpdateErr    error
	UpdateCalled bool
	GotPatch     updatepatch.UpdateBankSoalPatch

	DeleteErr    error
	DeleteCalled bool
	GotDeleteID  bank_soal.ID

	GetByAttemptData   ujian.ID
	GetByAttemptErr    error
	GetByAttemptCalled bool
	GotAttemptID       ujian.ID
}

func (f *FakeBankSoalRepo) GetBankSoal(_ context.Context, filter query.BankSoalFilter) ([]bank_soal.BankSoal, error) {
	f.GetCalled = true
	f.GotFilter = filter
	return f.GetBankSoalData, f.GetBankSoalErr
}

func (f *FakeBankSoalRepo) GetBankSoalUploaded(_ context.Context, filter query.BankSoalFilter) ([]bank_soal.BankSoal, error) {
	f.GetUploadedCalled = true
	f.GotUploadedFilter = filter
	return f.GetBankSoalUploadedData, f.GetBankSoalUploadedErr
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

func (f *FakeBankSoalRepo) CreateBankSoal(_ context.Context, data bank_soal.BankSoal) error {
	f.CreateCalled = true
	f.GotCreate = data
	return f.CreateErr
}

func (f *FakeBankSoalRepo) UpdateBankSoal(_ context.Context, idBankSoal bank_soal.ID, patch updatepatch.UpdateBankSoalPatch) error {
	f.UpdateCalled = true
	f.GotID = idBankSoal
	f.GotPatch = patch
	return f.UpdateErr
}

func (f *FakeBankSoalRepo) DeleteBankSoal(_ context.Context, idBankSoal bank_soal.ID) error {
	f.DeleteCalled = true
	f.GotDeleteID = idBankSoal
	return f.DeleteErr
}

func (f *FakeBankSoalRepo) GetIdBankSoalByAttemptId(_ context.Context, idAttempt ujian.ID) (ujian.ID, error) {
	f.GetByAttemptCalled = true
	f.GotAttemptID = idAttempt
	return f.GetByAttemptData, f.GetByAttemptErr
}
