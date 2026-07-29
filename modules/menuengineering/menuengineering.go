package menuengineering

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/menuengineering/handlers"
)

type MenuEngineeringModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *MenuEngineeringModule) RegisterHttpHandlers(router *mux.Router) *MenuEngineeringModule {
	prefix := "/menuengineering"

	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth_svc := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(prefix+"/api/analysis", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetMenuAnalysis(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	return m
}

func (m *MenuEngineeringModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Menu Engineering module started")
		return nil
	}
}

func (m *MenuEngineeringModule) OnEnd() func() {
	return func() {
		m.Logger.Info("Menu Engineering module stopped")
	}
}
