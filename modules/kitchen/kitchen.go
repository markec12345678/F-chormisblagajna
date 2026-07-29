package kitchen

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/kitchen/handlers"

	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
)

type KitchenModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *KitchenModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Kitchen module started")
		return nil
	}
}

func (m *KitchenModule) OnEnd() func() {
	return func() {}
}

func (m *KitchenModule) RegisterHttpHandlers(router *mux.Router, prefix string) {
	var auth_svc auth_mw.IAuthService

	if m.Config.Auth.Enabled {
		jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
		auth_svc = auth_mw.NewInternalAuth(m.Config, jwtUtil)
	} else {
		auth_svc = auth_mw.NewNoAuth(m.Config)
	}

	router.Handle(prefix+"/api/kitchen/stations", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetStations(m.Config, m.Logger), "admin", "chef"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/kitchen/stations", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.CreateStation(m.Config, m.Logger), "admin"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/kitchen/orders/{order_id}/items/status", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.UpdateItemStatus(m.Config, m.Logger), "admin", "chef"),
	)).Methods("PUT", "OPTIONS")

	router.Handle(prefix+"/api/kitchen/orders", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetOrdersByStation(m.Config, m.Logger), "admin", "chef"),
	)).Methods("GET", "OPTIONS")
}
