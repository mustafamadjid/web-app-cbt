package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	matapelajaran_get "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/mata_pelajaran/get"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	matapelajaranquery "github.com/mustafamadjid/web-app-cbt/internal/core/query/mata_pelajaran"
	matapelajaran_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/get"
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

func TestListMapel_Success(t *testing.T) {
	mockRepo := new(MockMapelRepo)
	svc := matapelajaran_get_service.NewGetMapelService(mockRepo)
	h := matapelajaran_get.NewGetMapelHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/mata-pelajaran?q=mat&limit=10&offset=1&tingkat_kelas=10&nama_mapel=Matematika", nil)
	w := httptest.NewRecorder()
	mockRepo.On("GetMapel", mock.Anything, mock.MatchedBy(func(f matapelajaranquery.ListMapelFilter) bool {
		return f.Search == "mat" && f.Limit == 10 && f.Offset == 1 && f.TingkatKelas != nil && *f.TingkatKelas == 10
	})).Return([]matapelajaran.MataPelajaran{{IdMapel: 1, NamaMapel: "Matematika"}}, nil).Once()
	h.ListMapel(w, req, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Matematika")
}

func TestListMapel_BadRequest_InvalidFilter(t *testing.T) {
	mockRepo := new(MockMapelRepo)
	svc := matapelajaran_get_service.NewGetMapelService(mockRepo)
	h := matapelajaran_get.NewGetMapelHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/mata-pelajaran?limit=abc", nil)
	w := httptest.NewRecorder()
	h.ListMapel(w, req, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "limit must be a number")
}

func TestGetMapelByID_Success(t *testing.T) {
	mockRepo := new(MockMapelRepo)
	svc := matapelajaran_get_service.NewGetMapelService(mockRepo)
	h := matapelajaran_get.NewGetMapelHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/mata-pelajaran/1", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idMapel", Value: "1"}}
	mockRepo.On("GetMapelById", mock.Anything, 1).Return(matapelajaran.MataPelajaran{IdMapel: 1, NamaMapel: "Matematika"}, nil).Once()
	h.GetMapelByID(w, req, ps)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetMapelByID_BadRequest_InvalidParam(t *testing.T) {
	mockRepo := new(MockMapelRepo)
	svc := matapelajaran_get_service.NewGetMapelService(mockRepo)
	h := matapelajaran_get.NewGetMapelHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/mata-pelajaran/abc", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idMapel", Value: "abc"}}
	h.GetMapelByID(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id mapel")
}

func TestGetMapelByID_NotFound(t *testing.T) {
	mockRepo := new(MockMapelRepo)
	svc := matapelajaran_get_service.NewGetMapelService(mockRepo)
	h := matapelajaran_get.NewGetMapelHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/mata-pelajaran/99", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idMapel", Value: "99"}}
	mockRepo.On("GetMapelById", mock.Anything, 99).Return(matapelajaran.MataPelajaran{}, coreerror.ErrNotFound).Once()
	h.GetMapelByID(w, req, ps)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
