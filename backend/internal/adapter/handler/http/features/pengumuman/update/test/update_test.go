package httpx_test

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/julienschmidt/httprouter"
	pengumuman_update "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/pengumuman/update"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	pengumuman_update_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUpdatePengumumanRepo struct{ mock.Mock }

func (m *mockUpdatePengumumanRepo) GetPengumumanActive(ctx context.Context) ([]pengumuman.Pengumuman, error) {
	a := m.Called(ctx)
	return a.Get(0).([]pengumuman.Pengumuman), a.Error(1)
}
func (m *mockUpdatePengumumanRepo) GetPengumumanNonActive(ctx context.Context) ([]pengumuman.Pengumuman, error) {
	a := m.Called(ctx)
	return a.Get(0).([]pengumuman.Pengumuman), a.Error(1)
}
func (m *mockUpdatePengumumanRepo) GetPengumumanIncoming(ctx context.Context) ([]pengumuman.Pengumuman, error) {
	a := m.Called(ctx)
	return a.Get(0).([]pengumuman.Pengumuman), a.Error(1)
}
func (m *mockUpdatePengumumanRepo) GetPengumumanById(ctx context.Context, id pengumuman.ID) (pengumuman.Pengumuman, error) {
	a := m.Called(ctx, id)
	return a.Get(0).(pengumuman.Pengumuman), a.Error(1)
}
func (m *mockUpdatePengumumanRepo) CreatePengumuman(ctx context.Context, p pengumuman.Pengumuman) error {
	return m.Called(ctx, p).Error(0)
}
func (m *mockUpdatePengumumanRepo) UpdatePengumuman(ctx context.Context, id pengumuman.ID, p updatepatch.PengumumanUpdatePatch) error {
	return m.Called(ctx, id, p).Error(0)
}
func (m *mockUpdatePengumumanRepo) DeletePengumuman(ctx context.Context, id pengumuman.ID) error {
	return m.Called(ctx, id).Error(0)
}

type mockDeletePengumumanFileRepo struct{ mock.Mock }

func (m *mockDeletePengumumanFileRepo) DeleteFile(ctx context.Context, filePath string) error {
	return m.Called(ctx, filePath).Error(0)
}

func newUpdatePengumumanHandler(t *testing.T, repo *mockUpdatePengumumanRepo, deleteRepo *mockDeletePengumumanFileRepo, maxBytes int64) *pengumuman_update.UpdatePengumumanHandler {
	t.Helper()
	store := httphelper.DocumentStore{Dir: t.TempDir(), Route: "/documents", MaxBytes: maxBytes}
	return pengumuman_update.NewUpdatePengumumanHandler(pengumuman_update_service.NewUpdatePengumumanService(repo, deleteRepo), store)
}

func newMultipartUpdatePengumumanRequest(t *testing.T, method string, values map[string]string, fileField string, fileName string, fileContent []byte, withActor bool) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for key, value := range values {
		assert.NoError(t, writer.WriteField(key, value))
	}
	if fileField != "" {
		part, err := writer.CreateFormFile(fileField, fileName)
		assert.NoError(t, err)
		_, err = part.Write(fileContent)
		assert.NoError(t, err)
	}
	assert.NoError(t, writer.Close())
	req := httptest.NewRequest(method, "/pengumuman/11", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if withActor {
		req = req.WithContext(middleware.WithActor(req.Context(), user.Actor{IdPengguna: 4, Role: user.ADMIN}))
	}
	return req
}

func TestUpdatePengumuman_SuccessTextOnly(t *testing.T) {
	repo := new(mockUpdatePengumumanRepo)
	deleteRepo := new(mockDeletePengumumanFileRepo)
	handler := newUpdatePengumumanHandler(t, repo, deleteRepo, 1024)
	req := newMultipartUpdatePengumumanRequest(t, http.MethodPatch, map[string]string{"judul_pengumuman": " Judul Baru "}, "", "", nil, true)
	rec := httptest.NewRecorder()

	repo.On("UpdatePengumuman", mock.Anything, pengumuman.ID(11), mock.MatchedBy(func(p updatepatch.PengumumanUpdatePatch) bool {
		return p.IdPengguna != nil && *p.IdPengguna == 4 &&
			p.JudulPengumuman != nil && *p.JudulPengumuman == "Judul Baru" &&
			p.DokumenPengumuman == nil
	})).Return(nil).Once()

	handler.UpdatePengumuman(rec, req, httprouter.Params{{Key: "idPengumuman", Value: "11"}})

	assert.Equal(t, http.StatusOK, rec.Code)
	repo.AssertExpectations(t)
	deleteRepo.AssertExpectations(t)
}

func TestUpdatePengumuman_SuccessWithDocument(t *testing.T) {
	repo := new(mockUpdatePengumumanRepo)
	deleteRepo := new(mockDeletePengumumanFileRepo)
	handler := newUpdatePengumumanHandler(t, repo, deleteRepo, 1024)
	req := newMultipartUpdatePengumumanRequest(t, http.MethodPatch, map[string]string{"judul_pengumuman": "Judul Baru"}, "dokumen_pengumuman", "pengumuman.pdf", []byte("%PDF-1.4\nbody"), true)
	rec := httptest.NewRecorder()

	repo.On("GetPengumumanById", mock.Anything, pengumuman.ID(11)).Return(pengumuman.Pengumuman{DokumenPengumuman: "/documents/old.pdf"}, nil).Once()
	deleteRepo.On("DeleteFile", mock.Anything, "/documents/old.pdf").Return(nil).Once()
	repo.On("UpdatePengumuman", mock.Anything, pengumuman.ID(11), mock.MatchedBy(func(p updatepatch.PengumumanUpdatePatch) bool {
		return p.DokumenPengumuman != nil && filepath.Ext(*p.DokumenPengumuman) == ".pdf"
	})).Return(nil).Once()

	handler.UpdatePengumuman(rec, req, httprouter.Params{{Key: "idPengumuman", Value: "11"}})

	assert.Equal(t, http.StatusOK, rec.Code)
	repo.AssertExpectations(t)
	deleteRepo.AssertExpectations(t)
}

