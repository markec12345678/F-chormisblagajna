package fiscal

import (
	"time"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules"
	"github.com/nutrixpos/pos/modules/fiscal/handlers"
	"github.com/nutrixpos/pos/modules/fiscal/services"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
)

// FiscalModule is the FURS ZAPOS fiscalization module.
type FiscalModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (f *FiscalModule) OnStart() func() error {
	return func() error {
		f.Logger.Info("Fiscal module started")
		return nil
	}
}

func (f *FiscalModule) OnEnd() func() {
	return func() {}
}

func (f *FiscalModule) RegisterHttpHandlers(router *mux.Router, prefix string) {
	var auth_svc auth_mw.IAuthService

	if f.Config.Auth.Enabled {
		jwtUtil := auth_mw.NewJWTUtil(f.Config.Auth.JWTSecret, f.Config.Auth.JWTExpireHrs)
		auth_svc = auth_mw.NewInternalAuth(f.Config, jwtUtil)
	} else {
		auth_svc = auth_mw.NewNoAuth(f.Config)
	}

	router.Handle(prefix+"/api/fiscal/echo", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.FiscalEchoHandler(f.Config, f.Logger), "admin"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/fiscal/invoice", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.FiscalizeOrderHandler(f.Config, f.Logger), "admin", "cashier"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/fiscal/storno", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.StornoHandler(f.Config, f.Logger), "admin"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/fiscal/settings", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetFiscalSettingsHandler(f.Config, f.Logger), "admin"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/fiscal/settings", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.UpdateFiscalSettingsHandler(f.Config, f.Logger), "admin"),
	)).Methods("PATCH", "OPTIONS")
}

func (f *FiscalModule) RegisterBackgroundWorkers() []modules.Worker {
	return []modules.Worker{
		{
			Interval: 5 * time.Minute,
			Task: func() {
				queue := services.NewOfflineFiscalQueue(f.Logger, f.Config)
				queue.ProcessPending()
			},
		},
	}
}
