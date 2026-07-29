package multilocation

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/multilocation/handlers"

	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
)

type MultiLocationModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *MultiLocationModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Multi-Location module started")
		return nil
	}
}

func (m *MultiLocationModule) OnEnd() func() {
	return func() {}
}

func (m *MultiLocationModule) RegisterHttpHandlers(router *mux.Router, prefix string) {
	var auth_svc auth_mw.IAuthService

	if m.Config.Auth.Enabled {
		jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
		auth_svc = auth_mw.NewInternalAuth(m.Config, jwtUtil)
	} else {
		auth_svc = auth_mw.NewNoAuth(m.Config)
	}

	router.Handle(prefix+"/api/multilocation/dashboard", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetDashboard(m.Config, m.Logger), "admin"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/multilocation/comparison", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetComparison(m.Config, m.Logger), "admin"),
	)).Methods("GET", "OPTIONS")
}
