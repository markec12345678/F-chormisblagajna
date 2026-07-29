package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	receipt_models "github.com/nutrixpos/pos/modules/receipt/models"
	"github.com/nutrixpos/pos/modules/receipt/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

func GetTemplates(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.ReceiptService{Logger: log, Config: cfg}
		templates, err := svc.GetTemplates()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to load templates", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: templates})
	}
}

func GetTemplate(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		svc := services.ReceiptService{Logger: log, Config: cfg}
		tpl, err := svc.GetTemplate(id)
		if err != nil {
			http.Error(w, "template not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: tpl})
	}
}

func SaveTemplate(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tpl receipt_models.ReceiptTemplate
		if err := json.NewDecoder(r.Body).Decode(&tpl); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if tpl.Name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		svc := services.ReceiptService{Logger: log, Config: cfg}
		if err := svc.SaveTemplate(&tpl); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to save template", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: tpl})
	}
}

func DeleteTemplate(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		svc := services.ReceiptService{Logger: log, Config: cfg}
		if err := svc.DeleteTemplate(id); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to delete template", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func GetPrintSettings(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.ReceiptService{Logger: log, Config: cfg}
		settings, err := svc.GetPrintSettings()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to load print settings", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: settings})
	}
}

func SavePrintSettings(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var settings receipt_models.PrintSettings
		if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		svc := services.ReceiptService{Logger: log, Config: cfg}
		if err := svc.SavePrintSettings(&settings); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to save print settings", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: settings})
	}
}
