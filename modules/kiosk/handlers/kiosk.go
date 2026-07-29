package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	kk_models "github.com/nutrixpos/pos/modules/kiosk/models"
	"github.com/nutrixpos/pos/modules/kiosk/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
}

func GetKioskConfigs(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.KioskService{Logger: log, Config: cfg}
		d, _ := svc.GetConfigs()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: d})
	}
}
func SaveKioskConfig(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var k kk_models.KioskConfig
		if json.NewDecoder(r.Body).Decode(&k) != nil {
			http.Error(w, "invalid", http.StatusBadRequest)
			return
		}
		svc := services.KioskService{Logger: log, Config: cfg}
		svc.SaveConfig(&k)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: k})
	}
}
func PlaceKioskOrder(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var o kk_models.KioskOrder
		if json.NewDecoder(r.Body).Decode(&o) != nil {
			http.Error(w, "invalid", http.StatusBadRequest)
			return
		}
		svc := services.KioskService{Logger: log, Config: cfg}
		if err := svc.PlaceOrder(&o); err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: o})
	}
}
func GetKioskOrders(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.KioskService{Logger: log, Config: cfg}
		d, _ := svc.GetOrders()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: d})
	}
}
func ServeKioskMenu(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		svc := services.KioskService{Logger: log, Config: cfg}
		configs, _ := svc.GetConfigs()
		for _, k := range configs {
			if k.Id == vars["id"] || k.Name == vars["id"] {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(JSONApiOkResponse{Data: k})
				return
			}
		}
		http.Error(w, "not found", http.StatusNotFound)
	}
}
