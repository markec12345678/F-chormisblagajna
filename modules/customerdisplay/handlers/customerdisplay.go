package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	cd_models "github.com/nutrixpos/pos/modules/customerdisplay/models"
	"github.com/nutrixpos/pos/modules/customerdisplay/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

func GetAllConfigs(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.CustomerDisplayService{Logger: log, Config: cfg}
		configs, err := svc.GetAllConfigs()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to load display configs", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: configs})
	}
}

func GetConfig(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		svc := services.CustomerDisplayService{Logger: log, Config: cfg}
		config, err := svc.GetConfig(id)
		if err != nil {
			http.Error(w, "config not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: config})
	}
}

func SaveConfig(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var config cd_models.DisplayConfig
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		svc := services.CustomerDisplayService{Logger: log, Config: cfg}
		if err := svc.SaveConfig(&config); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to save config", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: config})
	}
}

func DeleteConfig(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		svc := services.CustomerDisplayService{Logger: log, Config: cfg}
		if err := svc.DeleteConfig(id); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to delete config", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "deleted"}})
	}
}

func GetDisplayContent(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		displayId := vars["id"]

		svc := services.CustomerDisplayService{Logger: log, Config: cfg}
		content, err := svc.GetDisplayContent(displayId)
		if err != nil {
			http.Error(w, "display not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: content})
	}
}
