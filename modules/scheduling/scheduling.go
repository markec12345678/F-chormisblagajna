package scheduling

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/scheduling/handlers"

	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
)

type SchedulingModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *SchedulingModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Scheduling module started")
		return nil
	}
}

func (m *SchedulingModule) OnEnd() func() {
	return func() {}
}

func (m *SchedulingModule) RegisterHttpHandlers(router *mux.Router, prefix string) {
	var auth_svc auth_mw.IAuthService

	if m.Config.Auth.Enabled {
		jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
		auth_svc = auth_mw.NewInternalAuth(m.Config, jwtUtil)
	} else {
		auth_svc = auth_mw.NewNoAuth(m.Config)
	}

	router.Handle(prefix+"/api/shifts", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetShifts(m.Config, m.Logger), "admin", "cashier"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/shifts", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.CreateShift(m.Config, m.Logger), "admin"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/shifts/{id}", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.UpdateShift(m.Config, m.Logger), "admin"),
	)).Methods("PATCH", "OPTIONS")

	router.Handle(prefix+"/api/shifts/{id}", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.DeleteShift(m.Config, m.Logger), "admin"),
	)).Methods("DELETE", "OPTIONS")
}
