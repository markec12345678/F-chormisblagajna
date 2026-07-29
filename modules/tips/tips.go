package tips

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/tips/handlers"

	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
)

type TipsModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *TipsModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Tips module started")
		return nil
	}
}

func (m *TipsModule) OnEnd() func() {
	return func() {}
}

func (m *TipsModule) RegisterHttpHandlers(router *mux.Router, prefix string) {
	var auth_svc auth_mw.IAuthService

	if m.Config.Auth.Enabled {
		jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
		auth_svc = auth_mw.NewInternalAuth(m.Config, jwtUtil)
	} else {
		auth_svc = auth_mw.NewNoAuth(m.Config)
	}

	router.Handle(prefix+"/api/tips", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.RecordTip(m.Config, m.Logger), "admin", "cashier"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/tips/summary", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetTipsByEmployee(m.Config, m.Logger), "admin"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/tips/payout", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.PayoutTips(m.Config, m.Logger), "admin"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/tips/payouts", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetPayouts(m.Config, m.Logger), "admin"),
	)).Methods("GET", "OPTIONS")
}
