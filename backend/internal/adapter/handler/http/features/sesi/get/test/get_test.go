package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	sesi_get "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/sesi/get"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	authsession "github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/session"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	sesiquery "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
	sesi_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/get"
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

type MockSessionRepo struct{ mock.Mock }

func (m *MockSessionRepo) GetSession(ctx context.Context, sid string) (authsession.Session, error) {
	a := m.Called(ctx, sid); return a.Get(0).(authsession.Session), a.Error(1)
}
func (m *MockSessionRepo) GetSessionByUserId(ctx context.Context, uid user.ID) (authsession.Session, error) {
	a := m.Called(ctx, uid); return a.Get(0).(authsession.Session), a.Error(1)
}
func (m *MockSessionRepo) GetAllActiveSession(ctx context.Context) ([]authsession.SessionWithUser, error) {
	a := m.Called(ctx); return a.Get(0).([]authsession.SessionWithUser), a.Error(1)
}
func (m *MockSessionRepo) CreateSession(ctx context.Context, uid user.ID, role user.Role, exp time.Time) (string, error) {
	a := m.Called(ctx, uid, role, exp); return a.String(0), a.Error(1)
}
func (m *MockSessionRepo) RevokeSession(ctx context.Context, sid string) error { return m.Called(ctx, sid).Error(0) }
func (m *MockSessionRepo) RevokeSessionAllbyUser(ctx context.Context, uid user.ID) error { return m.Called(ctx, uid).Error(0) }
func (m *MockSessionRepo) RevokeExpiredSessions(ctx context.Context, uid user.ID) (bool, error) {
	a := m.Called(ctx, uid); return a.Bool(0), a.Error(1)
}
func (m *MockSessionRepo) HasActiveSession(ctx context.Context, uid user.ID) (bool, error) {
	a := m.Called(ctx, uid); return a.Bool(0), a.Error(1)
}

func TestListSesi_Success(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	mockSR := new(MockSessionRepo)
	svc := sesi_get_service.NewGetSesiService(mockRepo, mockSR)
	h := sesi_get.NewGetSesiHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/sesi?q=pagi&limit=10&offset=2", nil)
	w := httptest.NewRecorder()
	mockRepo.On("GetSesi", mock.Anything, mock.MatchedBy(func(f sesiquery.ListSesiFilter) bool {
		return f.Search == "pagi" && f.Limit == 10 && f.Offset == 2
	})).Return([]sesi.Sesi{{IdSesi: 1, NamaSesi: "Pagi", KodeSesi: "PAGI"}}, nil).Once()
	h.ListSesi(w, req, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Pagi")
}

func TestListSesi_BadRequest_InvalidFilter(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	mockSR := new(MockSessionRepo)
	svc := sesi_get_service.NewGetSesiService(mockRepo, mockSR)
	h := sesi_get.NewGetSesiHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/sesi?limit=abc", nil)
	w := httptest.NewRecorder()
	h.ListSesi(w, req, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "limit must be a number")
}

func TestGetSesiByID_Success(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	mockSR := new(MockSessionRepo)
	svc := sesi_get_service.NewGetSesiService(mockRepo, mockSR)
	h := sesi_get.NewGetSesiHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/sesi/1", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idSesi", Value: "1"}}
	mockRepo.On("GetSesiById", mock.Anything, 1).Return(sesi.Sesi{IdSesi: 1, NamaSesi: "Sesi 1"}, nil).Once()
	h.GetSesiByID(w, req, ps)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetSesiByID_BadRequest_InvalidParam(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	mockSR := new(MockSessionRepo)
	svc := sesi_get_service.NewGetSesiService(mockRepo, mockSR)
	h := sesi_get.NewGetSesiHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/sesi/abc", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idSesi", Value: "abc"}}
	h.GetSesiByID(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id sesi")
}

func TestGetSesiByID_NotFound(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	mockSR := new(MockSessionRepo)
	svc := sesi_get_service.NewGetSesiService(mockRepo, mockSR)
	h := sesi_get.NewGetSesiHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/sesi/99", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idSesi", Value: "99"}}
	mockRepo.On("GetSesiById", mock.Anything, 99).Return(sesi.Sesi{}, coreerror.ErrNotFound).Once()
	h.GetSesiByID(w, req, ps)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetSesiByKode_Success(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	mockSR := new(MockSessionRepo)
	svc := sesi_get_service.NewGetSesiService(mockRepo, mockSR)
	h := sesi_get.NewGetSesiHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/sesi/kode/PAGI", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "kodeSesi", Value: "PAGI"}}
	mockRepo.On("GetSesiByKode", mock.Anything, "PAGI").Return(sesi.Sesi{IdSesi: 1, KodeSesi: "PAGI", NamaSesi: "Pagi"}, nil).Once()
	h.GetSesiByKode(w, req, ps)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "PAGI")
}

func TestListActiveLoginSession_Success(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	mockSR := new(MockSessionRepo)
	svc := sesi_get_service.NewGetSesiService(mockRepo, mockSR)
	h := sesi_get.NewGetSesiHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/sesi/active-login", nil)
	w := httptest.NewRecorder()
	mockSR.On("GetAllActiveSession", mock.Anything).Return([]authsession.SessionWithUser{
		{Session: authsession.Session{SessionID: "session-1", UserID: 1, Role: user.ADMIN, ExpiresAt: time.Now().Add(time.Hour)}, Pengguna: user.Pengguna{ID: 1, Username: "admin", Role: user.ADMIN}},
	}, nil).Once()
	h.ListActiveLoginSession(w, req, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "session-1")
}
