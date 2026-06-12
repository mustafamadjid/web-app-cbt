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

	pengumuman_create "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/pengumuman/create"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	pengumuman_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/create"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCreatePengumumanRepo struct{ mock.Mock }

func (m *mockCreatePengumumanRepo) GetPengumumanActive(ctx context.Context) ([]pengumuman.Pengumuman, error) {
	a := m.Called(ctx)
	return a.Get(0).([]pengumuman.Pengumuman), a.Error(1)
}
func (m *mockCreatePengumumanRepo) GetPengumumanNonActive(ctx context.Context) ([]pengumuman.Pengumuman, error) {
	a := m.Called(ctx)
	return a.Get(0).([]pengumuman.Pengumuman), a.Error(1)
}
func (m *mockCreatePengumumanRepo) GetPengumumanIncoming(ctx context.Context) ([]pengumuman.Pengumuman, error) {
	a := m.Called(ctx)
	return a.Get(0).([]pengumuman.Pengumuman), a.Error(1)
}
func (m *mockCreatePengumumanRepo) GetPengumumanById(ctx context.Context, id pengumuman.ID) (pengumuman.Pengumuman, error) {
	a := m.Called(ctx, id)
	return a.Get(0).(pengumuman.Pengumuman), a.Error(1)
}
func (m *mockCreatePengumumanRepo) CreatePengumuman(ctx context.Context, p pengumuman.Pengumuman) error {
	return m.Called(ctx, p).Error(0)
}
func (m *mockCreatePengumumanRepo) UpdatePengumuman(ctx context.Context, id pengumuman.ID, p updatepatch.PengumumanUpdatePatch) error {
	return m.Called(ctx, id, p).Error(0)
}
func (m *mockCreatePengumumanRepo) DeletePengumuman(ctx context.Context, id pengumuman.ID) error {
	return m.Called(ctx, id).Error(0)
}

func newCreatePengumumanHandler(t *testing.T, repo *mockCreatePengumumanRepo, maxBytes int64) *pengumuman_create.CreatePengumumanHandler {
	t.Helper()
	store := httphelper.DocumentStore{Dir: t.TempDir(), Route: "/documents", MaxBytes: maxBytes}
	return pengumuman_create.NewCreatePengumumanHandler(pengumuman_create_service.NewCreatePengumumanRepo(repo), store)
}

