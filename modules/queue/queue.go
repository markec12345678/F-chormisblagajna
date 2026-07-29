package queue

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/queue/handlers"
)

type QueueModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *QueueModule) RegisterHttpHandlers(router *mux.Router) *QueueModule {
	p := "/queue"
	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(p+"/api/queue", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.GetQueue(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(p+"/api/queue", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.AddToQueue(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(p+"/api/queue/{id}/status", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.UpdateQueueStatus(m.Config, m.Logger), "admin"))).Methods("PUT", "OPTIONS")
	router.Handle(p+"/api/queue/{id}", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.RemoveFromQueue(m.Config, m.Logger), "admin"))).Methods("DELETE", "OPTIONS")

	return m
}
func (m *QueueModule) OnStart() func() error { return func() error { m.Logger.Info("Queue module started"); return nil } }
func (m *QueueModule) OnEnd() func() { return func() { m.Logger.Info("Queue module stopped") } }
