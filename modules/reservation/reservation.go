package reservation

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/reservation/handlers"

	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
)

type ReservationModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *ReservationModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Reservation module started")
		return nil
	}
}

func (m *ReservationModule) OnEnd() func() {
	return func() {}
}

func (m *ReservationModule) RegisterHttpHandlers(router *mux.Router, prefix string) {
	var auth_svc auth_mw.IAuthService

	if m.Config.Auth.Enabled {
		jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
		auth_svc = auth_mw.NewInternalAuth(m.Config, jwtUtil)
	} else {
		auth_svc = auth_mw.NewNoAuth(m.Config)
	}

	router.Handle(prefix+"/api/reservations", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetReservations(m.Config, m.Logger), "admin", "cashier"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/reservations", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.CreateReservation(m.Config, m.Logger), "admin", "cashier"),
	)).Methods("POST", "OPTIONS")

	router.Handle(prefix+"/api/reservations/{id}", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.GetReservation(m.Config, m.Logger), "admin", "cashier"),
	)).Methods("GET", "OPTIONS")

	router.Handle(prefix+"/api/reservations/{id}", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.UpdateReservation(m.Config, m.Logger), "admin", "cashier"),
	)).Methods("PATCH", "OPTIONS")

	router.Handle(prefix+"/api/reservations/{id}", core_middlewares.AllowCors(
		auth_svc.AllowAnyOfRoles(handlers.DeleteReservation(m.Config, m.Logger), "admin"),
	)).Methods("DELETE", "OPTIONS")
}
