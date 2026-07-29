package marketing

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/marketing/handlers"
)

type MarketingModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *MarketingModule) RegisterHttpHandlers(router *mux.Router) *MarketingModule {
	p := "/marketing"
	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(p+"/api/campaigns", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.GetAllCampaigns(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(p+"/api/campaigns", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.CreateCampaign(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(p+"/api/campaigns/{id}/toggle", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.ToggleCampaign(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(p+"/api/campaigns/{id}", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.DeleteCampaign(m.Config, m.Logger), "admin"))).Methods("DELETE", "OPTIONS")

	return m
}
func (m *MarketingModule) OnStart() func() error { return func() error { m.Logger.Info("Marketing module started"); return nil } }
func (m *MarketingModule) OnEnd() func() { return func() { m.Logger.Info("Marketing module stopped") } }
