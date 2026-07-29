package kiosk

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/kiosk/handlers"
)

type KioskModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *KioskModule) RegisterHttpHandlers(router *mux.Router) *KioskModule {
	p := "/kiosk"
	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth := auth_mw.NewInternalAuth(m.Config, jwtUtil)
	noAuth := auth_mw.NewNoAuth(m.Config)

	router.Handle(p+"/api/configs", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.GetKioskConfigs(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(p+"/api/configs", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.SaveKioskConfig(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(p+"/api/orders", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.GetKioskOrders(m.Config, m.Logger), "admin", "kitchen"))).Methods("GET", "OPTIONS")
	router.Handle(p+"/api/orders", core_middlewares.AllowCors(noAuth.AllowAnyOfRoles(handlers.PlaceKioskOrder(m.Config, m.Logger)))).Methods("POST", "OPTIONS")
	router.Handle(p+"/api/menu/{id}", core_middlewares.AllowCors(noAuth.AllowAnyOfRoles(handlers.ServeKioskMenu(m.Config, m.Logger)))).Methods("GET", "OPTIONS")

	return m
}
func (m *KioskModule) OnStart() func() error { return func() error { m.Logger.Info("Kiosk module started"); return nil } }
func (m *KioskModule) OnEnd() func() { return func() { m.Logger.Info("Kiosk module stopped") } }
