package httpx_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	active_attempt "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/attempt_ujian/active_attempt"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
	activeattempt_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/active_attempt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeActiveAttemptRepo struct {
	getFn func(ctx context.Context, idSiswa int, idJadwalUjian int) (ujian.AttemptUjian, error)
}

func (f *fakeActiveAttemptRepo) GetActiveUjianAttemptBySiswa(ctx context.Context, idSiswa int, idJadwalUjian int) (ujian.AttemptUjian, error) {
	if f.getFn != nil {
		return f.getFn(ctx, idSiswa, idJadwalUjian)
	}
	return ujian.AttemptUjian{}, nil
}

func (*fakeActiveAttemptRepo) ListUjianSiswa(context.Context, int, query.ListUjianFilter) ([]ujian.ListUjian, error) {
	return nil, nil
}

func (*fakeActiveAttemptRepo) GetWaktuSelesaiUjian(context.Context, int) (time.Time, error) {
	return time.Time{}, nil
}

var _ ujian_repo.UjianSiswaRepository = (*fakeActiveAttemptRepo)(nil)

type activeAttemptAPIResp struct {
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
	Error   *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func TestGetActiveAttemptUjianHandler(t *testing.T) {
	t.Parallel()

	waktuMulai := time.Date(2026, time.March, 16, 11, 30, 0, 0, time.UTC)
	deadlineAt := time.Date(2026, time.March, 16, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		url          string
		withActor    bool
		repo         *fakeActiveAttemptRepo
		wantStatus   int
		wantCode     string
		wantMessage  string
		assertResult func(t *testing.T, data active_attempt.GetActiveAttemptUjianResponse)
	}{
		{
			name:        "missing actor",
			url:         "/siswa/ujian/attempt/active?id_jadwal_ujian=21",
			repo:        &fakeActiveAttemptRepo{},
			wantStatus:  http.StatusInternalServerError,
			wantCode:    "INTERNAL_SERVER_ERROR",
			wantMessage: "internal server error : failed get actor from context",
		},
		{
			name:        "invalid query",
			url:         "/siswa/ujian/attempt/active?id_jadwal_ujian=abc",
			withActor:   true,
			repo:        &fakeActiveAttemptRepo{},
			wantStatus:  http.StatusBadRequest,
			wantCode:    "BAD_REQUEST",
			wantMessage: "id_jadwal_ujian must be a positive number",
		},
		{
			name:      "active attempt not found",
			url:       "/siswa/ujian/attempt/active?id_jadwal_ujian=21",
			withActor: true,
			repo: &fakeActiveAttemptRepo{
				getFn: func(context.Context, int, int) (ujian.AttemptUjian, error) {
					return ujian.AttemptUjian{}, sql.ErrNoRows
				},
			},
			wantStatus:  http.StatusNotFound,
			wantCode:    "NOT_FOUND",
			wantMessage: "active attempt not found",
		},
		{
			name:      "success",
			url:       "/siswa/ujian/attempt/active?id_jadwal_ujian=21",
			withActor: true,
			repo: &fakeActiveAttemptRepo{
				getFn: func(context.Context, int, int) (ujian.AttemptUjian, error) {
					return ujian.AttemptUjian{
						IdAttempt:      31,
						IdPesertaUjian: 17,
						StatusAttempt:  ujian.ATTEMPT_IN_PROGRESS,
						WaktuMulai:     &waktuMulai,
						DeadlineAt:     &deadlineAt,
					}, nil
				},
			},
			wantStatus:  http.StatusOK,
			wantMessage: "Success",
			assertResult: func(t *testing.T, data active_attempt.GetActiveAttemptUjianResponse) {
				t.Helper()
				require.NotNil(t, data.WaktuMulai)
				require.NotNil(t, data.DeadlineAt)
				assert.Equal(t, 31, data.IDAttempt)
				assert.Equal(t, 17, data.IDPesertaUjian)
				assert.Equal(t, "in_progress", data.StatusAttempt)
				assert.Equal(t, waktuMulai.Format(time.RFC3339), *data.WaktuMulai)
				assert.Equal(t, deadlineAt.Format(time.RFC3339), *data.DeadlineAt)
				assert.Nil(t, data.WaktuSubmit)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			handler := active_attempt.NewGetActiveAttemptUjianHandler(activeattempt_service.NewGetActiveAttemptUjianService(tc.repo))
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			if tc.withActor {
				req = req.WithContext(middleware.WithActor(req.Context(), user.Actor{IdPengguna: 9, Role: user.SISWA}))
			}

			rec := httptest.NewRecorder()
			handler.GetActiveAttemptUjian(rec, req, httprouter.Params{})

			assert.Equal(t, tc.wantStatus, rec.Code)

			var resp activeAttemptAPIResp
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			if tc.wantCode == "" {
				assert.Equal(t, tc.wantMessage, resp.Message)
				assert.Nil(t, resp.Error)
				if tc.assertResult != nil {
					tc.assertResult(t, decodeActiveAttemptResponse(t, resp))
				}
				return
			}

			require.NotNil(t, resp.Error)
			assert.Equal(t, tc.wantCode, resp.Error.Code)
			assert.Equal(t, tc.wantMessage, resp.Error.Message)
		})
	}
}

func decodeActiveAttemptResponse(t *testing.T, resp activeAttemptAPIResp) active_attempt.GetActiveAttemptUjianResponse {
	t.Helper()

	var data active_attempt.GetActiveAttemptUjianResponse
	require.NoError(t, json.Unmarshal(resp.Data, &data))
	return data
}
