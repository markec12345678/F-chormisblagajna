package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	mk_models "github.com/nutrixpos/pos/modules/marketing/models"
	"github.com/nutrixpos/pos/modules/marketing/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
}

func CreateCampaign(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var c mk_models.Campaign
		if json.NewDecoder(r.Body).Decode(&c) != nil || c.Name == "" {
			http.Error(w, "name required", http.StatusBadRequest)
			return
		}
		svc := services.MarketingService{Logger: log, Config: cfg}
		if err := svc.Create(&c); err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: c})
	}
}
func GetAllCampaigns(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.MarketingService{Logger: log, Config: cfg}
		d, _ := svc.GetAll()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: d})
	}
}
func ToggleCampaign(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.MarketingService{Logger: log, Config: cfg}
		svc.ToggleActive(mux.Vars(r)["id"])
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "toggled"}})
	}
}
func DeleteCampaign(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.MarketingService{Logger: log, Config: cfg}
		svc.Delete(mux.Vars(r)["id"])
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "deleted"}})
	}
}
