package delivery

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/delivery/handlers"
)

type DeliveryModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *DeliveryModule) RegisterHttpHandlers(router *mux.Router, prefix string) {
	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(prefix+"/api/zones", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.GetZones(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/zones", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.SaveZone(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/zones/{id}", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.DeleteZone(m.Config, m.Logger), "admin"))).Methods("DELETE", "OPTIONS")
	router.Handle(prefix+"/api/orders", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.GetDeliveryOrders(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/orders", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.CreateDeliveryOrder(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders/{id}/status", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.UpdateDeliveryStatus(m.Config, m.Logger), "admin"))).Methods("PUT", "OPTIONS")

}

func (m *DeliveryModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Delivery module started")
		return nil
	}
}

func (m *DeliveryModule) OnEnd() func() {
	return func() {
		m.Logger.Info("Delivery module stopped")
	}
}
