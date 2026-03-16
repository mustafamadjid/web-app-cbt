package httpx

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	attemptujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeExpireAttemptHandlerRepo struct {
	updateFn     func(ctx context.Context, idAttempt ujian.ID, patch updatepatch.UpdateAttemptUjianPatch) error
	updateCalled bool
	gotAttemptID ujian.ID
	gotPatch     updatepatch.UpdateAttemptUjianPatch
}

func (f *fakeExpireAttemptHandlerRepo) UpdateAttemptUjian(ctx context.Context, idAttempt ujian.ID, patch updatepatch.UpdateAttemptUjianPatch) error {
	f.updateCalled = true
	f.gotAttemptID = idAttempt
	f.gotPatch = patch
	if f.updateFn != nil {
		return f.updateFn(ctx, idAttempt, patch)
	}
	return nil
}

type expireAttemptAPIResp struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
	Error   *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func TestExpireAttemptUjianHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		method        string
		idAttempt     string
		withActor     bool
		repo          *fakeExpireAttemptHandlerRepo
		wantStatus    int
		wantCode      string
		wantMessage   string
		wantUpdate    bool
		assertPayload func(t *testing.T, repo *fakeExpireAttemptHandlerRepo)
	}{
		{
			name:        "wrong method",
			method:      http.MethodGet,
			idAttempt:   "31",
			repo:        &fakeExpireAttemptHandlerRepo{},
			wantStatus:  http.StatusMethodNotAllowed,
			wantCode:    "METHOD_NOT_ALLOWED",
			wantMessage: "method not allowed",
		},
		{
			name:        "missing actor",
			method:      http.MethodPatch,
			idAttempt:   "31",
			repo:        &fakeExpireAttemptHandlerRepo{},
			wantStatus:  http.StatusInternalServerError,
			wantCode:    "INTERNAL_SERVER_ERROR",
			wantMessage: "internal server error : failed get actor from context",
		},
		{
			name:        "invalid id attempt",
			method:      http.MethodPatch,
			idAttempt:   "abc",
			withActor:   true,
			repo:        &fakeExpireAttemptHandlerRepo{},
			wantStatus:  http.StatusBadRequest,
			wantCode:    "BAD_REQUEST",
			wantMessage: "bad request: invalid id attempt",
		},
		{
			name:      "not found",
			method:    http.MethodPatch,
			idAttempt: "31",
			withActor: true,
			repo: &fakeExpireAttemptHandlerRepo{
				updateFn: func(context.Context, ujian.ID, updatepatch.UpdateAttemptUjianPatch) error {
					return coreerror.ErrNotFound
				},
			},
			wantStatus:  http.StatusNotFound,
			wantCode:    "NOT_FOUND",
			wantMessage: "data not found",
			wantUpdate:  true,
		},
		{
			name:        "success",
			method:      http.MethodPatch,
			idAttempt:   "31",
			withActor:   true,
			repo:        &fakeExpireAttemptHandlerRepo{},
			wantStatus:  http.StatusOK,
			wantMessage: "Success",
			wantUpdate:  true,
			assertPayload: func(t *testing.T, repo *fakeExpireAttemptHandlerRepo) {
				t.Helper()
				assert.Equal(t, ujian.ID(31), repo.gotAttemptID)
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
			handler := NewExpireAttemptUjianHandler(attemptujian_service.NewExpireAttemptUjianService(updater))

			req := httptest.NewRequest(tc.method, "/admin/ujian/attempt/"+tc.idAttempt+"/expire", nil)
			if tc.withActor {
				req = req.WithContext(middleware.WithActor(req.Context(), user.Actor{IdPengguna: 1, Role: user.ADMIN}))
			}

			rec := httptest.NewRecorder()
			handler.ExpireAttemptUjian(rec, req, httprouter.Params{{Key: "idAttempt", Value: tc.idAttempt}})

			assert.Equal(t, tc.wantStatus, rec.Code)

			var resp expireAttemptAPIResp
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			if tc.wantCode == "" {
				assert.Equal(t, tc.wantMessage, resp.Message)
				assert.Nil(t, resp.Error)
			} else {
				require.NotNil(t, resp.Error)
				assert.Equal(t, tc.wantCode, resp.Error.Code)
				assert.Equal(t, tc.wantMessage, resp.Error.Message)
			}

			assert.Equal(t, tc.wantUpdate, tc.repo.updateCalled)
			if tc.assertPayload != nil {
				tc.assertPayload(t, tc.repo)
			}
		})
	}
}
