package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	siswaujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeAttemptHandlerChecker struct {
	checkFn    func(ctx context.Context, idSiswa int, idJadwalUjian int) (bool, int, error)
	tokenFn    func(ctx context.Context, token string, idJadwalUjian int) (bool, error)
	ownerFn    func(ctx context.Context, idSiswa int, idAttempt ujian.ID) (bool, error)
	deadlineFn func(ctx context.Context, idJadwalUjian int) (time.Time, error)
}

func (f *fakeAttemptHandlerChecker) CheckValidSiswaInPesertaUjianById(ctx context.Context, idSiswa int, idJadwalUjian int) (bool, int, error) {
	if f.checkFn != nil {
		return f.checkFn(ctx, idSiswa, idJadwalUjian)
	}
	return false, 0, nil
}

func (f *fakeAttemptHandlerChecker) CheckTokenUjian(ctx context.Context, token string, idJadwalUjian int) (bool, error) {
	if f.tokenFn != nil {
		return f.tokenFn(ctx, token, idJadwalUjian)
	}
	return false, nil
}

func (f *fakeAttemptHandlerChecker) CheckAttemptOwnershipBySiswa(ctx context.Context, idSiswa int, idAttempt ujian.ID) (bool, error) {
	if f.ownerFn != nil {
		return f.ownerFn(ctx, idSiswa, idAttempt)
	}
	return false, nil
}

func (f *fakeAttemptHandlerChecker) GetDeadlineUjian(ctx context.Context, idJadwalUjian int) (time.Time, error) {
	if f.deadlineFn != nil {
		return f.deadlineFn(ctx, idJadwalUjian)
	}
	return time.Time{}, nil
}

type fakeAttemptHandlerRepo struct {
	createFn func(ctx context.Context, data ujian.AttemptUjian) error
}

func (f *fakeAttemptHandlerRepo) GetAttemptUjianById(context.Context, ujian.ID) (ujian.AttemptUjian, error) {
	return ujian.AttemptUjian{}, nil
}

func (f *fakeAttemptHandlerRepo) CreateAttemptUjian(ctx context.Context, data ujian.AttemptUjian) error {
	if f.createFn != nil {
		return f.createFn(ctx, data)
	}
	return nil
}

func (f *fakeAttemptHandlerRepo) UpdateAttemptUjian(context.Context, ujian.ID, updatepatch.UpdateAttemptUjianPatch) error {
	return nil
}

func (f *fakeAttemptHandlerRepo) DeleteAttemptUjian(context.Context, ujian.ID) error {
	return nil
}

type createAttemptAPIResp struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
	Error   *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func TestAttemptUjianHandler(t *testing.T) {
	t.Parallel()

	deadline := time.Date(2026, time.March, 16, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		repo        *fakeAttemptHandlerRepo
		wantStatus  int
		wantCode    string
		wantMessage string
	}{
		{
			name: "double attempt not allowed",
			repo: &fakeAttemptHandlerRepo{
				createFn: func(context.Context, ujian.AttemptUjian) error {
					return coreerror.ErrSiswaHasActiveAttempt
				},
			},
			wantStatus:  http.StatusConflict,
			wantCode:    "DOUBLE_ATTEMPT_NOT_ALLOWED",
			wantMessage: "double attempt not allowed",
		},
		{
			name:        "success",
			repo:        &fakeAttemptHandlerRepo{},
			wantStatus:  http.StatusOK,
			wantMessage: "Success",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			handler := NewAttemptUjianHandler(siswaujian_service.NewAttemptUjianService(
				&fakeAttemptHandlerChecker{
					checkFn: func(context.Context, int, int) (bool, int, error) {
						return true, 17, nil
					},
					tokenFn: func(context.Context, string, int) (bool, error) {
						return true, nil
					},
					deadlineFn: func(context.Context, int) (time.Time, error) {
						return deadline, nil
					},
				},
				tc.repo,
			))

			req := httptest.NewRequest(http.MethodPost, "/siswa/ujian/attempt", bytes.NewBufferString(`{"id_siswa":9,"id_jadwal_ujian":21,"token_ujian":"token-123","waktu_mulai":"2026-03-16T11:30:00Z"}`))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			handler.AttemptUjian(rec, req, httprouter.Params{})

			assert.Equal(t, tc.wantStatus, rec.Code)

			var resp createAttemptAPIResp
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			if tc.wantCode == "" {
				assert.Equal(t, tc.wantMessage, resp.Message)
				assert.Nil(t, resp.Error)
				return
			}

			require.NotNil(t, resp.Error)
			assert.Equal(t, tc.wantCode, resp.Error.Code)
			assert.Equal(t, tc.wantMessage, resp.Error.Message)
		})
	}
}
