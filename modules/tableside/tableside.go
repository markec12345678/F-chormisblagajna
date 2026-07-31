package tableside

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/tableside/handlers"
)

type TablesideModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *TablesideModule) RegisterHttpHandlers(router *mux.Router, prefix string) {

	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth_svc := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(prefix+"/api/sessions", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetAllSessions(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/sessions", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CreateSession(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/sessions/{id}/close", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CloseSession(m.Config, m.Logger), "admin", "waiter"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders/{id}/status", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.UpdateOrderStatus(m.Config, m.Logger), "admin", "kitchen"))).Methods("PUT", "OPTIONS")
	router.Handle(prefix+"/api/sessions/{sessionId}/orders", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetOrdersBySession(m.Config, m.Logger), "admin", "waiter"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/qr/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetQrUrl(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")

	// public endpoints (no auth)
	noAuth := auth_mw.NewNoAuth(m.Config)
	router.Handle(prefix+"/api/menu/place-order", core_middlewares.AllowCors(noAuth.AllowAnyOfRoles(handlers.PlaceOrder(m.Config, m.Logger)))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/menu/{sessionId}/orders", core_middlewares.AllowCors(noAuth.AllowAnyOfRoles(handlers.GetOrdersBySession(m.Config, m.Logger)))).Methods("GET", "OPTIONS")

}

func (m *TablesideModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Tableside module started")
		return nil
	}
}

func (m *TablesideModule) OnEnd() func() {
	return func() {
		m.Logger.Info("Tableside module stopped")
	}
}
