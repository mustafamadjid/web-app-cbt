package dashboard_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/dashboard"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	dashboard_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/dashboard"
	"github.com/stretchr/testify/assert"
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

func TestDashboardServiceGetDashboardStatistik(t *testing.T) {
	t.Parallel()

	expected := dashboard.DashboardStatistik{
		TotalSiswa:           120,
		TotalGuru:            18,
		TotalUjianTerlaksana: 42,
		TotalBankSoal:        65,
		TotalMapelAktif:      8,
	}
	repoErr := errors.New("repo error")

	tests := []struct {
		name           string
		repo           *fakeDashboardStatistikRepo
		want           dashboard.DashboardStatistik
		wantErr        error
		wantErrorCalls int
	}{
		{
			name: "Branch 1 -> repository success returns statistik",
			repo: &fakeDashboardStatistikRepo{
				result: expected,
			},
			want: expected,
		},
		{
			name: "Branch 2 -> repository error returns zero value and logs error",
			repo: &fakeDashboardStatistikRepo{
				err: repoErr,
			},
			want:           dashboard.DashboardStatistik{},
			wantErr:        repoErr,
			wantErrorCalls: 1,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			logger := &recordingLogger{}
			ctx := corelog.WithLogger(context.Background(), logger)
			svc := dashboard_service.NewDashboardService(tc.repo)

			got, err := svc.GetDashboardStatistik(ctx)

			assert.True(t, tc.repo.called)
			assert.Same(t, ctx, tc.repo.receivedCtx)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantErrorCalls, logger.errorCalls)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Equal(t, "failed get dashboard statistik", logger.lastMsg)
				assert.Contains(t, logger.lastAttrs, "layer")
				assert.Contains(t, logger.lastAttrs, "core.service")
				return
			}

			assert.NoError(t, err)
		})
	}
}
