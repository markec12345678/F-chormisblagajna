package auditlog

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/auditlog/handlers"
)

type AuditLogModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *AuditLogModule) RegisterHttpHandlers(router *mux.Router) *AuditLogModule {
	prefix := "/auditlog"

	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth_svc := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(prefix+"/api/logs", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetAll(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/logs", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.Create(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/summary", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetSummary(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")

	return m
}

func (m *AuditLogModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("AuditLog module started")
		return nil
	}
}

func (m *AuditLogModule) OnEnd() func() {
	return func() {
		m.Logger.Info("AuditLog module stopped")
	}
}
