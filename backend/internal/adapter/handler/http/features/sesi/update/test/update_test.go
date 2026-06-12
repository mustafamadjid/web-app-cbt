package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	sesi_update "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/sesi/update"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	sesiquery "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
	sesi_update_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSesiRepo struct{ mock.Mock }

func (m *MockSesiRepo) GetSesi(ctx context.Context, f sesiquery.ListSesiFilter) ([]sesi.Sesi, error) {
	a := m.Called(ctx, f); return a.Get(0).([]sesi.Sesi), a.Error(1)
}
func (m *MockSesiRepo) GetSesiById(ctx context.Context, id int) (sesi.Sesi, error) {
	a := m.Called(ctx, id); return a.Get(0).(sesi.Sesi), a.Error(1)
}
func (m *MockSesiRepo) GetSesiByKode(ctx context.Context, k string) (sesi.Sesi, error) {
	a := m.Called(ctx, k); return a.Get(0).(sesi.Sesi), a.Error(1)
}
func (m *MockSesiRepo) ExistByKodeSesi(ctx context.Context, k string) (bool, error) {
	a := m.Called(ctx, k); return a.Bool(0), a.Error(1)
}
func (m *MockSesiRepo) CreateSesi(ctx context.Context, s sesi.Sesi) error { return m.Called(ctx, s).Error(0) }
func (m *MockSesiRepo) UpdateSesi(ctx context.Context, id int, s updatepatch.UpdateSesiPatch) error {
	return m.Called(ctx, id, s).Error(0)
}
func (m *MockSesiRepo) DeleteSesi(ctx context.Context, id int) error { return m.Called(ctx, id).Error(0) }

func TestUpdateSesi_Success(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	svc := sesi_update_service.NewUpdateSesiService(mockRepo)
	h := sesi_update.NewUpdateSesiHandler(svc)
	req := httptest.NewRequest(http.MethodPatch, "/sesi/1", strings.NewReader(`{"nama_sesi":"Siang"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idSesi", Value: "1"}}
	mockRepo.On("UpdateSesi", mock.Anything, 1, mock.MatchedBy(func(p updatepatch.UpdateSesiPatch) bool {
		return p.NamaSesi != nil && *p.NamaSesi == "Siang"
	})).Return(nil).Once()
	h.UpdateSesi(w, req, ps)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateSesi_BadRequest_InvalidParam(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	svc := sesi_update_service.NewUpdateSesiService(mockRepo)
	h := sesi_update.NewUpdateSesiHandler(svc)
	req := httptest.NewRequest(http.MethodPatch, "/sesi/abc", strings.NewReader(`{"nama_sesi":"Siang"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idSesi", Value: "abc"}}
	h.UpdateSesi(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id sesi")
}

func TestUpdateSesi_BadRequest_NoField(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	svc := sesi_update_service.NewUpdateSesiService(mockRepo)
	h := sesi_update.NewUpdateSesiHandler(svc)
	req := httptest.NewRequest(http.MethodPatch, "/sesi/1", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idSesi", Value: "1"}}
	h.UpdateSesi(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid field")
}

func TestUpdateSesi_NotFound(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	svc := sesi_update_service.NewUpdateSesiService(mockRepo)
	h := sesi_update.NewUpdateSesiHandler(svc)
	req := httptest.NewRequest(http.MethodPatch, "/sesi/99", strings.NewReader(`{"nama_sesi":"Siang"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idSesi", Value: "99"}}
	mockRepo.On("UpdateSesi", mock.Anything, 99, mock.Anything).Return(coreerror.ErrNotFound).Once()
	h.UpdateSesi(w, req, ps)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateSesi_BadRequest_MissingContentType(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	svc := sesi_update_service.NewUpdateSesiService(mockRepo)
	h := sesi_update.NewUpdateSesiHandler(svc)
	req := httptest.NewRequest(http.MethodPatch, "/sesi/1", strings.NewReader(`{"nama_sesi":"Siang"}`))
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idSesi", Value: "1"}}
	h.UpdateSesi(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "content type must be application/json")
}