func newMultipartPengumumanRequest(t *testing.T, method string, values map[string]string, fileField string, fileName string, fileContent []byte, withActor bool) *http.Request {
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
	req := httptest.NewRequest(method, "/pengumuman", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if withActor {
		req = req.WithContext(middleware.WithActor(req.Context(), user.Actor{IdPengguna: 9, Role: user.ADMIN}))
	}
	return req
}

func validPengumumanValues() map[string]string {
	return map[string]string{
		"judul_pengumuman":           " Pengumuman Baru ",
		"isi_pengumuman":             "Isi pengumuman",
		"tanggal_rilis_pengumuman":   "2026-03-01",
		"tanggal_selesai_pengumuman": "2026-03-31",
	}
}

func TestCreatePengumuman_SuccessWithoutDocument(t *testing.T) {
	repo := new(mockCreatePengumumanRepo)
	handler := newCreatePengumumanHandler(t, repo, 1024)
	req := newMultipartPengumumanRequest(t, http.MethodPost, validPengumumanValues(), "", "", nil, true)
	rec := httptest.NewRecorder()

	repo.On("CreatePengumuman", mock.Anything, mock.MatchedBy(func(p pengumuman.Pengumuman) bool {
		return p.IdPengguna == 9 &&
			p.JudulPengumuman == "Pengumuman Baru" &&
			p.IsiPengumuman == "Isi pengumuman" &&
			p.TanggalRilisPengumuman == "2026-03-01" &&
			p.TanggalSelesaiPengumuman == "2026-03-31" &&
			p.DokumenPengumuman == ""
	})).Return(nil).Once()

	handler.CreatePengumuman(rec, req, nil)

	assert.Equal(t, http.StatusOK, rec.Code)
	repo.AssertExpectations(t)
}

func TestCreatePengumuman_SuccessWithDocument(t *testing.T) {
	repo := new(mockCreatePengumumanRepo)
	handler := newCreatePengumumanHandler(t, repo, 1024)
	req := newMultipartPengumumanRequest(t, http.MethodPost, validPengumumanValues(), "dokumen_pengumuman", "pengumuman.pdf", []byte("%PDF-1.4\nbody"), true)
	rec := httptest.NewRecorder()

	repo.On("CreatePengumuman", mock.Anything, mock.MatchedBy(func(p pengumuman.Pengumuman) bool {
		return p.DokumenPengumuman != "" && filepath.Ext(p.DokumenPengumuman) == ".pdf"
	})).Return(nil).Once()

	handler.CreatePengumuman(rec, req, nil)

	assert.Equal(t, http.StatusOK, rec.Code)
	repo.AssertExpectations(t)
}

func TestCreatePengumuman_RequestValidation(t *testing.T) {
	tests := []struct {
		name       string
		req        func(t *testing.T) *http.Request
		wantStatus int
		wantBody   string
	}{
		{
			name: "method not allowed",
			req: func(t *testing.T) *http.Request {
				return newMultipartPengumumanRequest(t, http.MethodGet, validPengumumanValues(), "", "", nil, true)
			},
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "method not allowed",
		},
		{
			name: "missing multipart content type",
			req: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodPost, "/pengumuman", bytes.NewBufferString(""))
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "content type must be multipart/form-data",
		},
		{
			name: "missing fields",
			req: func(t *testing.T) *http.Request {
				return newMultipartPengumumanRequest(t, http.MethodPost, map[string]string{"judul_pengumuman": "Judul"}, "", "", nil, true)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "missing fields",
		},
		{
			name: "invalid title",
			req: func(t *testing.T) *http.Request {
				values := validPengumumanValues()
				values["judul_pengumuman"] = "bad\u0001"
				return newMultipartPengumumanRequest(t, http.MethodPost, values, "", "", nil, true)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "judul_pengumuman",
		},
		{
			name: "missing actor",
			req: func(t *testing.T) *http.Request {
				return newMultipartPengumumanRequest(t, http.MethodPost, validPengumumanValues(), "", "", nil, false)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "failed get actor from context",
		},
		{
			name: "invalid document",
			req: func(t *testing.T) *http.Request {
				return newMultipartPengumumanRequest(t, http.MethodPost, validPengumumanValues(), "dokumen_pengumuman", "bad.txt", []byte("plain text"), true)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid dokumen_pengumuman",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := newCreatePengumumanHandler(t, new(mockCreatePengumumanRepo), 1024)
			rec := httptest.NewRecorder()

			handler.CreatePengumuman(rec, tc.req(t), nil)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantBody)
		})
	}
}

func TestCreatePengumuman_FileTooLarge(t *testing.T) {
	handler := newCreatePengumumanHandler(t, new(mockCreatePengumumanRepo), 4)
	req := newMultipartPengumumanRequest(t, http.MethodPost, validPengumumanValues(), "dokumen_pengumuman", "pengumuman.pdf", []byte("%PDF-1.4\nbody"), true)
	rec := httptest.NewRecorder()

	handler.CreatePengumuman(rec, req, nil)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "file too large")
}

func TestCreatePengumuman_ServiceErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{name: "invalid date", err: coreerror.ErrInvalidDateFormat, wantStatus: http.StatusBadRequest, wantBody: "invalid date format"},
		{name: "invalid input", err: coreerror.ErrInvalidInput, wantStatus: http.StatusBadRequest, wantBody: "invalid input"},
		{name: "internal", err: errors.New("db down"), wantStatus: http.StatusInternalServerError, wantBody: "failed create pengumuman"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mockCreatePengumumanRepo)
			handler := newCreatePengumumanHandler(t, repo, 1024)
			req := newMultipartPengumumanRequest(t, http.MethodPost, validPengumumanValues(), "", "", nil, true)
			rec := httptest.NewRecorder()
			repo.On("CreatePengumuman", mock.Anything, mock.Anything).Return(tc.err).Once()

			handler.CreatePengumuman(rec, req, nil)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantBody)
			repo.AssertExpectations(t)
		})
	}
}
