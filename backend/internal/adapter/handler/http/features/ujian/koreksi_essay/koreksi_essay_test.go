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
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	essaygrading_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/grading_ujian/grading/essay_grading"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeEssayGradingRepo struct {
	gradeFn     func(ctx context.Context, jawaban []ujian.JawabanUjian, gradedBy ujian.ID) error
	gradeCalled bool
	gotJawaban  []ujian.JawabanUjian
	gotGradedBy ujian.ID
}

func (f *fakeEssayGradingRepo) UpsertNilaiToHasilUjian(context.Context, float64, ujian.HasilUjian) error {
	return nil
}

func (f *fakeEssayGradingRepo) UpsertJawabanBenarToStatistikSoal(context.Context, []ujian.StatistikSoal) error {
	return nil
}

func (f *fakeEssayGradingRepo) UpsertJawabanSalahToStatistikSoal(context.Context, []ujian.StatistikSoal) error {
	return nil
}

func (f *fakeEssayGradingRepo) UpdateAndGradingEssayUjian(ctx context.Context, jawaban []ujian.JawabanUjian, gradedBy ujian.ID) error {
	f.gradeCalled = true
	f.gotJawaban = jawaban
	f.gotGradedBy = gradedBy
	if f.gradeFn != nil {
		return f.gradeFn(ctx, jawaban, gradedBy)
	}
	return nil
}

