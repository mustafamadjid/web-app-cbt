package app

import (
	httpget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/dashboard/get"
	dashboard_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/dashboard"
)

type DashboardModule struct {
	Service    *dashboard_service.DashboardService
	GetHandler *httpget.GetDashboardStatistikHandler
}

func BuildDashboardModule(infra *InfraModule) *DashboardModule {
	svc := dashboard_service.NewDashboardService(infra.dashboardRepo)
	handler := httpget.NewGetDashboardStatistikHandler(svc)

	return &DashboardModule{
		Service:    svc,
		GetHandler: handler,
	}
}
