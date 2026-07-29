package purchase

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/purchase/handlers"
)

type PurchaseModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *PurchaseModule) RegisterHttpHandlers(router *mux.Router) *PurchaseModule {
	p := "/purchase"
	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(p+"/api/orders", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.GetAllPOs(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(p+"/api/orders", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.CreatePO(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(p+"/api/orders/{id}/receive", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.MarkReceived(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(p+"/api/orders/{id}/cancel", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.CancelPO(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")

	return m
}
func (m *PurchaseModule) OnStart() func() error { return func() error { m.Logger.Info("Purchase module started"); return nil } }
func (m *PurchaseModule) OnEnd() func() { return func() { m.Logger.Info("Purchase module stopped") } }
