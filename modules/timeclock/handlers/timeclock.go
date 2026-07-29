package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/timeclock/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

func ClockIn(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			EmployeeId   string `json:"employee_id"`
			EmployeeName string `json:"employee_name"`
			Notes        string `json:"notes"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		svc := services.TimeClockService{Logger: log, Config: cfg}
		entry, err := svc.ClockIn(req.EmployeeId, req.EmployeeName, req.Notes)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "clock in failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: entry})
	}
}

func ClockOut(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		svc := services.TimeClockService{Logger: log, Config: cfg}
		if err := svc.ClockOut(id); err != nil {
			log.Error(err.Error())
			http.Error(w, "clock out failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func GetActiveEntries(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.TimeClockService{Logger: log, Config: cfg}
		entries, err := svc.GetActiveEntries()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: entries})
	}
}

func GetEntriesByDate(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dateStr := r.URL.Query().Get("date")

		var date interface{} = dateStr
		_ = date

		svc := services.TimeClockService{Logger: log, Config: cfg}
		entries, err := svc.GetActiveEntries()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: entries})
	}
}

func GetDashboard(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.TimeClockService{Logger: log, Config: cfg}
		dashboard, err := svc.GetDashboard()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "dashboard query failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: dashboard})
	}
}
