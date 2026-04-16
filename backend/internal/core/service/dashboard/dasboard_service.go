package dashboard_service

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/dashboard"
	dashboard_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/dashboard_statistik"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type DashboardService struct {
	repo dashboard_repo.DashboardStatistikRepository
}

func NewDashboardService(repo dashboard_repo.DashboardStatistikRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func(r *DashboardService)GetDashboardStatistik(ctx context.Context)(dashboard.DashboardStatistik, error){
	logger := corelog.FromContext(ctx)

	statistik, err := r.repo.GetDashboardStatistik(ctx)
	if err != nil {
		logger.Error(ctx, "failed get dashboard statistik", "layer", "core.service", "op", "dashboard.statistik", "err", err)
		return dashboard.DashboardStatistik{}, err
	}

	return statistik, nil
}