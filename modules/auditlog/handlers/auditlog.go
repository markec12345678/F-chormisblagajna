package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auditlog_models "github.com/nutrixpos/pos/modules/auditlog/models"
	"github.com/nutrixpos/pos/modules/auditlog/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

func GetAll(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		action := r.URL.Query().Get("action")
		resource := r.URL.Query().Get("resource")
		userId := r.URL.Query().Get("user_id")

		svc := services.AuditLogService{Logger: log, Config: cfg}
		entries, err := svc.GetAll(action, resource, userId, 200)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to load audit logs", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: entries})
	}
}

func Create(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var entry auditlog_models.AuditLogEntry
		if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if entry.Action == "" || entry.Resource == "" {
			http.Error(w, "action and resource are required", http.StatusBadRequest)
			return
		}

		entry.IpAddress = r.RemoteAddr

		svc := services.AuditLogService{Logger: log, Config: cfg}
		if err := svc.Create(&entry); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to create audit log entry", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: entry})
	}
}

func GetSummary(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.AuditLogService{Logger: log, Config: cfg}
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
