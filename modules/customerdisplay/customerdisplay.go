package customerdisplay

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/customerdisplay/handlers"
)

type CustomerDisplayModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *CustomerDisplayModule) RegisterHttpHandlers(router *mux.Router, prefix string) {

	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth_svc := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(prefix+"/api/configs", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetAllConfigs(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/configs/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetConfig(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/configs", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.SaveConfig(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/configs/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteConfig(m.Config, m.Logger), "admin"))).Methods("DELETE", "OPTIONS")
	router.Handle(prefix+"/api/display/{id}", core_middlewares.AllowCors(handlers.GetDisplayContent(m.Config, m.Logger))).Methods("GET", "OPTIONS")

}

func (m *CustomerDisplayModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("CustomerDisplay module started")
		return nil
	}
}

func (m *CustomerDisplayModule) OnEnd() func() {
	return func() {
		m.Logger.Info("CustomerDisplay module stopped")
	}
}
