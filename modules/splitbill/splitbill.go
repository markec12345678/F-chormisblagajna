package splitbill

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/splitbill/handlers"

	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
)

type SplitBillModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *SplitBillModule) OnStart() func() error {
	return func() error {
		s.Logger.Info("SplitBill module started")
		return nil
	}
}

func (s *SplitBillModule) OnEnd() func() {
	return func() {}
}

func (s *SplitBillModule) RegisterHttpHandlers(router *mux.Router, prefix string) {
	var auth_svc auth_mw.IAuthService

	if s.Config.Auth.Enabled {
		jwtUtil := auth_mw.NewJWTUtil(s.Config.Auth.JWTSecret, s.Config.Auth.JWTExpireHrs)
		auth_svc = auth_mw.NewInternalAuth(s.Config, jwtUtil)
	} else {
		auth_svc = auth_mw.NewNoAuth(s.Config)
	}

	router.Handle(prefix+"/api/split-bills", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.CreateSplitBill(s.Config, s.Logger), "admin", "cashier"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/split-bills", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetSplitBills(s.Config, s.Logger), "admin", "cashier"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/split-bills/{id}", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetSplitBill(s.Config, s.Logger), "admin", "cashier"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/split-bills/order/{orderId}", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetSplitBillByOrder(s.Config, s.Logger), "admin", "cashier"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/split-bills/{id}/pay", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.PaySplitPart(s.Config, s.Logger), "admin", "cashier"),
	)).Methods("POST", "OPTIONS")
}
