package training

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/training/handlers"
)

type TrainingModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *TrainingModule) RegisterHttpHandlers(router *mux.Router, prefix string) {

	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth_svc := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(prefix+"/api/modules", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetModules(m.Config, m.Logger), "admin", "staff"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/modules/{module}/steps", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetSteps(m.Config, m.Logger), "admin", "staff"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/sessions", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.StartSession(m.Config, m.Logger), "admin", "staff"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/sessions/{id}/step", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CompleteStep(m.Config, m.Logger), "admin", "staff"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/sessions/{id}/complete", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CompleteSession(m.Config, m.Logger), "admin", "staff"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/users/{userId}/progress", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetUserProgress(m.Config, m.Logger), "admin", "staff"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/users/{userId}/sessions", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetUserSessions(m.Config, m.Logger), "admin", "staff"))).Methods("GET", "OPTIONS")

}

func (m *TrainingModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Training module started")
		return nil
	}
}

func (m *TrainingModule) OnEnd() func() {
	return func() {
		m.Logger.Info("Training module stopped")
	}
}
