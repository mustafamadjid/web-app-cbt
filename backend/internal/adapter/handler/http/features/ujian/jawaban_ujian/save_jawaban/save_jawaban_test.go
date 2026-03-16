package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	savejawaban_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/jawaban_ujian/save_jawaban"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeJawabanRepo struct {
	saveFn     func(ctx context.Context, idAttempt ujian.ID, jawaban []ujian.JawabanUjian) error
	saveCalled bool
	gotAttempt ujian.ID
	gotJawaban []ujian.JawabanUjian
}

func (f *fakeJawabanRepo) GetJawabanUjianByAttemptId(context.Context, ujian.ID) ([]ujian.JawabanUjian, error) {
	return nil, nil
}

func (f *fakeJawabanRepo) SaveJawabanUjian(ctx context.Context, idAttempt ujian.ID, jawaban []ujian.JawabanUjian) error {
	f.saveCalled = true
	f.gotAttempt = idAttempt
	f.gotJawaban = jawaban
	if f.saveFn != nil {
		return f.saveFn(ctx, idAttempt, jawaban)
	}
	return nil
}

type apiErrResp struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
	Error   *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func TestSaveJawabanUjianHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		method        string
		contentType   string
		body          string
		repo          *fakeJawabanRepo
		wantStatus    int
		wantCode      string
		wantMessage   string
		wantSave      bool
		assertPayload func(t *testing.T, repo *fakeJawabanRepo)
	}{
		{
			name:        "wrong method",
			method:      http.MethodGet,
			wantStatus:  http.StatusMethodNotAllowed,
			wantCode:    "METHOD_NOT_ALLOWED",
			wantMessage: "method not allowed",
			repo:        &fakeJawabanRepo{},
		},
		{
			name:        "invalid content type",
			method:      http.MethodPost,
			contentType: "text/plain",
			body:        `{"id_attempt":1,"jawaban":[]}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "BAD_REQUEST",
			wantMessage: "bad request: content type must be application/json",
			repo:        &fakeJawabanRepo{},
		},
		{
			name:        "invalid request body",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"id_attempt":1`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "BAD_REQUEST",
			wantMessage: "bad request: invalid request body",
			repo:        &fakeJawabanRepo{},
		},
		{
			name:        "validation error",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"id_attempt":1,"jawaban":[{"id_pilihan":2}]}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "BAD_REQUEST",
			wantMessage: "bad request: invalid jawaban ujian payload",
			repo:        &fakeJawabanRepo{},
		},
		{
			name:        "success",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"id_attempt":3,"jawaban":[{"id_soal":11,"jawaban_essay":"  essay siswa  "}]}`,
			wantStatus:  http.StatusOK,
			wantMessage: "Success",
			wantSave:    true,
			repo:        &fakeJawabanRepo{},
			assertPayload: func(t *testing.T, repo *fakeJawabanRepo) {
				t.Helper()
				require.Len(t, repo.gotJawaban, 1)
				assert.Equal(t, ujian.ID(3), repo.gotAttempt)
				assert.Equal(t, ujian.ID(0), repo.gotJawaban[0].IdJawaban)
				assert.Equal(t, ujian.ID(11), repo.gotJawaban[0].IdSoal)
				require.NotNil(t, repo.gotJawaban[0].JawabanEssay)
				assert.Equal(t, "essay siswa", *repo.gotJawaban[0].JawabanEssay)
			},
		},
		{
			name:        "repo error becomes internal error",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"id_attempt":3,"jawaban":[{"id_soal":11,"id_pilihan":7}]}`,
			wantStatus:  http.StatusInternalServerError,
			wantCode:    "INTERNAL_SERVER_ERROR",
			wantMessage: "internal server error",
			wantSave:    true,
			repo: &fakeJawabanRepo{
				saveFn: func(context.Context, ujian.ID, []ujian.JawabanUjian) error {
					return errors.New("repo error")
				},
			},
		},
		{
			name:        "repo invalid input stays bad request",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"id_attempt":3,"jawaban":[{"id_soal":11,"id_pilihan":7}]}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "BAD_REQUEST",
			wantMessage: "bad request: invalid jawaban ujian payload",
			wantSave:    true,
			repo: &fakeJawabanRepo{
				saveFn: func(context.Context, ujian.ID, []ujian.JawabanUjian) error {
					return coreerror.ErrInvalidInput
				},
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			handler := NewSaveJawabanUjianHandler(savejawaban_service.NewJawabanUjianService(tc.repo))
			req := httptest.NewRequest(tc.method, "/siswa/ujian/jawaban", bytes.NewBufferString(tc.body))
			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}

			rec := httptest.NewRecorder()
			handler.SaveJawabanUjian(rec, req, httprouter.Params{})

			assert.Equal(t, tc.wantStatus, rec.Code)

			var resp apiErrResp
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			if tc.wantCode == "" {
				assert.Equal(t, tc.wantMessage, resp.Message)
				assert.Nil(t, resp.Error)
			} else {
				require.NotNil(t, resp.Error)
				assert.Equal(t, tc.wantCode, resp.Error.Code)
				assert.Equal(t, tc.wantMessage, resp.Error.Message)
			}
			assert.Equal(t, tc.wantSave, tc.repo.saveCalled)

			if tc.assertPayload != nil {
				tc.assertPayload(t, tc.repo)
			}
		})
	}
}
