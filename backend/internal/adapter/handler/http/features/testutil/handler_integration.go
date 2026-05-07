package handlertestutil

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	repotestutil "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/security/bcrypt"
	"github.com/mustafamadjid/web-app-cbt/internal/app"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/mustafamadjid/web-app-cbt/internal/infra/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Envelope struct {
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
	Error   *APIError       `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func BuildApp(t *testing.T) *app.App {
	t.Helper()

	dbURL := strings.TrimSpace(os.Getenv("POSTGRES_TEST_DBURL"))
	if dbURL == "" {
		t.Skip("POSTGRES_TEST_DBURL is not set")
	}

	uploadDir := t.TempDir()
	cfg := app.Config{
		HTTP: app.HTTPConfig{Addr: ":0"},
		JWT: app.JWTConfig{
			Issuer:        "handler-integration-test",
			AccessSecret:  "handler-access-secret",
			RefreshSecret: "handler-refresh-secret",
			AccessTTL:     time.Hour,
			RefreshTTL:    24 * time.Hour,
		},
		Cookie: app.CookieConfig{
			AccessName:  "access_token",
			RefreshName: "refresh_token",
			SameSite:    http.SameSiteLaxMode,
		},
		ImageStore: app.ImageStoreConfig{
			Dir:      uploadDir,
			Route:    "/uploads/image",
			MaxBytes: 2 << 20,
		},
		DocumentStore: app.DocumentStoreConfig{
			Dir:      uploadDir,
			Route:    "/uploads/document",
			MaxBytes: 2 << 20,
		},
	}

	testApp, err := app.Build(context.Background(), cfg, dbURL, bcrypt.NewHasher(4), logging.NewLogger("test"), app.BuildDeleteFileModule(uploadDir))
	require.NoError(t, err)
	t.Cleanup(testApp.Infra.Pool.Close)
	return testApp
}

func AuthCookie(t *testing.T, testApp *app.App, fixtures *repotestutil.Fixtures, id user.ID, role user.Role, username string) *http.Cookie {
	t.Helper()

	fixtures.CreateSession(id, role, time.Now().Add(time.Hour), nil)
	token, err := testApp.Tokens.AccessTokenSvc.GenerateAccessToken(id, role, username, time.Hour)
	require.NoError(t, err)

	return &http.Cookie{
		Name:  "access_token",
		Value: token,
	}
}

func DoJSON(t *testing.T, handler http.Handler, method, path string, payload any, cookie *http.Cookie) *httptest.ResponseRecorder {
	t.Helper()

	var body io.Reader = bytes.NewReader(nil)
	if payload != nil {
		raw, err := json.Marshal(payload)
		require.NoError(t, err)
		body = bytes.NewReader(raw)
	}

	return DoRaw(t, handler, method, path, body, "application/json", cookie)
}

func DoRaw(t *testing.T, handler http.Handler, method, path string, body io.Reader, contentType string, cookie *http.Cookie) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, body)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if cookie != nil {
		req.AddCookie(cookie)
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func DoMultipart(t *testing.T, handler http.Handler, method, path string, fields map[string]string, cookie *http.Cookie) *httptest.ResponseRecorder {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for key, value := range fields {
		require.NoError(t, writer.WriteField(key, value))
	}
	require.NoError(t, writer.Close())

	return DoRaw(t, handler, method, path, &body, writer.FormDataContentType(), cookie)
}

func DecodeEnvelope(t *testing.T, rec *httptest.ResponseRecorder) Envelope {
	t.Helper()

	var env Envelope
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &env))
	return env
}

func AssertSuccess(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()

	env := DecodeEnvelope(t, rec)
	assert.Nil(t, env.Error)
	assert.Equal(t, "Success", env.Message)
}
