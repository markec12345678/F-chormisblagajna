package notification

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/notification/handlers"

	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
)

type NotificationModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *NotificationModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Notification module started")
		return nil
	}
}

func (m *NotificationModule) OnEnd() func() {
	return func() {}
}

func (m *NotificationModule) RegisterHttpHandlers(router *mux.Router, prefix string) {
	var auth_svc auth_mw.IAuthService

	if m.Config.Auth.Enabled {
		jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
		auth_svc = auth_mw.NewInternalAuth(m.Config, jwtUtil)
	} else {
		auth_svc = auth_mw.NewNoAuth(m.Config)
	}

	router.Handle(prefix+"/api/notifications", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetNotifications(m.Config, m.Logger), "admin", "cashier", "chef"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/notifications/{id}/read", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.MarkAsRead(m.Config, m.Logger), "admin", "cashier", "chef"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/notifications/read-all", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.MarkAllAsRead(m.Config, m.Logger), "admin", "cashier", "chef"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/notifications/unread-count", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetUnreadCount(m.Config, m.Logger), "admin", "cashier", "chef"),
	)).Methods("GET", "OPTIONS")
}
