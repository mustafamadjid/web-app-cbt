package bank_soal_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
	"github.com/stretchr/testify/mock"
)

type MockBankSoalRepository struct {
	mock.Mock
}

func (m *MockBankSoalRepository) GetBankSoal(ctx context.Context, filter query.BankSoalFilter) ([]bank_soal.BankSoal, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]bank_soal.BankSoal), args.Error(1)
}

func (m *MockBankSoalRepository) GetBankSoalUploaded(ctx context.Context, filter query.BankSoalFilter) ([]bank_soal.BankSoal, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]bank_soal.BankSoal), args.Error(1)
}

func (m *MockBankSoalRepository) GetBankSoalByGuru(ctx context.Context, idPengguna bank_soal.ID) ([]bank_soal.BankSoal, error) {
	args := m.Called(ctx, idPengguna)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]bank_soal.BankSoal), args.Error(1)
}

func (m *MockBankSoalRepository) GetBankSoalById(ctx context.Context, idBankSoal bank_soal.ID) (bank_soal.BankSoal, error) {
	args := m.Called(ctx, idBankSoal)
	return args.Get(0).(bank_soal.BankSoal), args.Error(1)
}

func (m *MockBankSoalRepository) CreateBankSoal(ctx context.Context, bs bank_soal.BankSoal) error {
	return m.Called(ctx, bs).Error(0)
}

func (m *MockBankSoalRepository) UpdateBankSoal(ctx context.Context, id bank_soal.ID, bs updatepatch.UpdateBankSoalPatch) error {
	return m.Called(ctx, id, bs).Error(0)
}

func (m *MockBankSoalRepository) DeleteBankSoal(ctx context.Context, id bank_soal.ID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockBankSoalRepository) GetIdBankSoalByAttemptId(ctx context.Context, id ujian.ID) (ujian.ID, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(ujian.ID), args.Error(1)
}

type MockImportSoalJobRepo struct {
	mock.Mock
}

func (m *MockImportSoalJobRepo) CreateJob(ctx context.Context, job importsoal.ImportSoalJob) (int64, error) {
	args := m.Called(ctx, job)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockImportSoalJobRepo) GetPendingJobs(ctx context.Context, limit int) ([]importsoal.ImportSoalJob, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]importsoal.ImportSoalJob), args.Error(1)
}

func (m *MockImportSoalJobRepo) UpdateJobStatus(ctx context.Context, jobID int64, status importsoal.JobStatus, errorMsg, warningMsg string, totalSoal int) error {
	return m.Called(ctx, jobID, status, errorMsg, warningMsg, totalSoal).Error(0)
}

func (m *MockImportSoalJobRepo) GetJobByID(ctx context.Context, jobID int64) (importsoal.ImportSoalJob, error) {
	args := m.Called(ctx, jobID)
	return args.Get(0).(importsoal.ImportSoalJob), args.Error(1)
}

func (m *MockImportSoalJobRepo) GetJobsByBankSoal(ctx context.Context, bankSoalID int64) ([]importsoal.ImportSoalJob, error) {
	args := m.Called(ctx, bankSoalID)
	return args.Get(0).([]importsoal.ImportSoalJob), args.Error(1)
}
