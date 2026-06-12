package httpx_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	dashboardget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/dashboard/get"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/dashboard"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	dashboard_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/dashboard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeDashboardStatistikRepo struct {
	result      dashboard.DashboardStatistik
	err         error
	called      bool
	receivedCtx context.Context
}

func (f *fakeDashboardStatistikRepo) GetDashboardStatistik(ctx context.Context) (dashboard.DashboardStatistik, error) {
	f.called = true
	f.receivedCtx = ctx
	return f.result, f.err
}

type recordingLogger struct {
	errorCalls int
	lastMsg    string
	lastAttrs  []any
}

func (l *recordingLogger) With(_ ...any) corelog.Logger {
	return l
}

func (l *recordingLogger) Info(_ context.Context, _ string, _ ...any) {}

func (l *recordingLogger) Error(_ context.Context, msg string, attrs ...any) {
	l.errorCalls++
	l.lastMsg = msg
	l.lastAttrs = attrs
}

type successEnvelope struct {
	Data    dashboardget.GetDashboardStatistikResponse `json:"data"`
	Message string                                     `json:"message"`
	Error   any                                        `json:"error"`
}

type errorEnvelope struct {
	Data  bool `json:"data"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func TestGetDashboardStatistikHandler(t *testing.T) {
	t.Parallel()

	expected := dashboard.DashboardStatistik{
		TotalSiswa:           90,
		TotalGuru:            12,
		TotalUjianTerlaksana: 25,
		TotalBankSoal:        40,
		TotalMapelAktif:      6,
	}
	repoErr := errors.New("repo error")

	tests := []struct {
		name             string
		method           string
		repo             *fakeDashboardStatistikRepo
		wantStatus       int
		wantMessage      string
		wantErrorCode    string
		wantErrorMessage string
		wantData         dashboardget.GetDashboardStatistikResponse
		wantRepoCalled   bool
		wantLoggerErrors int
	}{
		{
			name:             "Branch 1 -> invalid method returns method not allowed",
			method:           http.MethodPost,
			repo:             &fakeDashboardStatistikRepo{},
			wantStatus:       http.StatusMethodNotAllowed,
			wantErrorCode:    "METHOD_NOT_ALLOWED",
			wantErrorMessage: "method not allowed",
		},
		{
			name:             "Branch 2 -> service error returns internal server error",
			method:           http.MethodGet,
			repo:             &fakeDashboardStatistikRepo{err: repoErr},
			wantStatus:       http.StatusInternalServerError,
			wantErrorCode:    "INTERNAL_SERVER_ERROR",
			wantErrorMessage: "internal server error: failed get dashboard statistik",
			wantRepoCalled:   true,
			wantLoggerErrors: 2,
		},
		{
			name:        "Branch 3 -> success returns mapped dashboard statistik",
			method:      http.MethodGet,
			repo:        &fakeDashboardStatistikRepo{result: expected},
			wantStatus:  http.StatusOK,
			wantMessage: "Success",
			wantData: dashboardget.GetDashboardStatistikResponse{
				TotalSiswa:           expected.TotalSiswa,
				TotalGuru:            expected.TotalGuru,
				TotalUjianTerlaksana: expected.TotalUjianTerlaksana,
				TotalBankSoal:        expected.TotalBankSoal,
				TotalMapelAktif:      expected.TotalMapelAktif,
			},
			wantRepoCalled: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			logger := &recordingLogger{}
			ctx := corelog.WithLogger(context.Background(), logger)
			req := httptest.NewRequest(tc.method, "/dashboard", nil).WithContext(ctx)
			rec := httptest.NewRecorder()

			svc := dashboard_service.NewDashboardService(tc.repo)
			handler := dashboardget.NewGetDashboardStatistikHandler(svc)

			handler.GetDashboardStatistik(rec, req, nil)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
			assert.Equal(t, tc.wantRepoCalled, tc.repo.called)

			if tc.wantRepoCalled {
				assert.Same(t, ctx, tc.repo.receivedCtx)
			}

			if tc.wantErrorCode != "" {
				var body errorEnvelope
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
				assert.True(t, body.Data)
				assert.Equal(t, tc.wantErrorCode, body.Error.Code)
				assert.Equal(t, tc.wantErrorMessage, body.Error.Message)
				assert.Equal(t, tc.wantLoggerErrors, logger.errorCalls)
				if tc.wantLoggerErrors > 0 {
					assert.Equal(t, "failed get dashboard statistik", logger.lastMsg)
					assert.Contains(t, logger.lastAttrs, "adapter.http.handler")
				}
				return
			}

			var body successEnvelope
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			assert.Equal(t, tc.wantMessage, body.Message)
			assert.Nil(t, body.Error)
			assert.Equal(t, tc.wantData, body.Data)
			assert.Zero(t, logger.errorCalls)
		})
	}
}
