package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	matapelajaran_update "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/mata_pelajaran/update"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	matapelajaranquery "github.com/mustafamadjid/web-app-cbt/internal/core/query/mata_pelajaran"
	matapelajaran_update_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMapelRepo struct{ mock.Mock }

func (m *MockMapelRepo) GetMapel(ctx context.Context, f matapelajaranquery.ListMapelFilter) ([]matapelajaran.MataPelajaran, error) {
	a := m.Called(ctx, f); return a.Get(0).([]matapelajaran.MataPelajaran), a.Error(1)
}
func (m *MockMapelRepo) GetMapelById(ctx context.Context, id int) (matapelajaran.MataPelajaran, error) {
	a := m.Called(ctx, id); return a.Get(0).(matapelajaran.MataPelajaran), a.Error(1)
}
func (m *MockMapelRepo) CreateMapel(ctx context.Context, s matapelajaran.MataPelajaran) error {
	return m.Called(ctx, s).Error(0)
}
func (m *MockMapelRepo) UpdateMapel(ctx context.Context, id int, s updatepatch.UpdateMapelPatch) error {
	return m.Called(ctx, id, s).Error(0)
}
func (m *MockMapelRepo) DeleteMapel(ctx context.Context, id int) error { return m.Called(ctx, id).Error(0) }
func (m *MockMapelRepo) ExistKodeMapel(ctx context.Context, k string) (bool, error) {
	a := m.Called(ctx, k); return a.Bool(0), a.Error(1)
}

func TestUpdateMapel_Success(t *testing.T) {
	mockRepo := new(MockMapelRepo)
	svc := matapelajaran_update_service.NewUpdateMapelService(mockRepo)
	h := matapelajaran_update.NewUpdateMapelHandler(svc)
	req := httptest.NewRequest(http.MethodPatch, "/mata-pelajaran/1", strings.NewReader(`{"nama_mapel":"Fisika"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idMapel", Value: "1"}}
	mockRepo.On("UpdateMapel", mock.Anything, 1, mock.MatchedBy(func(p updatepatch.UpdateMapelPatch) bool {
		return p.NamaMapel != nil && *p.NamaMapel == "Fisika"
	})).Return(nil).Once()
	h.UpdateMapel(w, req, ps)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestUpdateMapel_BadRequest_InvalidParam(t *testing.T) {
	mockRepo := new(MockMapelRepo)
	svc := matapelajaran_update_service.NewUpdateMapelService(mockRepo)
	h := matapelajaran_update.NewUpdateMapelHandler(svc)
	req := httptest.NewRequest(http.MethodPatch, "/mata-pelajaran/0", strings.NewReader(`{"nama_mapel":"Fisika"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idMapel", Value: "0"}}
	h.UpdateMapel(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id mapel")
}

func TestUpdateMapel_BadRequest_NoField(t *testing.T) {
	mockRepo := new(MockMapelRepo)
	svc := matapelajaran_update_service.NewUpdateMapelService(mockRepo)
	h := matapelajaran_update.NewUpdateMapelHandler(svc)
	req := httptest.NewRequest(http.MethodPatch, "/mata-pelajaran/1", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idMapel", Value: "1"}}
	h.UpdateMapel(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid field")
}

func TestUpdateMapel_NotFound(t *testing.T) {
	mockRepo := new(MockMapelRepo)
	svc := matapelajaran_update_service.NewUpdateMapelService(mockRepo)
	h := matapelajaran_update.NewUpdateMapelHandler(svc)
	req := httptest.NewRequest(http.MethodPatch, "/mata-pelajaran/99", strings.NewReader(`{"nama_mapel":"Fisika"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idMapel", Value: "99"}}
	mockRepo.On("UpdateMapel", mock.Anything, 99, mock.Anything).Return(coreerror.ErrNotFound).Once()
	h.UpdateMapel(w, req, ps)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateMapel_BadRequest_MissingContentType(t *testing.T) {
	mockRepo := new(MockMapelRepo)
	svc := matapelajaran_update_service.NewUpdateMapelService(mockRepo)
	h := matapelajaran_update.NewUpdateMapelHandler(svc)
	req := httptest.NewRequest(http.MethodPatch, "/mata-pelajaran/1", strings.NewReader(`{"nama_mapel":"Fisika"}`))
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idMapel", Value: "1"}}
	h.UpdateMapel(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "content type must be application/json")
}
