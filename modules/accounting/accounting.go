package accounting

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/accounting/handlers"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
)

type AccountingModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *AccountingModule) RegisterHttpHandlers(router *mux.Router) *AccountingModule {
	prefix := "/accounting"

	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth_svc := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(prefix+"/api/export/quickbooks", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.ExportQuickBooksCSV(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/export/xero", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.ExportXeroCSV(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	return m
}

func (m *AccountingModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Accounting module started")
		return nil
	}
}

func (m *AccountingModule) OnEnd() func() {
	return func() {
		m.Logger.Info("Accounting module stopped")
	}
}
