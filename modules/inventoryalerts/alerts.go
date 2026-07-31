package inventoryalerts

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/inventoryalerts/handlers"
)

type InventoryAlertsModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *InventoryAlertsModule) RegisterHttpHandlers(router *mux.Router, prefix string) {

	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth_svc := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(prefix+"/api/rules", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetRules(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/rules", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.SaveRule(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/rules/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteRule(m.Config, m.Logger), "admin"))).Methods("DELETE", "OPTIONS")
	router.Handle(prefix+"/api/alerts", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetAlerts(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/alerts/unread", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetUnreadAlerts(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/alerts/{id}/read", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.MarkAsRead(m.Config, m.Logger), "admin"))).Methods("PUT", "OPTIONS")
	router.Handle(prefix+"/api/summary", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetSummary(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")

}

func (m *InventoryAlertsModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("InventoryAlerts module started")
		return nil
	}
}

func (m *InventoryAlertsModule) OnEnd() func() {
	return func() {
		m.Logger.Info("InventoryAlerts module stopped")
	}
}
