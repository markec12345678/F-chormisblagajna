package multipayment

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/multipayment/handlers"
)

type MultiPaymentModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *MultiPaymentModule) RegisterHttpHandlers(router *mux.Router, prefix string) {

	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth_svc := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(prefix+"/api/payments", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.AddPayment(m.Config, m.Logger), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders/{orderId}/payments", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetPaymentsByOrder(m.Config, m.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/orders/{orderId}/summary", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetPaymentSummary(m.Config, m.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/daily", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetDailyPayments(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/orders/{orderId}/settle", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.SettleOrder(m.Config, m.Logger), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/payments/{id}/refund", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.RefundPayment(m.Config, m.Logger), "admin"))).Methods("DELETE", "OPTIONS")

}

func (m *MultiPaymentModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("MultiPayment module started")
		return nil
	}
}

func (m *MultiPaymentModule) OnEnd() func() {
	return func() {
		m.Logger.Info("MultiPayment module stopped")
	}
}
