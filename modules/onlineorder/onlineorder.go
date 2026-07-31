package onlineorder

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/onlineorder/handlers"
)

type OnlineOrderModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *OnlineOrderModule) RegisterHttpHandlers(router *mux.Router, prefix string) {

	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth_svc := auth_mw.NewNoAuth(m.Config)
	admin_auth := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	// Public endpoints (no auth required)
	router.Handle(prefix+"/api/menu", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetMenu(m.Config, m.Logger)))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/order", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CreateOrder(m.Config, m.Logger)))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/track/{displayId}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.TrackOrder(m.Config, m.Logger)))).Methods("GET", "OPTIONS")

	// Admin endpoints (auth required)
	router.Handle(prefix+"/api/orders", core_middlewares.AllowCors(admin_auth.AllowAnyOfRoles(handlers.GetAllOrders(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/orders/{id}/status", core_middlewares.AllowCors(admin_auth.AllowAnyOfRoles(handlers.UpdateOrderStatus(m.Config, m.Logger), "admin"))).Methods("PUT", "OPTIONS")

}

func (m *OnlineOrderModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("OnlineOrder module started")
		return nil
	}
}

func (m *OnlineOrderModule) OnEnd() func() {
	return func() {
		m.Logger.Info("OnlineOrder module stopped")
	}
}
