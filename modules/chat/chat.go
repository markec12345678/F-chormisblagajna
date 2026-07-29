package chat

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/chat/handlers"
)

type ChatModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *ChatModule) RegisterHttpHandlers(router *mux.Router) *ChatModule {
	prefix := "/chat"

	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth_svc := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(prefix+"/api/channels", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetChannels(m.Config, m.Logger), "admin", "superuser", "chef", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/channels", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CreateChannel(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/messages", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetMessages(m.Config, m.Logger), "admin", "superuser", "chef", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/messages", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.SendMessage(m.Config, m.Logger), "admin", "superuser", "chef", "cashier"))).Methods("POST", "OPTIONS")

	return m
}

func (m *ChatModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Chat module started")
		return nil
	}
}

func (m *ChatModule) OnEnd() func() {
	return func() {
		m.Logger.Info("Chat module stopped")
	}
}
