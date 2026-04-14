package sesi_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	sesi_create "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/sesi/create"
	sesi_get "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/sesi/get"
	authsession "github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/session"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	sesi_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/create"
	sesi_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/get"
	sesiquery "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSesiRepository struct {
	mock.Mock
}

func (m *MockSesiRepository) GetSesi(ctx context.Context, filter sesiquery.ListSesiFilter) ([]sesi.Sesi, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]sesi.Sesi), args.Error(1)
}

func (m *MockSesiRepository) GetSesiById(ctx context.Context, idSesi int) (sesi.Sesi, error) {
	args := m.Called(ctx, idSesi)
	return args.Get(0).(sesi.Sesi), args.Error(1)
}

func (m *MockSesiRepository) GetSesiByKode(ctx context.Context, kodeSesi string) (sesi.Sesi, error) {
	args := m.Called(ctx, kodeSesi)
	return args.Get(0).(sesi.Sesi), args.Error(1)
}

func (m *MockSesiRepository) ExistByKodeSesi(ctx context.Context, kodeSesi string) (bool, error) {
	args := m.Called(ctx, kodeSesi)
	return args.Bool(0), args.Error(1)
}

func (m *MockSesiRepository) CreateSesi(ctx context.Context, s sesi.Sesi) error {
	return m.Called(ctx, s).Error(0)
}

func (m *MockSesiRepository) UpdateSesi(ctx context.Context, idSesi int, s updatepatch.UpdateSesiPatch) error {
	return m.Called(ctx, idSesi, s).Error(0)
}

func (m *MockSesiRepository) DeleteSesi(ctx context.Context, idSesi int) error {
	return m.Called(ctx, idSesi).Error(0)
}

type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) GetSession(ctx context.Context, sessionID string) (authsession.Session, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).(authsession.Session), args.Error(1)
}

func (m *MockSessionRepository) GetSessionByUserId(ctx context.Context, userId user.ID) (authsession.Session, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(authsession.Session), args.Error(1)
}

func (m *MockSessionRepository) GetAllActiveSession(ctx context.Context) ([]authsession.SessionWithUser, error) {
	args := m.Called(ctx)
	return args.Get(0).([]authsession.SessionWithUser), args.Error(1)
}

func (m *MockSessionRepository) CreateSession(ctx context.Context, userID user.ID, role user.Role, expiresAt time.Time) (string, error) {
	args := m.Called(ctx, userID, role, expiresAt)
	return args.String(0), args.Error(1)
}

func (m *MockSessionRepository) RevokeSession(ctx context.Context, sessionID string) error {
	return m.Called(ctx, sessionID).Error(0)
}

func (m *MockSessionRepository) RevokeSessionAllbyUser(ctx context.Context, userID user.ID) error {
	return m.Called(ctx, userID).Error(0)
}

func (m *MockSessionRepository) RevokeExpiredSessions(ctx context.Context, userID user.ID) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSessionRepository) HasActiveSession(ctx context.Context, userID user.ID) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func TestSesiHandlers(t *testing.T) {
	mockRepo := new(MockSesiRepository)
	mockSessionRepo := new(MockSessionRepository)
	createSvc := sesi_create_service.NewCreateSesiService(mockRepo)
	getSvc := sesi_get_service.NewGetSesiService(mockRepo, mockSessionRepo)

	createHandler := sesi_create.NewCreateSesiHandler(createSvc)
	getHandler := sesi_get.NewGetSesiHandler(getSvc)

	t.Run("Create Sesi Success", func(t *testing.T) {
		reqBody := `{"nama_sesi":"Sesi 1", "kode_sesi":"SESI1"}`
		req := httptest.NewRequest(http.MethodPost, "/sesi", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockRepo.On("ExistByKodeSesi", mock.Anything, "SESI1").Return(false, nil).Once()
		mockRepo.On("CreateSesi", mock.Anything, mock.Anything).Return(nil).Once()

		createHandler.CreateSesiHandler(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Get Sesi By ID Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sesi/1", nil)
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "idSesi", Value: "1"}}

		mockRepo.On("GetSesiById", mock.Anything, 1).Return(sesi.Sesi{IdSesi: 1, NamaSesi: "Sesi 1"}, nil).Once()

		getHandler.GetSesiByID(w, req, params)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
