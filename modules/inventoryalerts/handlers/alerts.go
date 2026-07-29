package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	alert_models "github.com/nutrixpos/pos/modules/inventoryalerts/models"
	"github.com/nutrixpos/pos/modules/inventoryalerts/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

func GetRules(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.InventoryAlertsService{Logger: log, Config: cfg}
		rules, err := svc.GetRules()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to load rules", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: rules})
	}
}

func SaveRule(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rule alert_models.InventoryAlertRule
		if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		svc := services.InventoryAlertsService{Logger: log, Config: cfg}
		if err := svc.SaveRule(&rule); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to save rule", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: rule})
	}
}

func DeleteRule(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		svc := services.InventoryAlertsService{Logger: log, Config: cfg}
		if err := svc.DeleteRule(id); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to delete rule", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func GetAlerts(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.InventoryAlertsService{Logger: log, Config: cfg}
		alerts, err := svc.GetAlerts()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to load alerts", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: alerts})
	}
}

func GetUnreadAlerts(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.InventoryAlertsService{Logger: log, Config: cfg}
		alerts, err := svc.GetUnreadAlerts()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to load unread alerts", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: alerts})
	}
}

func MarkAsRead(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		svc := services.InventoryAlertsService{Logger: log, Config: cfg}
		if err := svc.MarkAsRead(id); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to mark as read", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "marked"}})
	}
}

func GetSummary(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.InventoryAlertsService{Logger: log, Config: cfg}
		summary, err := svc.GetSummary()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "summary query failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: summary})
	}
}
