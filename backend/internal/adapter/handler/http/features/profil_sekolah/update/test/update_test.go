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

	profil_update "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/profil_sekolah/update"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/profil_sekolah"
	profil_update_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/profil_sekolah/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockProfilSekolahRepo struct{ mock.Mock }

func (m *mockProfilSekolahRepo) GetProfilSekolah(ctx context.Context) (profil_sekolah.ProfilSekolah, error) {
	a := m.Called(ctx)
	return a.Get(0).(profil_sekolah.ProfilSekolah), a.Error(1)
}

func (m *mockProfilSekolahRepo) UpdateProfilSekolah(ctx context.Context, id profil_sekolah.IDProfil, profil profil_sekolah.ProfilSekolah) error {
	return m.Called(ctx, id, profil).Error(0)
}

func newUpdateProfilSekolahHandler(t *testing.T, repo *mockProfilSekolahRepo, maxBytes int64) *profil_update.UpdateProfilSekolahHandler {
	t.Helper()
	store := httphelper.ImageStore{Dir: t.TempDir(), Route: "/images", MaxBytes: maxBytes}
	return profil_update.NewUpdateProfilSekolahHandler(profil_update_service.NewUpdateProfilSekolahService(repo), store)
}

func newMultipartProfilRequest(t *testing.T, method string, values map[string]string, fileField string, fileName string, fileContent []byte) *http.Request {
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
	req := httptest.NewRequest(method, "/profil-sekolah", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func pngBytes() []byte {
	return []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n', 0, 0, 0, 0}
}

func existingProfilSekolah() profil_sekolah.ProfilSekolah {
	logo := "/images/old.png"
	return profil_sekolah.ProfilSekolah{
		IDProfil:      1,
		EmailSekolah:  "old@example.com",
		NoTelpSekolah: "08123",
		KepalaSekolah: "Kepala Lama",
		WakaSekolah:   "Waka Lama",
		NamaSekolah:   "Sekolah Lama",
		AlamatSekolah: "Alamat Lama",
		LogoSekolah:   &logo,
	}
}

func TestUpdateProfilSekolah_SuccessTextOnly(t *testing.T) {
	repo := new(mockProfilSekolahRepo)
	handler := newUpdateProfilSekolahHandler(t, repo, 1024)
	req := newMultipartProfilRequest(t, http.MethodPatch, map[string]string{
		"email_sekolah": " admin@example.com ",
		"nama_sekolah":  " SMA Baru ",
	}, "", "", nil)
	rec := httptest.NewRecorder()

	repo.On("GetProfilSekolah", mock.Anything).Return(existingProfilSekolah(), nil).Once()
	repo.On("UpdateProfilSekolah", mock.Anything, profil_sekolah.IDProfil(1), mock.MatchedBy(func(p profil_sekolah.ProfilSekolah) bool {
		return p.EmailSekolah == "admin@example.com" &&
			p.NamaSekolah == "SMA Baru" &&
			p.NoTelpSekolah == "08123" &&
			p.LogoSekolah != nil && *p.LogoSekolah == "/images/old.png"
	})).Return(nil).Once()

	handler.UpdateProfilSekolah(rec, req, nil)

	assert.Equal(t, http.StatusOK, rec.Code)
	repo.AssertExpectations(t)
}

func TestUpdateProfilSekolah_SuccessWithLogo(t *testing.T) {
	repo := new(mockProfilSekolahRepo)
	handler := newUpdateProfilSekolahHandler(t, repo, 1024)
	req := newMultipartProfilRequest(t, http.MethodPatch, map[string]string{"nama_sekolah": "SMA Baru"}, "logo_sekolah", "logo.png", pngBytes())
	rec := httptest.NewRecorder()

	repo.On("GetProfilSekolah", mock.Anything).Return(existingProfilSekolah(), nil).Once()
	repo.On("UpdateProfilSekolah", mock.Anything, profil_sekolah.IDProfil(1), mock.MatchedBy(func(p profil_sekolah.ProfilSekolah) bool {
		return p.LogoSekolah != nil && filepath.Ext(*p.LogoSekolah) == ".png"
	})).Return(nil).Once()

	handler.UpdateProfilSekolah(rec, req, nil)

	assert.Equal(t, http.StatusOK, rec.Code)
	repo.AssertExpectations(t)
}

func TestUpdateProfilSekolah_RequestValidation(t *testing.T) {
	tests := []struct {
		name       string
		req        func(t *testing.T) *http.Request
		wantStatus int
		wantBody   string
	}{
		{
			name: "method not allowed",
			req: func(t *testing.T) *http.Request {
				return newMultipartProfilRequest(t, http.MethodPost, map[string]string{}, "", "", nil)
			},
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "method not allowed",
		},
		{
			name: "missing multipart content type",
			req: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodPatch, "/profil-sekolah", bytes.NewBufferString(""))
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "content type must be multipart/form-data",
		},
		{
			name: "empty optional email",
			req: func(t *testing.T) *http.Request {
				return newMultipartProfilRequest(t, http.MethodPatch, map[string]string{"email_sekolah": ""}, "", "", nil)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "email_sekolah is required",
		},
		{
			name: "invalid email",
			req: func(t *testing.T) *http.Request {
				return newMultipartProfilRequest(t, http.MethodPatch, map[string]string{"email_sekolah": "not-email"}, "", "", nil)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "email_sekolah",
		},
		{
			name: "invalid printable text",
			req: func(t *testing.T) *http.Request {
				return newMultipartProfilRequest(t, http.MethodPatch, map[string]string{"nama_sekolah": "bad\u0001"}, "", "", nil)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "nama_sekolah",
		},
		{
			name: "invalid logo",
			req: func(t *testing.T) *http.Request {
				return newMultipartProfilRequest(t, http.MethodPatch, map[string]string{"nama_sekolah": "SMA"}, "logo_sekolah", "logo.txt", []byte("plain text"))
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid logo_sekolah",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := newUpdateProfilSekolahHandler(t, new(mockProfilSekolahRepo), 1024)
			rec := httptest.NewRecorder()

			handler.UpdateProfilSekolah(rec, tc.req(t), nil)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantBody)
		})
	}
}

func TestUpdateProfilSekolah_FileTooLarge(t *testing.T) {
	handler := newUpdateProfilSekolahHandler(t, new(mockProfilSekolahRepo), 4)
	req := newMultipartProfilRequest(t, http.MethodPatch, map[string]string{"nama_sekolah": "SMA"}, "logo_sekolah", "logo.png", pngBytes())
	rec := httptest.NewRecorder()

	handler.UpdateProfilSekolah(rec, req, nil)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "file too large")
}

func TestUpdateProfilSekolah_ServiceErrors(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(repo *mockProfilSekolahRepo)
		wantStatus int
		wantBody   string
	}{
		{
			name: "no fields",
			setup: func(repo *mockProfilSekolahRepo) {
				repo.On("GetProfilSekolah", mock.Anything).Return(profil_sekolah.ProfilSekolah{}, coreerror.ErrNoFieldToUpdate).Maybe()
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "no fields to update",
		},
		{
			name: "invalid input",
			setup: func(repo *mockProfilSekolahRepo) {
				repo.On("GetProfilSekolah", mock.Anything).Return(existingProfilSekolah(), nil).Once()
				repo.On("UpdateProfilSekolah", mock.Anything, profil_sekolah.IDProfil(1), mock.Anything).Return(coreerror.ErrInvalidInput).Once()
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid input",
		},
		{
			name: "internal",
			setup: func(repo *mockProfilSekolahRepo) {
				repo.On("GetProfilSekolah", mock.Anything).Return(existingProfilSekolah(), nil).Once()
				repo.On("UpdateProfilSekolah", mock.Anything, profil_sekolah.IDProfil(1), mock.Anything).Return(errors.New("db down")).Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "failed update profil sekolah",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mockProfilSekolahRepo)
			handler := newUpdateProfilSekolahHandler(t, repo, 1024)
			values := map[string]string{"nama_sekolah": "SMA Baru"}
			if tc.name == "no fields" {
				values = map[string]string{}
			}
			req := newMultipartProfilRequest(t, http.MethodPatch, values, "", "", nil)
			rec := httptest.NewRecorder()
			tc.setup(repo)

			handler.UpdateProfilSekolah(rec, req, nil)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantBody)
			repo.AssertExpectations(t)
		})
	}
}
