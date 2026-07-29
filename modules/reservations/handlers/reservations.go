package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	rs_models "github.com/nutrixpos/pos/modules/reservations/models"
	"github.com/nutrixpos/pos/modules/reservations/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

func CreateReservation(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res rs_models.Reservation
		if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if res.CustomerName == "" || res.ReservationDate == "" || res.ReservationTime == "" {
			http.Error(w, "name, date, time required", http.StatusBadRequest)
			return
		}
		svc := services.ReservationService{Logger: log, Config: cfg}
		if err := svc.Create(&res); err != nil {
			log.Error(err.Error())
			http.Error(w, "creation failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: res})
	}
}

func GetAllReservations(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.ReservationService{Logger: log, Config: cfg}
		data, err := svc.GetAll()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "load failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: data})
	}
}

func UpdateReservationStatus(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var body struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Status == "" {
			http.Error(w, "status required", http.StatusBadRequest)
			return
		}
		svc := services.ReservationService{Logger: log, Config: cfg}
		if err := svc.UpdateStatus(vars["id"], body.Status); err != nil {
			log.Error(err.Error())
			http.Error(w, "update failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "updated"}})
	}
}

func AssignTable(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var body struct {
			Table string `json:"table"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Table == "" {
			http.Error(w, "table required", http.StatusBadRequest)
			return
		}
		svc := services.ReservationService{Logger: log, Config: cfg}
		if err := svc.AssignTable(vars["id"], body.Table); err != nil {
			log.Error(err.Error())
			http.Error(w, "assign failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "assigned"}})
	}
}

func DeleteReservation(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		svc := services.ReservationService{Logger: log, Config: cfg}
		if err := svc.Delete(vars["id"]); err != nil {
			log.Error(err.Error())
			http.Error(w, "delete failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "deleted"}})
	}
}
