package fiscal_hr

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules"
	"github.com/nutrixpos/pos/modules/fiscal_hr/handlers"
	"github.com/nutrixpos/pos/modules/fiscal_hr/services"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
)

// FiscalModuleHR is the Croatian Fina eRačun fiscalization module.
type FiscalModuleHR struct {
	Logger logger.ILogger
	Config config.Config
}

func (f *FiscalModuleHR) OnStart() func() error {
	return func() error {
		f.Logger.Info("Croatian fiscal module (HR) started")
		return nil
	}
}

func (f *FiscalModuleHR) OnEnd() func() {
	return func() {}
}

func (f *FiscalModuleHR) RegisterHttpHandlers(router *mux.Router, prefix string) {
	var auth_svc auth_mw.IAuthService

	if f.Config.Auth.Enabled {
		jwtUtil := auth_mw.NewJWTUtil(f.Config.Auth.JWTSecret, f.Config.Auth.JWTExpireHrs)
		auth_svc = auth_mw.NewInternalAuth(f.Config, jwtUtil)
	} else {
		auth_svc = auth_mw.NewNoAuth(f.Config)
	}

	router.Handle(prefix+"/api/fiscal_hr/echo", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.FiscalEchoHandlerHR(f.Config, f.Logger), "admin"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/fiscal_hr/invoice", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.FiscalizeOrderHandlerHR(f.Config, f.Logger), "admin", "cashier"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/fiscal_hr/settings", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetFiscalSettingsHandlerHR(f.Config, f.Logger), "admin"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/fiscal_hr/settings", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.UpdateFiscalSettingsHandlerHR(f.Config, f.Logger), "admin"),
	)).Methods("PATCH", "OPTIONS")
}

func (f *FiscalModuleHR) RegisterBackgroundWorkers() []modules.Worker {
	// Croatian fiscal doesn't need offline retry — simpler architecture
	return []modules.Worker{}
}

// NewCISClientForTest creates a CISClient using NewCISClientFromKey (exported for tests).
func NewCISClientForTest(services *services.CISClient) *services.CISClient {
	return services
}
