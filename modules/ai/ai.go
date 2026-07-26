package ai

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/ai/handlers"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
)

type AIModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (a *AIModule) OnStart() func() error {
	return func() error {
		a.Logger.Info("AI module started")
		return nil
	}
}

func (a *AIModule) OnEnd() func() {
	return func() {}
}

func (a *AIModule) RegisterHttpHandlers(router *mux.Router, prefix string) {
	var auth_svc auth_mw.IAuthService

	if a.Config.Auth.Enabled {
		jwtUtil := auth_mw.NewJWTUtil(a.Config.Auth.JWTSecret, a.Config.Auth.JWTExpireHrs)
		auth_svc = auth_mw.NewInternalAuth(a.Config, jwtUtil)
	} else {
		auth_svc = auth_mw.NewNoAuth(a.Config)
	}

	router.Handle(prefix+"/api/ai/search", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.AISearch(a.Config, a.Logger), "admin", "cashier", "chef"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/ai/voice", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.VoiceOrder(a.Config, a.Logger), "admin", "cashier"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/ai/suggestions", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.SmartSuggestions(a.Config, a.Logger), "admin", "cashier", "chef"),
	)).Methods("GET", "OPTIONS")
}
