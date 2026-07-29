package inventorytransfer

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/inventorytransfer/handlers"

	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
)

type InventoryTransferModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *InventoryTransferModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Inventory Transfer module started")
		return nil
	}
}

func (m *InventoryTransferModule) OnEnd() func() {
	return func() {}
}

func (m *InventoryTransferModule) RegisterHttpHandlers(router *mux.Router, prefix string) {
	var auth_svc auth_mw.IAuthService

	if m.Config.Auth.Enabled {
		jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
		auth_svc = auth_mw.NewInternalAuth(m.Config, jwtUtil)
	} else {
		auth_svc = auth_mw.NewNoAuth(m.Config)
	}

	router.Handle(prefix+"/api/transfers", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetAllTransfers(m.Config, m.Logger), "admin"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/transfers", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.CreateTransfer(m.Config, m.Logger), "admin"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/transfers/{id}/status", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.UpdateTransferStatus(m.Config, m.Logger), "admin"),
	)).Methods("PUT", "OPTIONS")
}
