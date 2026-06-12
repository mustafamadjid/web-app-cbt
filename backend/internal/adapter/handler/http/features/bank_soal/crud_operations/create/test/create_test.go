package httpx_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	bank_soal_create "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/crud_operations/create"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
	bank_soal_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/create"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCreateBankSoalRepo struct{ mock.Mock }

func (m *mockCreateBankSoalRepo) GetBankSoal(ctx context.Context, f query.BankSoalFilter) ([]bank_soal.BankSoal, error) {
	a := m.Called(ctx, f)
	return a.Get(0).([]bank_soal.BankSoal), a.Error(1)
}
func (m *mockCreateBankSoalRepo) GetBankSoalUploaded(ctx context.Context, f query.BankSoalFilter) ([]bank_soal.BankSoal, error) {
	a := m.Called(ctx, f)
	return a.Get(0).([]bank_soal.BankSoal), a.Error(1)
}
func (m *mockCreateBankSoalRepo) GetBankSoalByGuru(ctx context.Context, id bank_soal.ID) ([]bank_soal.BankSoal, error) {
	a := m.Called(ctx, id)
	return a.Get(0).([]bank_soal.BankSoal), a.Error(1)
}
func (m *mockCreateBankSoalRepo) GetBankSoalById(ctx context.Context, id bank_soal.ID) (bank_soal.BankSoal, error) {
	a := m.Called(ctx, id)
	return a.Get(0).(bank_soal.BankSoal), a.Error(1)
}
func (m *mockCreateBankSoalRepo) CreateBankSoal(ctx context.Context, b bank_soal.BankSoal) error {
	return m.Called(ctx, b).Error(0)
}
func (m *mockCreateBankSoalRepo) UpdateBankSoal(ctx context.Context, id bank_soal.ID, b updatepatch.UpdateBankSoalPatch) error {
	return m.Called(ctx, id, b).Error(0)
}
func (m *mockCreateBankSoalRepo) DeleteBankSoal(ctx context.Context, id bank_soal.ID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockCreateBankSoalRepo) GetIdBankSoalByAttemptId(ctx context.Context, id ujian.ID) (ujian.ID, error) {
	a := m.Called(ctx, id)
	return a.Get(0).(ujian.ID), a.Error(1)
}

func newCreateBankSoalHandler(repo *mockCreateBankSoalRepo) *bank_soal_create.CreateBankSoalHandler {
	return bank_soal_create.NewCreateBankSoalHandler(bank_soal_create_service.NewCreateBankSoalService(repo))
}

func createBankSoalRequest(method, body string, withActor bool) *http.Request {
	req := httptest.NewRequest(method, "/bank-soal", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if withActor {
		req = req.WithContext(middleware.WithActor(req.Context(), user.Actor{IdPengguna: 7, Role: user.GURU}))
	}
	return req
}

func TestCreateBankSoal_Success(t *testing.T) {
	repo := new(mockCreateBankSoalRepo)
	handler := newCreateBankSoalHandler(repo)
	body := `{"id_mapel":2,"id_kelas":3,"nama_bank_soal":" Bank Soal ","deskripsi":" Desc ","materi":" Materi "}`
	req := createBankSoalRequest(http.MethodPost, body, true)
	rec := httptest.NewRecorder()

	repo.On("CreateBankSoal", mock.Anything, mock.MatchedBy(func(b bank_soal.BankSoal) bool {
		return b.IdMapel == 2 && b.IdKelas == 3 && b.IdPengguna == 7 &&
			b.NamaBankSoal == "Bank Soal" && b.Deskripsi == "Desc" && b.Materi == "Materi"
	})).Return(nil).Once()

	handler.CreateBankSoal(rec, req, nil)

	assert.Equal(t, http.StatusOK, rec.Code)
	repo.AssertExpectations(t)
}

func TestCreateBankSoal_BadRequests(t *testing.T) {
	tests := []struct {
		name       string
		req        *http.Request
		wantBody   string
		withActor  bool
		wantStatus int
	}{
		{
			name:       "method not allowed",
			req:        createBankSoalRequest(http.MethodGet, `{}`, true),
			wantBody:   "method not allowed",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name: "missing content type",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/bank-soal", strings.NewReader(`{}`))
				return req.WithContext(middleware.WithActor(req.Context(), user.Actor{IdPengguna: 7, Role: user.GURU}))
			}(),
			wantBody:   "content type must be application/json",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			req:        createBankSoalRequest(http.MethodPost, `{"id_mapel":2`, true),
			wantBody:   "invalid request body",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing nama bank soal",
			req:        createBankSoalRequest(http.MethodPost, `{"id_mapel":2,"id_kelas":3,"nama_bank_soal":" ","deskripsi":"Desc","materi":"Materi"}`, true),
			wantBody:   "nama_bank_soal",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing id kelas",
			req:        createBankSoalRequest(http.MethodPost, `{"id_mapel":2,"id_kelas":0,"nama_bank_soal":"Bank","deskripsi":"Desc","materi":"Materi"}`, true),
			wantBody:   "id kelas is required",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing id mapel",
			req:        createBankSoalRequest(http.MethodPost, `{"id_mapel":0,"id_kelas":3,"nama_bank_soal":"Bank","deskripsi":"Desc","materi":"Materi"}`, true),
			wantBody:   "id mapel is required",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := newCreateBankSoalHandler(new(mockCreateBankSoalRepo))
			rec := httptest.NewRecorder()

			handler.CreateBankSoal(rec, tc.req, nil)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantBody)
		})
	}
}

func TestCreateBankSoal_MissingActor(t *testing.T) {
	handler := newCreateBankSoalHandler(new(mockCreateBankSoalRepo))
	req := createBankSoalRequest(http.MethodPost, `{"id_mapel":2,"id_kelas":3,"nama_bank_soal":"Bank","deskripsi":"Desc","materi":"Materi"}`, false)
	rec := httptest.NewRecorder()

	handler.CreateBankSoal(rec, req, nil)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed get actor from context")
}

func TestCreateBankSoal_ServiceErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{name: "invalid input", err: coreerror.ErrInvalidInput, wantStatus: http.StatusBadRequest, wantBody: "invalid input"},
		{name: "internal error", err: errors.New("db down"), wantStatus: http.StatusInternalServerError, wantBody: "internal server error"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mockCreateBankSoalRepo)
			handler := newCreateBankSoalHandler(repo)
			req := createBankSoalRequest(http.MethodPost, `{"id_mapel":2,"id_kelas":3,"nama_bank_soal":"Bank","deskripsi":"Desc","materi":"Materi"}`, true)
			rec := httptest.NewRecorder()
			repo.On("CreateBankSoal", mock.Anything, mock.Anything).Return(tc.err).Once()

			handler.CreateBankSoal(rec, req, nil)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantBody)
			repo.AssertExpectations(t)
		})
	}
}
