package loyalty

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/loyalty/handlers"
)

type LoyaltyModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *LoyaltyModule) RegisterHttpHandlers(router *mux.Router, prefix string) {

	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth_svc := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(prefix+"/api/cards", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetAllCards(m.Config, m.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/cards", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CreateCard(m.Config, m.Logger), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/cards/{cardId}/points", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.AddPoints(m.Config, m.Logger), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/cards/{cardId}/redeem", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.RedeemPoints(m.Config, m.Logger), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/rewards", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetAllRewards(m.Config, m.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/rewards", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CreateReward(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/settings", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetSettings(m.Config, m.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")

}

func (m *LoyaltyModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Loyalty module started")
		return nil
	}
}

func (m *LoyaltyModule) OnEnd() func() {
	return func() {
		m.Logger.Info("Loyalty module stopped")
	}
}
