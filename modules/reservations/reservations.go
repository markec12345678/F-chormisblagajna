package reservations

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	core_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/reservations/handlers"
)

type ReservationModule struct {
	Logger logger.ILogger
	Config config.Config
}

func (m *ReservationModule) RegisterHttpHandlers(router *mux.Router, prefix string) {

	jwtUtil := auth_mw.NewJWTUtil(m.Config.Auth.JWTSecret, m.Config.Auth.JWTExpireHrs)
	auth := auth_mw.NewInternalAuth(m.Config, jwtUtil)

	router.Handle(prefix+"/api/reservations", core_middlewares.AllowCors(auth_mw.NewNoAuth(m.Config).AllowAnyOfRoles(handlers.CreateReservation(m.Config, m.Logger)))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/reservations", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.GetAllReservations(m.Config, m.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/reservations/{id}/status", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.UpdateReservationStatus(m.Config, m.Logger), "admin"))).Methods("PUT", "OPTIONS")
	router.Handle(prefix+"/api/reservations/{id}/assign", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.AssignTable(m.Config, m.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/reservations/{id}", core_middlewares.AllowCors(auth.AllowAnyOfRoles(handlers.DeleteReservation(m.Config, m.Logger), "admin"))).Methods("DELETE", "OPTIONS")

}

func (m *ReservationModule) OnStart() func() error {
	return func() error {
		m.Logger.Info("Reservations module started")
		return nil
	}
}

func (m *ReservationModule) OnEnd() func() {
	return func() {
		m.Logger.Info("Reservations module stopped")
	}
}
