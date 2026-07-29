package giftcard

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/giftcard/handlers"

	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
)

type GiftCardModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *GiftCardModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Gift Card module started")
		return nil
	}
}

func (m *GiftCardModule) OnEnd() func() {
	return func() {}
}

func (m *GiftCardModule) RegisterHttpHandlers(router *mux.Router, prefix string) {
	var auth_svc auth_mw.IAuthService

	if m.Config.Auth.Enabled {
		jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
		auth_svc = auth_mw.NewInternalAuth(m.Config, jwtUtil)
	} else {
		auth_svc = auth_mw.NewNoAuth(m.Config)
	}

	router.Handle(prefix+"/api/giftcards", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetAllGiftCards(m.Config, m.Logger), "admin"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/giftcards", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.CreateGiftCard(m.Config, m.Logger), "admin"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/giftcards/{code}", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetGiftCard(m.Config, m.Logger), "admin"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/giftcards/redeem", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.RedeemGiftCard(m.Config, m.Logger), "admin", "cashier"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/giftcards/{id}/transactions", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetTransactions(m.Config, m.Logger), "admin"),
	)).Methods("GET", "OPTIONS")
}