type koreksiEssayAPIResp struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
	Error   *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func TestKoreksiEssayHandler(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("repo error")

	tests := []struct {
		name          string
		method        string
		contentType   string
		body          string
		withActor     bool
		repo          *fakeEssayGradingRepo
		wantStatus    int
		wantCode      string
		wantMessage   string
		wantGrade     bool
		assertPayload func(t *testing.T, repo *fakeEssayGradingRepo)
	}{
		{
			name:        "wrong method",
			method:      http.MethodPost,
			wantStatus:  http.StatusMethodNotAllowed,
			wantCode:    "METHOD_NOT_ALLOWED",
			wantMessage: "method not allowed",
			repo:        &fakeEssayGradingRepo{},
		},
		{
			name:        "missing actor",
			method:      http.MethodPatch,
			contentType: "application/json",
			body:        `{"jawaban":[{"id_jawaban":11,"essay_is_benar":true}]}`,
			wantStatus:  http.StatusInternalServerError,
			wantCode:    "INTERNAL_SERVER_ERROR",
			wantMessage: "internal server error : failed get actor from context",
			repo:        &fakeEssayGradingRepo{},
		},
		{
			name:        "invalid content type",
			method:      http.MethodPatch,
			contentType: "text/plain",
			body:        `{"jawaban":[{"id_jawaban":11,"essay_is_benar":true}]}`,
			withActor:   true,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "BAD_REQUEST",
			wantMessage: "bad request: content type must be application/json",
			repo:        &fakeEssayGradingRepo{},
		},
		{
			name:        "invalid request body",
			method:      http.MethodPatch,
			contentType: "application/json",
			body:        `{"jawaban":[{"id_jawaban":11,"essay_is_benar":true}]`,
			withActor:   true,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "BAD_REQUEST",
			wantMessage: "bad request: invalid request body",
			repo:        &fakeEssayGradingRepo{},
		},
		{
			name:        "unknown field rejected",
			method:      http.MethodPatch,
			contentType: "application/json",
			body:        `{"jawaban":[{"id_jawaban":11,"essay_is_benar":true,"unknown":"x"}]}`,
			withActor:   true,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "BAD_REQUEST",
			wantMessage: "bad request: invalid request body",
			repo:        &fakeEssayGradingRepo{},
		},
		{
			name:        "missing essay grading field",
			method:      http.MethodPatch,
			contentType: "application/json",
			body:        `{"jawaban":[{"id_jawaban":11}]}`,
			withActor:   true,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "BAD_REQUEST",
			wantMessage: "bad request: invalid essay grading payload",
			repo:        &fakeEssayGradingRepo{},
		},
		{
			name:        "repo invalid input stays bad request",
			method:      http.MethodPatch,
			contentType: "application/json",
			body:        `{"jawaban":[{"id_jawaban":11,"essay_is_benar":true}]}`,
			withActor:   true,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "BAD_REQUEST",
			wantMessage: "bad request: invalid essay grading payload",
			wantGrade:   true,
			repo: &fakeEssayGradingRepo{
				gradeFn: func(context.Context, []ujian.JawabanUjian, ujian.ID) error {
					return coreerror.ErrInvalidInput
				},
			},
		},
		{
			name:        "repo not found becomes not found",
			method:      http.MethodPatch,
			contentType: "application/json",
			body:        `{"jawaban":[{"id_jawaban":11,"essay_is_benar":true}]}`,
			withActor:   true,
			wantStatus:  http.StatusNotFound,
			wantCode:    "NOT_FOUND",
			wantMessage: "data not found",
			wantGrade:   true,
			repo: &fakeEssayGradingRepo{
				gradeFn: func(context.Context, []ujian.JawabanUjian, ujian.ID) error {
					return coreerror.ErrNotFound
				},
			},
		},
		{
			name:        "repo conflict becomes conflict",
			method:      http.MethodPatch,
			contentType: "application/json",
			body:        `{"jawaban":[{"id_jawaban":11,"essay_is_benar":true}]}`,
			withActor:   true,
			wantStatus:  http.StatusConflict,
			wantCode:    "CONFLICT",
			wantMessage: "conflict",
			wantGrade:   true,
			repo: &fakeEssayGradingRepo{
				gradeFn: func(context.Context, []ujian.JawabanUjian, ujian.ID) error {
					return coreerror.ErrConflict
				},
			},
		},
		{
			name:        "repo error becomes internal error",
			method:      http.MethodPatch,
			contentType: "application/json",
			body:        `{"jawaban":[{"id_jawaban":11,"essay_is_benar":true}]}`,
			withActor:   true,
			wantStatus:  http.StatusInternalServerError,
			wantCode:    "INTERNAL_SERVER_ERROR",
			wantMessage: "internal server error",
			wantGrade:   true,
			repo: &fakeEssayGradingRepo{
				gradeFn: func(context.Context, []ujian.JawabanUjian, ujian.ID) error {
					return repoErr
				},
			},
		},
		{
			name:        "success",
			method:      http.MethodPatch,
			contentType: "application/json",
			body:        `{"jawaban":[{"id_jawaban":11,"essay_is_benar":true},{"id_jawaban":12,"essay_is_benar":false}]}`,
			withActor:   true,
			wantStatus:  http.StatusOK,
			wantMessage: "Success",
			wantGrade:   true,
			repo:        &fakeEssayGradingRepo{},
			assertPayload: func(t *testing.T, repo *fakeEssayGradingRepo) {
				t.Helper()
				require.Len(t, repo.gotJawaban, 2)
				assert.Equal(t, ujian.ID(7), repo.gotGradedBy)
				assert.Equal(t, ujian.ID(11), repo.gotJawaban[0].IdJawaban)
				require.NotNil(t, repo.gotJawaban[0].EssayIsBenar)
				assert.True(t, *repo.gotJawaban[0].EssayIsBenar)
				assert.Equal(t, ujian.ID(12), repo.gotJawaban[1].IdJawaban)
				require.NotNil(t, repo.gotJawaban[1].EssayIsBenar)
				assert.False(t, *repo.gotJawaban[1].EssayIsBenar)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			handler := NewKoreksiEssayHandler(essaygrading_service.NewEssayGradingUjianService(tc.repo))
			req := httptest.NewRequest(tc.method, "/ujian/koreksi-essay", bytes.NewBufferString(tc.body))
			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}
			if tc.withActor {
				req = req.WithContext(middleware.WithActor(req.Context(), user.Actor{IdPengguna: 7, Role: user.GURU}))
			}

			rec := httptest.NewRecorder()
			handler.KoreksiEssay(rec, req, httprouter.Params{})

			assert.Equal(t, tc.wantStatus, rec.Code)

			var resp koreksiEssayAPIResp
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			if tc.wantCode == "" {
				assert.Equal(t, tc.wantMessage, resp.Message)
				assert.Nil(t, resp.Error)
			} else {
				require.NotNil(t, resp.Error)
				assert.Equal(t, tc.wantCode, resp.Error.Code)
				assert.Equal(t, tc.wantMessage, resp.Error.Message)
			}

			assert.Equal(t, tc.wantGrade, tc.repo.gradeCalled)
			if tc.assertPayload != nil {
				tc.assertPayload(t, tc.repo)
			}
		})
	}
}
