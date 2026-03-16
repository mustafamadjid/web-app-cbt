package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	attemptujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeUpdateAttemptHandlerRepo struct {
	updateFn     func(ctx context.Context, idAttempt ujian.ID, patch updatepatch.UpdateAttemptUjianPatch) error
	updateCalled bool
	gotAttemptID ujian.ID
	gotPatch     updatepatch.UpdateAttemptUjianPatch
}

func (f *fakeUpdateAttemptHandlerRepo) UpdateAttemptUjian(ctx context.Context, idAttempt ujian.ID, patch updatepatch.UpdateAttemptUjianPatch) error {
	f.updateCalled = true
	f.gotAttemptID = idAttempt
	f.gotPatch = patch
	if f.updateFn != nil {
		return f.updateFn(ctx, idAttempt, patch)
	}
	return nil
}

type fakeUpdateAttemptOwnershipChecker struct {
	checkFn      func(ctx context.Context, idSiswa int, idAttempt ujian.ID) (bool, error)
	checkCalled  bool
	gotSiswaID   int
	gotAttemptID ujian.ID
}

func (f *fakeUpdateAttemptOwnershipChecker) CheckAttemptOwnershipBySiswa(ctx context.Context, idSiswa int, idAttempt ujian.ID) (bool, error) {
	f.checkCalled = true
	f.gotSiswaID = idSiswa
	f.gotAttemptID = idAttempt
	if f.checkFn != nil {
		return f.checkFn(ctx, idSiswa, idAttempt)
	}
	return false, nil
}

