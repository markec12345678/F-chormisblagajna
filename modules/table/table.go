package table

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/table/handlers"

	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
)

type TableModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (t *TableModule) OnStart() func() error {
	return func() error {
		t.Logger.Info("Table module started")
		return nil
	}
}

func (t *TableModule) OnEnd() func() {
	return func() {}
}

func (t *TableModule) RegisterHttpHandlers(router *mux.Router, prefix string) {
	var auth_svc auth_mw.IAuthService

	if t.Config.Auth.Enabled {
		jwtUtil := auth_mw.NewJWTUtil(t.Config.Auth.JWTSecret, t.Config.Auth.JWTExpireHrs)
		auth_svc = auth_mw.NewInternalAuth(t.Config, jwtUtil)
	} else {
		auth_svc = auth_mw.NewNoAuth(t.Config)
	}

	router.Handle(prefix+"/api/tables", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetTables(t.Config, t.Logger), "admin", "cashier", "chef"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/tables", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.CreateTable(t.Config, t.Logger), "admin"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/tables/{id}", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetTable(t.Config, t.Logger), "admin", "cashier", "chef"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/tables/{id}", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.UpdateTable(t.Config, t.Logger), "admin"),
	)).Methods("PATCH", "OPTIONS")

	router.Handle(prefix+"/api/tables/{id}", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.DeleteTable(t.Config, t.Logger), "admin"),
	)).Methods("DELETE", "OPTIONS")

	router.Handle(prefix+"/api/tables/{id}/status", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.UpdateTableStatus(t.Config, t.Logger), "admin", "cashier"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/tables/{id}/qr", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetQRCode(t.Config, t.Logger), "admin", "cashier"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/tables/transfer", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.TransferTable(t.Config, t.Logger), "admin", "cashier"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/tables/merge", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.MergeTables(t.Config, t.Logger), "admin"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/tables/floorplan", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetFloorPlan(t.Config, t.Logger), "admin", "cashier", "chef"),
	)).Methods("GET", "OPTIONS")
}
