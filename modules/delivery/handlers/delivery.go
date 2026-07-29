package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	dv_models "github.com/nutrixpos/pos/modules/delivery/models"
	"github.com/nutrixpos/pos/modules/delivery/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
}

func GetZones(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.DeliveryService{Logger: log, Config: cfg}
		d, _ := svc.GetZones()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: d})
	}
}
func SaveZone(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var z dv_models.DeliveryZone
		if json.NewDecoder(r.Body).Decode(&z) != nil {
			http.Error(w, "invalid", http.StatusBadRequest)
			return
		}
		svc := services.DeliveryService{Logger: log, Config: cfg}
		if err := svc.SaveZone(&z); err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: z})
	}
}
func DeleteZone(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.DeliveryService{Logger: log, Config: cfg}
		svc.DeleteZone(mux.Vars(r)["id"])
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "deleted"}})
	}
}
func GetDeliveryOrders(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.DeliveryService{Logger: log, Config: cfg}
		d, _ := svc.GetOrders()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: d})
	}
}
func CreateDeliveryOrder(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var o dv_models.DeliveryOrder
		if json.NewDecoder(r.Body).Decode(&o) != nil {
			http.Error(w, "invalid", http.StatusBadRequest)
			return
		}
		svc := services.DeliveryService{Logger: log, Config: cfg}
		if err := svc.CreateOrder(&o); err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: o})
	}
}
func UpdateDeliveryStatus(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct{ Status string `json:"status"` }
		if json.NewDecoder(r.Body).Decode(&body) != nil || body.Status == "" {
			http.Error(w, "status required", http.StatusBadRequest)
			return
		}
		svc := services.DeliveryService{Logger: log, Config: cfg}
		svc.UpdateOrderStatus(mux.Vars(r)["id"], body.Status)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "updated"}})
	}
}
