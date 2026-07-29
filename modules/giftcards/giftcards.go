package giftcards

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/giftcards/handlers"
)

type GiftCardModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *GiftCardModule) RegisterHttpHandlers(router *mux.Router) *GiftCardModule {
	p := "/giftcards"
	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(p+"/api/cards", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.GetAllCards(m.Config, m.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(p+"/api/cards", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.IssueCard(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(p+"/api/redeem", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.RedeemCard(m.Config, m.Logger), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(p+"/api/cards/{id}/deactivate", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.DeactivateCard(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")

	return m
}
func (m *GiftCardModule) OnStart() func() error { return func() error { m.Logger.Info("GiftCards module started"); return nil } }
func (m *GiftCardModule) OnEnd() func() { return func() { m.Logger.Info("GiftCards module stopped") } }
