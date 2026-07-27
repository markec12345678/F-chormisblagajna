package branch

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/branch/handlers"

	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
)

type BranchModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (b *BranchModule) OnStart() func() error {
	return func() error {
		b.Logger.Info("Branch module started")
		return nil
	}
}

func (b *BranchModule) OnEnd() func() {
	return func() {}
}

func (b *BranchModule) RegisterHttpHandlers(router *mux.Router, prefix string) {
	var auth_svc auth_mw.IAuthService

	if b.Config.Auth.Enabled {
		jwtUtil := auth_mw.NewJWTUtil(b.Config.Auth.JWTSecret, b.Config.Auth.JWTExpireHrs)
		auth_svc = auth_mw.NewInternalAuth(b.Config, jwtUtil)
	} else {
		auth_svc = auth_mw.NewNoAuth(b.Config)
	}

	router.Handle(prefix+"/api/branches", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetBranches(b.Config, b.Logger), "admin", "cashier"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/branches", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.CreateBranch(b.Config, b.Logger), "admin"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/branches/{id}", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetBranch(b.Config, b.Logger), "admin", "cashier"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/branches/{id}", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.UpdateBranch(b.Config, b.Logger), "admin"),
	)).Methods("PATCH", "OPTIONS")

	router.Handle(prefix+"/api/branches/{id}", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.DeleteBranch(b.Config, b.Logger), "admin"),
	)).Methods("DELETE", "OPTIONS")

	router.Handle(prefix+"/api/branches/{id}/stats", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetBranchStats(b.Config, b.Logger), "admin"),
	)).Methods("GET", "OPTIONS")
}
