package promotion

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/promotion/handlers"

	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
)

type PromotionModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *PromotionModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Promotion module started")
		return nil
	}
}

func (m *PromotionModule) OnEnd() func() {
	return func() {}
}

func (m *PromotionModule) RegisterHttpHandlers(router *mux.Router, prefix string) {
	var auth_svc auth_mw.IAuthService

	if m.Config.Auth.Enabled {
		jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
		auth_svc = auth_mw.NewInternalAuth(m.Config, jwtUtil)
	} else {
		auth_svc = auth_mw.NewNoAuth(m.Config)
	}

	router.Handle(prefix+"/api/promotions", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetPromotions(m.Config, m.Logger), "admin", "cashier"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/promotions/validate", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.ValidatePromotion(m.Config, m.Logger), "admin", "cashier"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/promotions", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.CreatePromotion(m.Config, m.Logger), "admin"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/promotions/{id}", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.UpdatePromotion(m.Config, m.Logger), "admin"),
	)).Methods("PATCH", "OPTIONS")

	router.Handle(prefix+"/api/promotions/{id}", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.DeletePromotion(m.Config, m.Logger), "admin"),
	)).Methods("DELETE", "OPTIONS")
}
