package employee

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/employee/handlers"
)

type EmployeeModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *EmployeeModule) RegisterHttpHandlers(router *mux.Router, prefix string) {

	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth_svc := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(prefix+"/api/performance", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetEmployeePerformance(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
}

func (m *EmployeeModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Employee module started")
		return nil
	}
}

func (m *EmployeeModule) OnEnd() func() {
	return func() {
		m.Logger.Info("Employee module stopped")
	}
}
