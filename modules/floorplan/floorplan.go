package floorplan

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/floorplan/handlers"
)

type FloorplanModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *FloorplanModule) RegisterHttpHandlers(router *mux.Router) *FloorplanModule {
	p := "/floorplan"
	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(p+"/api/tables", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.GetTables(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(p+"/api/tables", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.SaveTable(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(p+"/api/tables/{id}", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.DeleteTable(m.Config, m.Logger), "admin"))).Methods("DELETE", "OPTIONS")
	router.Handle(p+"/api/zones", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.GetZones(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(p+"/api/zones", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.SaveZone(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")

	return m
}
func (m *FloorplanModule) OnStart() func() error { return func() error { m.Logger.Info("Floorplan module started"); return nil } }
func (m *FloorplanModule) OnEnd() func() { return func() { m.Logger.Info("Floorplan module stopped") } }
