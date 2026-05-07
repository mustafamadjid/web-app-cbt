package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	ruangujian_get "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ruang_ujian/get"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ruangujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ruang_ujian"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	ruangujianquery "github.com/mustafamadjid/web-app-cbt/internal/core/query/ruang_ujian"
	ruangujian_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/get"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRuangUjianRepo struct{ mock.Mock }

func (m *MockRuangUjianRepo) GetRuangUjian(ctx context.Context, f ruangujianquery.ListRuangUjianFilter) ([]ruangujian.RuangUjian, error) {
	a := m.Called(ctx, f); return a.Get(0).([]ruangujian.RuangUjian), a.Error(1)
}
func (m *MockRuangUjianRepo) GetRuangUjianById(ctx context.Context, id int) (ruangujian.RuangUjian, error) {
	a := m.Called(ctx, id); return a.Get(0).(ruangujian.RuangUjian), a.Error(1)
}
func (m *MockRuangUjianRepo) GetRuangUjianByKode(ctx context.Context, k string) (ruangujian.RuangUjian, error) {
	a := m.Called(ctx, k); return a.Get(0).(ruangujian.RuangUjian), a.Error(1)
}
func (m *MockRuangUjianRepo) ExistByKodeRuang(ctx context.Context, k string) (bool, error) {
	a := m.Called(ctx, k); return a.Bool(0), a.Error(1)
}
func (m *MockRuangUjianRepo) CreateRuangUjian(ctx context.Context, r ruangujian.RuangUjian) error {
	return m.Called(ctx, r).Error(0)
}
func (m *MockRuangUjianRepo) UpdateRuangUjian(ctx context.Context, id int, r updatepatch.UpdateRuangUjianPatch) error {
	return m.Called(ctx, id, r).Error(0)
}
func (m *MockRuangUjianRepo) DeleteRuangUjian(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

func TestListRuangUjian_Success(t *testing.T) {
	mockRepo := new(MockRuangUjianRepo)
	svc := ruangujian_get_service.NewGetRuangUjianService(mockRepo)
	h := ruangujian_get.NewGetRuangUjianHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/ruang-ujian?q=lab&limit=5&offset=2", nil)
	w := httptest.NewRecorder()
	mockRepo.On("GetRuangUjian", mock.Anything, mock.MatchedBy(func(f ruangujianquery.ListRuangUjianFilter) bool {
		return f.Search == "lab" && f.Limit == 5 && f.Offset == 2
	})).Return([]ruangujian.RuangUjian{{IdRuangan: 1, NamaRuangan: "Lab 1", KodeRuang: "LAB1"}}, nil).Once()
	h.GetRuangUjian(w, req, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Lab 1")
}

func TestListRuangUjian_BadRequest_InvalidFilter(t *testing.T) {
	mockRepo := new(MockRuangUjianRepo)
	svc := ruangujian_get_service.NewGetRuangUjianService(mockRepo)
	h := ruangujian_get.NewGetRuangUjianHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/ruang-ujian?offset=abc", nil)
	w := httptest.NewRecorder()
	h.GetRuangUjian(w, req, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "offset must be a number")
}

func TestGetRuangUjianByID_Success(t *testing.T) {
	mockRepo := new(MockRuangUjianRepo)
	svc := ruangujian_get_service.NewGetRuangUjianService(mockRepo)
	h := ruangujian_get.NewGetRuangUjianHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/ruang-ujian/1", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "IdRuangan", Value: "1"}}
	mockRepo.On("GetRuangUjianById", mock.Anything, 1).Return(ruangujian.RuangUjian{IdRuangan: 1, NamaRuangan: "Ruangan 1"}, nil).Once()
	h.GetRuangUjianByID(w, req, ps)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetRuangUjianByID_BadRequest_InvalidParam(t *testing.T) {
	mockRepo := new(MockRuangUjianRepo)
	svc := ruangujian_get_service.NewGetRuangUjianService(mockRepo)
	h := ruangujian_get.NewGetRuangUjianHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/ruang-ujian/abc", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "IdRuangan", Value: "abc"}}
	h.GetRuangUjianByID(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id ruangan")
}

func TestGetRuangUjianByID_NotFound(t *testing.T) {
	mockRepo := new(MockRuangUjianRepo)
	svc := ruangujian_get_service.NewGetRuangUjianService(mockRepo)
	h := ruangujian_get.NewGetRuangUjianHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/ruang-ujian/99", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "IdRuangan", Value: "99"}}
	mockRepo.On("GetRuangUjianById", mock.Anything, 99).Return(ruangujian.RuangUjian{}, coreerror.ErrNotFound).Once()
	h.GetRuangUjianByID(w, req, ps)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
