package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	bank_soal_delete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/crud_operations/delete"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
	bank_soal_delete_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/delete"
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

func TestDeleteBankSoal_Success(t *testing.T) {
	mockRepo := new(MockBankSoalRepo)
	svc := bank_soal_delete_service.NewDeleteBankSoalService(mockRepo)
	h := bank_soal_delete.NewDeleteBankSoalHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/bank-soal/1", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idBankSoal", Value: "1"}}
	mockRepo.On("DeleteBankSoal", mock.Anything, bank_soal.ID(1)).Return(nil).Once()
	h.DeleteBankSoal(w, req, ps)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteBankSoal_BadRequest_InvalidParam(t *testing.T) {
	mockRepo := new(MockBankSoalRepo)
	svc := bank_soal_delete_service.NewDeleteBankSoalService(mockRepo)
	h := bank_soal_delete.NewDeleteBankSoalHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/bank-soal/abc", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idBankSoal", Value: "abc"}}
	h.DeleteBankSoal(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id bank soal")
}

func TestDeleteBankSoal_NotFound(t *testing.T) {
	mockRepo := new(MockBankSoalRepo)
	svc := bank_soal_delete_service.NewDeleteBankSoalService(mockRepo)
	h := bank_soal_delete.NewDeleteBankSoalHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/bank-soal/99", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idBankSoal", Value: "99"}}
	mockRepo.On("DeleteBankSoal", mock.Anything, bank_soal.ID(99)).Return(coreerror.ErrNotFound).Once()
	h.DeleteBankSoal(w, req, ps)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteBankSoal_Restricted(t *testing.T) {
	mockRepo := new(MockBankSoalRepo)
	svc := bank_soal_delete_service.NewDeleteBankSoalService(mockRepo)
	h := bank_soal_delete.NewDeleteBankSoalHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/bank-soal/2", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idBankSoal", Value: "2"}}
	mockRepo.On("DeleteBankSoal", mock.Anything, bank_soal.ID(2)).Return(coreerror.ErrDeleteRestricted).Once()
	h.DeleteBankSoal(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "DELETE_RESTRICTED")
}
