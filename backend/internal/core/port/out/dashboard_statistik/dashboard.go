package dashboard_repo

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/dashboard"
)

type DashboardStatistikRepository interface {
	GetDashboardStatistik(ctx context.Context)(dashboard.DashboardStatistik, error)
}