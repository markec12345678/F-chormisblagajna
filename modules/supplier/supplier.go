package supplier

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/supplier/handlers"
)

type SupplierModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *SupplierModule) RegisterHttpHandlers(router *mux.Router) *SupplierModule {
	prefix := "/supplier"

	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth_svc := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(prefix+"/api/suppliers", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetAllSuppliers(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/suppliers", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CreateSupplier(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/suppliers/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetSupplier(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/suppliers/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.UpdateSupplier(m.Config, m.Logger), "admin"))).Methods("PUT", "OPTIONS")
	router.Handle(prefix+"/api/suppliers/{id}", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteSupplier(m.Config, m.Logger), "admin"))).Methods("DELETE", "OPTIONS")
	router.Handle(prefix+"/api/suppliers/{supplier_id}/orders", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetSupplierOrders(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/supplier-orders", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CreateSupplierOrder(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/supplier-orders/{id}/status", core_middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.UpdateSupplierOrderStatus(m.Config, m.Logger), "admin"))).Methods("PUT", "OPTIONS")
	return m
}

func (m *SupplierModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Supplier module started")
		return nil
	}
}

func (m *SupplierModule) OnEnd() func() {
	return func() {
		m.Logger.Info("Supplier module stopped")
	}
}
