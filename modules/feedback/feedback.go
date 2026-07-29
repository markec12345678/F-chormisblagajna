package feedback

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/feedback/handlers"
)

type FeedbackModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *FeedbackModule) RegisterHttpHandlers(router *mux.Router) *FeedbackModule {
	prefix := "/feedback"

	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth_svc := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(prefix+"/api/submit", core_middlewares.AllowCors(auth_mw.NewNoAuth(m.Config).AllowAnyOfRoles(handlers.SubmitFeedback(m.Config, m.Logger)))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/list", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetAllFeedbacks(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/summary", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetFeedbackSummary(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/{id}/respond", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.RespondToFeedback(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteFeedback(m.Config, m.Logger), "admin"))).Methods("DELETE", "OPTIONS")

	return m
}

func (m *FeedbackModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Feedback module started")
		return nil
	}
}

func (m *FeedbackModule) OnEnd() func() {
	return func() {
		m.Logger.Info("Feedback module stopped")
	}
}
