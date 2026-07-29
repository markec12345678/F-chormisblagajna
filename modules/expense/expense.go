package expense

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/expense/handlers"
)

type ExpenseModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *ExpenseModule) RegisterHttpHandlers(router *mux.Router) *ExpenseModule {
	prefix := "/expense"

	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth_svc := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(prefix+"/api/expenses", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetAllExpenses(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/expenses", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CreateExpense(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/expenses/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.UpdateExpense(m.Config, m.Logger), "admin"))).Methods("PUT", "OPTIONS")
	router.Handle(prefix+"/api/expenses/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteExpense(m.Config, m.Logger), "admin"))).Methods("DELETE", "OPTIONS")
	router.Handle(prefix+"/api/expenses/summary", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetExpenseSummary(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	return m
}

func (m *ExpenseModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Expense module started")
		return nil
	}
}

func (m *ExpenseModule) OnEnd() func() {
	return func() {
		m.Logger.Info("Expense module stopped")
	}
}
