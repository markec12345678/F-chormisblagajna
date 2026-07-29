package report

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/report/handlers"

	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
)

type ReportModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *ReportModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Report module started")
		return nil
	}
}

func (m *ReportModule) OnEnd() func() {
	return func() {}
}

func (m *ReportModule) RegisterHttpHandlers(router *mux.Router, prefix string) {
	var auth_svc auth_mw.IAuthService

	if m.Config.Auth.Enabled {
		jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
		auth_svc = auth_mw.NewInternalAuth(m.Config, jwtUtil)
	} else {
		auth_svc = auth_mw.NewNoAuth(m.Config)
	}

	router.Handle(prefix+"/api/reports/sales", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetSalesReport(m.Config, m.Logger), "admin"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/reports/inventory", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetInventoryReport(m.Config, m.Logger), "admin"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/reports/dashboard", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetDashboardStats(m.Config, m.Logger), "admin", "cashier"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/reports/sales/export", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.ExportSalesReportCSV(m.Config, m.Logger), "admin"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/reports/inventory/export", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.ExportInventoryReportCSV(m.Config, m.Logger), "admin"),
	)).Methods("GET", "OPTIONS")
}
