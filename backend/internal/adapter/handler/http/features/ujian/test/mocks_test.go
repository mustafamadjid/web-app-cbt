package ujian_test

import (
	"context"
	"time"

	ujiandomain "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	ujianquery "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
	"github.com/stretchr/testify/mock"
)

type MockUjianRepository struct {
	mock.Mock
}

func (m *MockUjianRepository) GetIdUjianByAttempt(ctx context.Context, idAttempt ujiandomain.ID) (ujiandomain.ID, error) {
	args := m.Called(ctx, idAttempt)
	return args.Get(0).(ujiandomain.ID), args.Error(1)
}

func (m *MockUjianRepository) CreateUjian(ctx context.Context, u ujiandomain.PenjadwalanUjian) error {
	return m.Called(ctx, u).Error(0)
}

func (m *MockUjianRepository) UpdateUjian(ctx context.Context, id ujiandomain.ID, payload updatepatch.UpdatePenjadwalanUjian) error {
	return m.Called(ctx, id, payload).Error(0)
}

func (m *MockUjianRepository) DeleteUjian(ctx context.Context, id ujiandomain.ID) error {
	return m.Called(ctx, id).Error(0)
}

type MockAttemptUjianRepository struct {
	mock.Mock
}

func (m *MockAttemptUjianRepository) GetAttemptUjianById(ctx context.Context, idAttempt ujiandomain.ID) (ujiandomain.AttemptUjian, error) {
	args := m.Called(ctx, idAttempt)
	return args.Get(0).(ujiandomain.AttemptUjian), args.Error(1)
}

func (m *MockAttemptUjianRepository) CreateAttemptUjian(ctx context.Context, data ujiandomain.AttemptUjian) error {
	return m.Called(ctx, data).Error(0)
}

func (m *MockAttemptUjianRepository) UpdateAttemptUjian(ctx context.Context, idAttempt ujiandomain.ID, data updatepatch.UpdateAttemptUjianPatch) error {
	return m.Called(ctx, idAttempt, data).Error(0)
}

func (m *MockAttemptUjianRepository) DeleteAttemptUjian(ctx context.Context, idAttempt ujiandomain.ID) error {
	return m.Called(ctx, idAttempt).Error(0)
}

func (m *MockAttemptUjianRepository) SubmitAttemptUjian(ctx context.Context, idAttempt ujiandomain.ID) error {
	return m.Called(ctx, idAttempt).Error(0)
}

func (m *MockAttemptUjianRepository) ListPesertaUjianAttemptSubmittedByIdJadwalUjian(ctx context.Context, idJadwalUjian ujiandomain.ID) ([]ujiandomain.PesertaUjianSubmitted, error) {
	args := m.Called(ctx, idJadwalUjian)
	return args.Get(0).([]ujiandomain.PesertaUjianSubmitted), args.Error(1)
}

type MockJawabanUjianRepository struct {
	mock.Mock
}

func (m *MockJawabanUjianRepository) GetJawabanUjianByAttemptId(ctx context.Context, idAttempt ujiandomain.ID) ([]ujiandomain.JawabanUjian, error) {
	args := m.Called(ctx, idAttempt)
	return args.Get(0).([]ujiandomain.JawabanUjian), args.Error(1)
}

func (m *MockJawabanUjianRepository) SaveJawabanUjian(ctx context.Context, idAttempt ujiandomain.ID, jawaban []ujiandomain.JawabanUjian) error {
	return m.Called(ctx, idAttempt, jawaban).Error(0)
}

func (m *MockJawabanUjianRepository) ListHasilJawabanUjianByIdAttempt(ctx context.Context, idAttempt ujiandomain.ID) ([]ujiandomain.HasilJawabanUjian, error) {
	args := m.Called(ctx, idAttempt)
	return args.Get(0).([]ujiandomain.HasilJawabanUjian), args.Error(1)
}

type MockSiswaUjianChecker struct {
	mock.Mock
}

func (m *MockSiswaUjianChecker) CheckAttemptOwnershipBySiswa(ctx context.Context, idSiswa int, idAttempt ujiandomain.ID) (bool, error) {
	args := m.Called(ctx, idSiswa, idAttempt)
	return args.Bool(0), args.Error(1)
}

func (m *MockSiswaUjianChecker) CheckValidSiswaInPesertaUjianById(ctx context.Context, idSiswa int, idJadwalUjian int) (bool, int, error) {
	args := m.Called(ctx, idSiswa, idJadwalUjian)
	return args.Bool(0), args.Int(1), args.Error(2)
}

func (m *MockSiswaUjianChecker) CheckTokenUjian(ctx context.Context, token string, idJadwalUjian int) (bool, error) {
	args := m.Called(ctx, token, idJadwalUjian)
	return args.Bool(0), args.Error(1)
}

func (m *MockSiswaUjianChecker) GetDeadlineUjian(ctx context.Context, idJadwalUjian int) (time.Time, error) {
	args := m.Called(ctx, idJadwalUjian)
	return args.Get(0).(time.Time), args.Error(1)
}

type MockUjianSiswaRepository struct {
	mock.Mock
}

func (m *MockUjianSiswaRepository) ListUjianSiswa(ctx context.Context, idSiswa int, filter ujianquery.ListUjianFilter) ([]ujiandomain.ListUjian, error) {
	args := m.Called(ctx, idSiswa, filter)
	return args.Get(0).([]ujiandomain.ListUjian), args.Error(1)
}

func (m *MockUjianSiswaRepository) GetWaktuSelesaiUjian(ctx context.Context, idJadwalUjian int) (time.Time, error) {
	args := m.Called(ctx, idJadwalUjian)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockUjianSiswaRepository) GetActiveUjianAttemptBySiswa(ctx context.Context, idSiswa int, idJadwalUjian int) (ujiandomain.AttemptUjian, error) {
	args := m.Called(ctx, idSiswa, idJadwalUjian)
	return args.Get(0).(ujiandomain.AttemptUjian), args.Error(1)
}
