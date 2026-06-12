package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	bank_soal_get "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/crud_operations/get"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
	bank_soal_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/get"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBankSoalRepo struct{ mock.Mock }

func (m *MockBankSoalRepo) GetBankSoal(ctx context.Context, f query.BankSoalFilter) ([]bank_soal.BankSoal, error) {
	a := m.Called(ctx, f); return a.Get(0).([]bank_soal.BankSoal), a.Error(1)
}
func (m *MockBankSoalRepo) GetBankSoalUploaded(ctx context.Context, f query.BankSoalFilter) ([]bank_soal.BankSoal, error) {
	a := m.Called(ctx, f); return a.Get(0).([]bank_soal.BankSoal), a.Error(1)
}
func (m *MockBankSoalRepo) GetBankSoalByGuru(ctx context.Context, id bank_soal.ID) ([]bank_soal.BankSoal, error) {
	a := m.Called(ctx, id); return a.Get(0).([]bank_soal.BankSoal), a.Error(1)
}
func (m *MockBankSoalRepo) GetBankSoalById(ctx context.Context, id bank_soal.ID) (bank_soal.BankSoal, error) {
	a := m.Called(ctx, id); return a.Get(0).(bank_soal.BankSoal), a.Error(1)
}
func (m *MockBankSoalRepo) CreateBankSoal(ctx context.Context, b bank_soal.BankSoal) error {
	return m.Called(ctx, b).Error(0)
}
func (m *MockBankSoalRepo) UpdateBankSoal(ctx context.Context, id bank_soal.ID, b updatepatch.UpdateBankSoalPatch) error {
	return m.Called(ctx, id, b).Error(0)
}
func (m *MockBankSoalRepo) DeleteBankSoal(ctx context.Context, id bank_soal.ID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockBankSoalRepo) GetIdBankSoalByAttemptId(ctx context.Context, id ujian.ID) (ujian.ID, error) {
	a := m.Called(ctx, id); return a.Get(0).(ujian.ID), a.Error(1)
}

func TestGetBankSoal_Success(t *testing.T) {
	mockRepo := new(MockBankSoalRepo)
	svc := bank_soal_get_service.NewGetBankSoalService(mockRepo)
	h := bank_soal_get.NewGetBankSoalHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/bank-soal?q=fisika&limit=10&offset=0", nil)
	w := httptest.NewRecorder()
	mockRepo.On("GetBankSoal", mock.Anything, mock.MatchedBy(func(f query.BankSoalFilter) bool {
		return f.Search == "fisika" && f.Limit == 10 && f.Offset == 0
	})).Return([]bank_soal.BankSoal{{IdBankSoal: 1, NamaBankSoal: "Fisika"}}, nil).Once()
	h.GetBankSoal(w, req, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Fisika")
}

func TestGetBankSoal_BadRequest_InvalidFilter(t *testing.T) {
	mockRepo := new(MockBankSoalRepo)
	svc := bank_soal_get_service.NewGetBankSoalService(mockRepo)
	h := bank_soal_get.NewGetBankSoalHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/bank-soal?limit=abc", nil)
	w := httptest.NewRecorder()
	h.GetBankSoal(w, req, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetBankSoalByID_Success(t *testing.T) {
	mockRepo := new(MockBankSoalRepo)
	svc := bank_soal_get_service.NewGetBankSoalService(mockRepo)
	h := bank_soal_get.NewGetBankSoalHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/bank-soal/1", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idBankSoal", Value: "1"}}
	mockRepo.On("GetBankSoalById", mock.Anything, bank_soal.ID(1)).Return(bank_soal.BankSoal{IdBankSoal: 1, NamaBankSoal: "Fisika"}, nil).Once()
	h.GetBankSoalByID(w, req, ps)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetBankSoalByID_BadRequest_InvalidParam(t *testing.T) {
	mockRepo := new(MockBankSoalRepo)
	svc := bank_soal_get_service.NewGetBankSoalService(mockRepo)
	h := bank_soal_get.NewGetBankSoalHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/bank-soal/abc", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idBankSoal", Value: "abc"}}
	h.GetBankSoalByID(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id bank soal")
}

func TestGetBankSoalByID_NotFound(t *testing.T) {
	mockRepo := new(MockBankSoalRepo)
	svc := bank_soal_get_service.NewGetBankSoalService(mockRepo)
	h := bank_soal_get.NewGetBankSoalHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/bank-soal/99", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idBankSoal", Value: "99"}}
	mockRepo.On("GetBankSoalById", mock.Anything, bank_soal.ID(99)).Return(bank_soal.BankSoal{}, coreerror.ErrNotFound).Once()
	h.GetBankSoalByID(w, req, ps)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