func TestUpdatePengumuman_RequestValidation(t *testing.T) {
	tests := []struct {
		name       string
		req        func(t *testing.T) *http.Request
		params     httprouter.Params
		wantStatus int
		wantBody   string
	}{
		{
			name: "method not allowed",
			req: func(t *testing.T) *http.Request {
				return newMultipartUpdatePengumumanRequest(t, http.MethodPost, map[string]string{}, "", "", nil, true)
			},
			params:     httprouter.Params{{Key: "idPengumuman", Value: "11"}},
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "method not allowed",
		},
		{
			name: "invalid id",
			req: func(t *testing.T) *http.Request {
				return newMultipartUpdatePengumumanRequest(t, http.MethodPatch, map[string]string{}, "", "", nil, true)
			},
			params:     httprouter.Params{{Key: "idPengumuman", Value: "abc"}},
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid id pengumuman",
		},
		{
			name: "missing multipart content type",
			req: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodPatch, "/pengumuman/11", bytes.NewBufferString(""))
			},
			params:     httprouter.Params{{Key: "idPengumuman", Value: "11"}},
			wantStatus: http.StatusBadRequest,
			wantBody:   "content type must be multipart/form-data",
		},
		{
			name: "empty optional title",
			req: func(t *testing.T) *http.Request {
				return newMultipartUpdatePengumumanRequest(t, http.MethodPatch, map[string]string{"judul_pengumuman": ""}, "", "", nil, true)
			},
			params:     httprouter.Params{{Key: "idPengumuman", Value: "11"}},
			wantStatus: http.StatusBadRequest,
			wantBody:   "judul_pengumuman is required",
		},
		{
			name: "invalid title",
			req: func(t *testing.T) *http.Request {
				return newMultipartUpdatePengumumanRequest(t, http.MethodPatch, map[string]string{"judul_pengumuman": "bad\u0001"}, "", "", nil, true)
			},
			params:     httprouter.Params{{Key: "idPengumuman", Value: "11"}},
			wantStatus: http.StatusBadRequest,
			wantBody:   "judul_pengumuman",
		},
		{
			name: "missing actor",
			req: func(t *testing.T) *http.Request {
				return newMultipartUpdatePengumumanRequest(t, http.MethodPatch, map[string]string{"judul_pengumuman": "Judul"}, "", "", nil, false)
			},
			params:     httprouter.Params{{Key: "idPengumuman", Value: "11"}},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "failed get actor from context",
		},
		{
			name: "invalid document",
			req: func(t *testing.T) *http.Request {
				return newMultipartUpdatePengumumanRequest(t, http.MethodPatch, map[string]string{"judul_pengumuman": "Judul"}, "dokumen_pengumuman", "bad.txt", []byte("plain text"), true)
			},
			params:     httprouter.Params{{Key: "idPengumuman", Value: "11"}},
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid dokumen_pengumuman",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := newUpdatePengumumanHandler(t, new(mockUpdatePengumumanRepo), new(mockDeletePengumumanFileRepo), 1024)
			rec := httptest.NewRecorder()

			handler.UpdatePengumuman(rec, tc.req(t), tc.params)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantBody)
		})
	}
}

func TestUpdatePengumuman_FileTooLarge(t *testing.T) {
	handler := newUpdatePengumumanHandler(t, new(mockUpdatePengumumanRepo), new(mockDeletePengumumanFileRepo), 4)
	req := newMultipartUpdatePengumumanRequest(t, http.MethodPatch, map[string]string{"judul_pengumuman": "Judul"}, "dokumen_pengumuman", "pengumuman.pdf", []byte("%PDF-1.4\nbody"), true)
	rec := httptest.NewRecorder()

	handler.UpdatePengumuman(rec, req, httprouter.Params{{Key: "idPengumuman", Value: "11"}})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "file too large")
}

func TestUpdatePengumuman_ServiceErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{name: "no fields", err: coreerror.ErrNoFieldToUpdate, wantStatus: http.StatusBadRequest, wantBody: "missing fields"},
		{name: "invalid date", err: coreerror.ErrInvalidDateFormat, wantStatus: http.StatusBadRequest, wantBody: "invalid date format"},
		{name: "not found", err: coreerror.ErrNotFound, wantStatus: http.StatusNotFound, wantBody: "data not found"},
		{name: "internal", err: errors.New("db down"), wantStatus: http.StatusInternalServerError, wantBody: "failed update pengumuman"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mockUpdatePengumumanRepo)
			deleteRepo := new(mockDeletePengumumanFileRepo)
			handler := newUpdatePengumumanHandler(t, repo, deleteRepo, 1024)
			req := newMultipartUpdatePengumumanRequest(t, http.MethodPatch, map[string]string{"judul_pengumuman": "Judul"}, "", "", nil, true)
			rec := httptest.NewRecorder()
			repo.On("UpdatePengumuman", mock.Anything, pengumuman.ID(11), mock.Anything).Return(tc.err).Once()

			handler.UpdatePengumuman(rec, req, httprouter.Params{{Key: "idPengumuman", Value: "11"}})

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantBody)
			repo.AssertExpectations(t)
		})
	}
}
