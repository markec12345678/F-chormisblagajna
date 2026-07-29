package receipt

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/receipt/handlers"
)

type ReceiptModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *ReceiptModule) RegisterHttpHandlers(router *mux.Router) *ReceiptModule {
	prefix := "/receipt"

	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth_svc := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(prefix+"/api/templates", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetTemplates(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/templates/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetTemplate(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/templates", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.SaveTemplate(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/templates/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteTemplate(m.Config, m.Logger), "admin"))).Methods("DELETE", "OPTIONS")
	router.Handle(prefix+"/api/print-settings", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetPrintSettings(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/print-settings", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.SavePrintSettings(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")

	return m
}

func (m *ReceiptModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Receipt module started")
		return nil
	}
}

func (m *ReceiptModule) OnEnd() func() {
	return func() {
		m.Logger.Info("Receipt module stopped")
	}
}