type updateAttemptAPIResp struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
	Error   *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func TestUpdateAttemptUjianHandler(t *testing.T) {
	t.Parallel()

	submitTime := time.Date(2026, time.March, 16, 11, 45, 0, 0, time.UTC)
	repoErr := errors.New("repo error")

	tests := []struct {
		name          string
		method        string
		contentType   string
		body          string
		idAttempt     string
		withActor     bool
		checker       *fakeUpdateAttemptOwnershipChecker
		repo          *fakeUpdateAttemptHandlerRepo
		wantStatus    int
		wantCode      string
		wantMessage   string
		wantCheck     bool
		wantUpdate    bool
		assertPayload func(t *testing.T, repo *fakeUpdateAttemptHandlerRepo)
	}{
		{
			name:        "wrong method",
			method:      http.MethodGet,
			idAttempt:   "17",
			wantStatus:  http.StatusMethodNotAllowed,
			wantCode:    "METHOD_NOT_ALLOWED",
			wantMessage: "method not allowed",
			checker:     &fakeUpdateAttemptOwnershipChecker{},
			repo:        &fakeUpdateAttemptHandlerRepo{},
		},
		{
			name:        "missing actor",
			method:      http.MethodPatch,
			idAttempt:   "17",
			contentType: "application/json",
			body:        `{"status_attempt":"submitted","waktu_submit":"2026-03-16T11:45:00Z"}`,
			wantStatus:  http.StatusInternalServerError,
			wantCode:    "INTERNAL_SERVER_ERROR",
			wantMessage: "internal server error : failed get actor from context",
			checker:     &fakeUpdateAttemptOwnershipChecker{},
			repo:        &fakeUpdateAttemptHandlerRepo{},
		},
		{
			name:        "invalid id attempt",
			method:      http.MethodPatch,
			idAttempt:   "abc",
			withActor:   true,
			contentType: "application/json",
			body:        `{"status_attempt":"submitted","waktu_submit":"2026-03-16T11:45:00Z"}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "BAD_REQUEST",
			wantMessage: "bad request: invalid id attempt",
			checker:     &fakeUpdateAttemptOwnershipChecker{},
			repo:        &fakeUpdateAttemptHandlerRepo{},
		},
		{
			name:        "invalid content type",
			method:      http.MethodPatch,
			idAttempt:   "17",
			withActor:   true,
			contentType: "text/plain",
			body:        `{"status_attempt":"submitted","waktu_submit":"2026-03-16T11:45:00Z"}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "BAD_REQUEST",
			wantMessage: "bad request: content type must be application/json",
			checker:     &fakeUpdateAttemptOwnershipChecker{},
			repo:        &fakeUpdateAttemptHandlerRepo{},
		},
		{
			name:        "invalid request body",
			method:      http.MethodPatch,
			idAttempt:   "17",
			withActor:   true,
			contentType: "application/json",
			body:        `{"status_attempt":"submitted"`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "BAD_REQUEST",
			wantMessage: "bad request: invalid request body",
			checker:     &fakeUpdateAttemptOwnershipChecker{},
			repo:        &fakeUpdateAttemptHandlerRepo{},
		},
		{
			name:        "unknown field rejected",
			method:      http.MethodPatch,
			idAttempt:   "17",
			withActor:   true,
			contentType: "application/json",
			body:        `{"status_attempt":"submitted","updated_at":"2026-03-16T11:45:00Z"}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "BAD_REQUEST",
			wantMessage: "bad request: invalid request body",
			checker:     &fakeUpdateAttemptOwnershipChecker{},
			repo:        &fakeUpdateAttemptHandlerRepo{},
		},
		{
			name:        "ownership check returns not found",
			method:      http.MethodPatch,
			idAttempt:   "17",
			withActor:   true,
			contentType: "application/json",
			body:        `{"status_attempt":"expired"}`,
			checker: &fakeUpdateAttemptOwnershipChecker{
				checkFn: func(context.Context, int, ujian.ID) (bool, error) {
					return false, nil
				},
			},
			repo:        &fakeUpdateAttemptHandlerRepo{},
			wantStatus:  http.StatusNotFound,
			wantCode:    "NOT_FOUND",
			wantMessage: "data not found",
			wantCheck:   true,
		},
		{
			name:        "submitted requires submit time",
			method:      http.MethodPatch,
			idAttempt:   "17",
			withActor:   true,
			contentType: "application/json",
			body:        `{"status_attempt":"submitted"}`,
			checker:     &fakeUpdateAttemptOwnershipChecker{},
			repo:        &fakeUpdateAttemptHandlerRepo{},
			wantStatus:  http.StatusBadRequest,
			wantCode:    "BAD_REQUEST",
			wantMessage: "bad request: invalid attempt update payload",
		},
		{
			name:        "expired rejects submit time",
			method:      http.MethodPatch,
			idAttempt:   "17",
			withActor:   true,
			contentType: "application/json",
			body:        `{"status_attempt":"expired","waktu_submit":"2026-03-16T11:45:00Z"}`,
			checker:     &fakeUpdateAttemptOwnershipChecker{},
			repo:        &fakeUpdateAttemptHandlerRepo{},
			wantStatus:  http.StatusBadRequest,
			wantCode:    "BAD_REQUEST",
			wantMessage: "bad request: invalid attempt update payload",
		},
		{
			name:        "invalid status",
			method:      http.MethodPatch,
			idAttempt:   "17",
			withActor:   true,
			contentType: "application/json",
			body:        `{"status_attempt":"cancelled"}`,
			checker:     &fakeUpdateAttemptOwnershipChecker{},
			repo:        &fakeUpdateAttemptHandlerRepo{},
			wantStatus:  http.StatusBadRequest,
			wantCode:    "BAD_REQUEST",
			wantMessage: "bad request: invalid attempt update payload",
		},
		{
			name:        "repo error becomes internal server error",
			method:      http.MethodPatch,
			idAttempt:   "17",
			withActor:   true,
			contentType: "application/json",
			body:        `{"status_attempt":"submitted","waktu_submit":"2026-03-16T11:45:00Z"}`,
			checker: &fakeUpdateAttemptOwnershipChecker{
				checkFn: func(context.Context, int, ujian.ID) (bool, error) {
					return true, nil
				},
			},
			repo: &fakeUpdateAttemptHandlerRepo{
				updateFn: func(context.Context, ujian.ID, updatepatch.UpdateAttemptUjianPatch) error {
					return repoErr
				},
			},
			wantStatus:  http.StatusInternalServerError,
			wantCode:    "INTERNAL_SERVER_ERROR",
			wantMessage: "internal server error",
			wantCheck:   true,
			wantUpdate:  true,
		},
		{
			name:        "success submitted",
			method:      http.MethodPatch,
			idAttempt:   "17",
			withActor:   true,
			contentType: "application/json",
			body:        `{"status_attempt":" submitted ","waktu_submit":"2026-03-16T11:45:00Z"}`,
			checker: &fakeUpdateAttemptOwnershipChecker{
				checkFn: func(context.Context, int, ujian.ID) (bool, error) {
					return true, nil
				},
			},
			repo:        &fakeUpdateAttemptHandlerRepo{},
			wantStatus:  http.StatusOK,
			wantMessage: "Success",
			wantCheck:   true,
			wantUpdate:  true,
			assertPayload: func(t *testing.T, repo *fakeUpdateAttemptHandlerRepo) {
				t.Helper()
				assert.Equal(t, ujian.ID(17), repo.gotAttemptID)
				assert.Equal(t, &submitTime, repo.gotPatch.WaktuSubmit)
				if assert.NotNil(t, repo.gotPatch.StatusAttempt) {
					assert.Equal(t, ujian.ATTEMPT_SUBMITTED, *repo.gotPatch.StatusAttempt)
				}
			},
		},
		{
			name:        "success expired",
			method:      http.MethodPatch,
			idAttempt:   "17",
			withActor:   true,
			contentType: "application/json",
			body:        `{"status_attempt":"expired"}`,
			checker: &fakeUpdateAttemptOwnershipChecker{
				checkFn: func(context.Context, int, ujian.ID) (bool, error) {
					return true, nil
				},
			},
			repo:        &fakeUpdateAttemptHandlerRepo{},
			wantStatus:  http.StatusOK,
			wantMessage: "Success",
			wantCheck:   true,
			wantUpdate:  true,
			assertPayload: func(t *testing.T, repo *fakeUpdateAttemptHandlerRepo) {
				t.Helper()
				assert.Equal(t, ujian.ID(17), repo.gotAttemptID)
				assert.Nil(t, repo.gotPatch.WaktuSubmit)
				if assert.NotNil(t, repo.gotPatch.StatusAttempt) {
					assert.Equal(t, ujian.ATTEMPT_EXPIRED, *repo.gotPatch.StatusAttempt)
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			updater := attemptujian_service.NewUpdateAttemptUjianService(tc.repo)
			svc := attemptujian_service.NewSiswaUpdateAttemptUjianService(tc.checker, updater)
			handler := NewUpdateAttemptUjianHandler(svc)

			req := httptest.NewRequest(tc.method, "/siswa/ujian/attempt/"+tc.idAttempt, bytes.NewBufferString(tc.body))
			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}
			if tc.withActor {
				req = req.WithContext(middleware.WithActor(req.Context(), user.Actor{IdPengguna: 9, Role: user.SISWA}))
			}

			rec := httptest.NewRecorder()
			handler.UpdateAttemptUjian(rec, req, httprouter.Params{{Key: "idAttempt", Value: tc.idAttempt}})

			assert.Equal(t, tc.wantStatus, rec.Code)

			var resp updateAttemptAPIResp
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			if tc.wantCode == "" {
				assert.Equal(t, tc.wantMessage, resp.Message)
				assert.Nil(t, resp.Error)
			} else {
				require.NotNil(t, resp.Error)
				assert.Equal(t, tc.wantCode, resp.Error.Code)
				assert.Equal(t, tc.wantMessage, resp.Error.Message)
			}

			assert.Equal(t, tc.wantCheck, tc.checker.checkCalled)
			assert.Equal(t, tc.wantUpdate, tc.repo.updateCalled)
			if tc.assertPayload != nil {
				tc.assertPayload(t, tc.repo)
			}
		})
	}
}
