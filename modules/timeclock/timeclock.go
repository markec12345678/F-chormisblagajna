package timeclock

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/timeclock/handlers"
)

type TimeClockModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *TimeClockModule) RegisterHttpHandlers(router *mux.Router, prefix string) {

	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth_svc := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(prefix+"/api/clock-in", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.ClockIn(m.Config, m.Logger), "admin", "cashier", "chef"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/clock-out/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.ClockOut(m.Config, m.Logger), "admin", "cashier", "chef"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/active", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetActiveEntries(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/dashboard", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetDashboard(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
}

func (m *TimeClockModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Time Clock module started")
		return nil
	}
}

func (m *TimeClockModule) OnEnd() func() {
	return func() {
		m.Logger.Info("Time Clock module stopped")
	}
}
