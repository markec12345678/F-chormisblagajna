package waste

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/waste/handlers"
)

type WasteModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *WasteModule) RegisterHttpHandlers(router *mux.Router, prefix string) {

	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth_svc := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(prefix+"/api/waste", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetAllWaste(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/waste", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CreateWaste(m.Config, m.Logger), "admin", "chef"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/waste/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteWaste(m.Config, m.Logger), "admin"))).Methods("DELETE", "OPTIONS")
	router.Handle(prefix+"/api/waste/summary", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetWasteSummary(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
}

func (m *WasteModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Waste module started")
		return nil
	}
}

func (m *WasteModule) OnEnd() func() {
	return func() {
		m.Logger.Info("Waste module stopped")
	}
}
