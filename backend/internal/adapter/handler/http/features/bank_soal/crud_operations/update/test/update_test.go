package httpx_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	bank_soal_update "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/crud_operations/update"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
	bank_soal_update_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUpdateBankSoalRepo struct{ mock.Mock }

func (m *mockUpdateBankSoalRepo) GetBankSoal(ctx context.Context, f query.BankSoalFilter) ([]bank_soal.BankSoal, error) {
	a := m.Called(ctx, f)
	return a.Get(0).([]bank_soal.BankSoal), a.Error(1)
}
func (m *mockUpdateBankSoalRepo) GetBankSoalUploaded(ctx context.Context, f query.BankSoalFilter) ([]bank_soal.BankSoal, error) {
	a := m.Called(ctx, f)
	return a.Get(0).([]bank_soal.BankSoal), a.Error(1)
}
func (m *mockUpdateBankSoalRepo) GetBankSoalByGuru(ctx context.Context, id bank_soal.ID) ([]bank_soal.BankSoal, error) {
	a := m.Called(ctx, id)
	return a.Get(0).([]bank_soal.BankSoal), a.Error(1)
}
func (m *mockUpdateBankSoalRepo) GetBankSoalById(ctx context.Context, id bank_soal.ID) (bank_soal.BankSoal, error) {
	a := m.Called(ctx, id)
	return a.Get(0).(bank_soal.BankSoal), a.Error(1)
}
func (m *mockUpdateBankSoalRepo) CreateBankSoal(ctx context.Context, b bank_soal.BankSoal) error {
	return m.Called(ctx, b).Error(0)
}
func (m *mockUpdateBankSoalRepo) UpdateBankSoal(ctx context.Context, id bank_soal.ID, b updatepatch.UpdateBankSoalPatch) error {
	return m.Called(ctx, id, b).Error(0)
}
func (m *mockUpdateBankSoalRepo) DeleteBankSoal(ctx context.Context, id bank_soal.ID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockUpdateBankSoalRepo) GetIdBankSoalByAttemptId(ctx context.Context, id ujian.ID) (ujian.ID, error) {
	a := m.Called(ctx, id)
	return a.Get(0).(ujian.ID), a.Error(1)
}

func newUpdateBankSoalHandler(repo *mockUpdateBankSoalRepo) *bank_soal_update.UpdateBankSoalHandler {
	return bank_soal_update.NewUpdateBankSoalHandler(bank_soal_update_service.NewUpdateBankSoalService(repo))
}

func updateBankSoalRequest(method, body string) *http.Request {
	req := httptest.NewRequest(method, "/bank-soal/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestUpdateBankSoal_SuccessPartialPayload(t *testing.T) {
	repo := new(mockUpdateBankSoalRepo)
	handler := newUpdateBankSoalHandler(repo)
	req := updateBankSoalRequest(http.MethodPatch, `{"nama_bank_soal":" Bank Baru ","id_mapel":5}`)
	rec := httptest.NewRecorder()
	params := httprouter.Params{{Key: "idBankSoal", Value: "12"}}

	repo.On("UpdateBankSoal", mock.Anything, bank_soal.ID(12), mock.MatchedBy(func(p updatepatch.UpdateBankSoalPatch) bool {
		return p.NamaBankSoal != nil && *p.NamaBankSoal == "Bank Baru" &&
			p.IdMapel != nil && *p.IdMapel == 5 &&
			p.IdKelas == nil && p.Deskripsi == nil && p.Materi == nil
	})).Return(nil).Once()

	handler.UpdateBankSoal(rec, req, params)

	assert.Equal(t, http.StatusOK, rec.Code)
	repo.AssertExpectations(t)
}

func TestUpdateBankSoal_RequestValidation(t *testing.T) {
	tests := []struct {
		name       string
		req        *http.Request
		params     httprouter.Params
		wantStatus int
		wantBody   string
	}{
		{
			name:       "method not allowed",
			req:        updateBankSoalRequest(http.MethodPost, `{}`),
			params:     httprouter.Params{{Key: "idBankSoal", Value: "1"}},
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "method not allowed",
		},
		{
			name:       "invalid id",
			req:        updateBankSoalRequest(http.MethodPatch, `{}`),
			params:     httprouter.Params{{Key: "idBankSoal", Value: "abc"}},
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid id bank soal",
		},
		{
			name: "missing content type",
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodPatch, "/bank-soal/1", strings.NewReader(`{}`))
			}(),
			params:     httprouter.Params{{Key: "idBankSoal", Value: "1"}},
			wantStatus: http.StatusBadRequest,
			wantBody:   "content type must be application/json",
		},
		{
			name:       "invalid json",
			req:        updateBankSoalRequest(http.MethodPatch, `{"nama_bank_soal":`),
			params:     httprouter.Params{{Key: "idBankSoal", Value: "1"}},
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid request body",
		},
		{
			name:       "invalid printable text",
			req:        updateBankSoalRequest(http.MethodPatch, "{\"nama_bank_soal\":\"bad\\u0001\"}"),
			params:     httprouter.Params{{Key: "idBankSoal", Value: "1"}},
			wantStatus: http.StatusBadRequest,
			wantBody:   "nama_bank_soal",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := newUpdateBankSoalHandler(new(mockUpdateBankSoalRepo))
			rec := httptest.NewRecorder()

			handler.UpdateBankSoal(rec, tc.req, tc.params)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantBody)
		})
	}
}

func TestUpdateBankSoal_ServiceErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{name: "no fields", err: coreerror.ErrNoFieldToUpdate, wantStatus: http.StatusBadRequest, wantBody: "invalid update payload"},
		{name: "not found", err: coreerror.ErrNotFound, wantStatus: http.StatusNotFound, wantBody: "data not found"},
		{name: "internal", err: errors.New("db down"), wantStatus: http.StatusInternalServerError, wantBody: "internal server error"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mockUpdateBankSoalRepo)
			handler := newUpdateBankSoalHandler(repo)
			req := updateBankSoalRequest(http.MethodPatch, `{"nama_bank_soal":"Bank Baru"}`)
			rec := httptest.NewRecorder()
			repo.On("UpdateBankSoal", mock.Anything, bank_soal.ID(1), mock.Anything).Return(tc.err).Once()

			handler.UpdateBankSoal(rec, req, httprouter.Params{{Key: "idBankSoal", Value: "1"}})

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantBody)
			repo.AssertExpectations(t)
		})
	}
}
