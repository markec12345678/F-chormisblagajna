package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	pc_models "github.com/nutrixpos/pos/modules/purchase/models"
	"github.com/nutrixpos/pos/modules/purchase/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
}

func CreatePO(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var po pc_models.PurchaseOrder
		if json.NewDecoder(r.Body).Decode(&po) != nil {
			http.Error(w, "invalid", http.StatusBadRequest)
			return
		}
		svc := services.PurchaseService{Logger: log, Config: cfg}
		if err := svc.Create(&po); err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: po})
	}
}
func GetAllPOs(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.PurchaseService{Logger: log, Config: cfg}
		d, _ := svc.GetAll()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: d})
	}
}
func MarkReceived(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.PurchaseService{Logger: log, Config: cfg}
		svc.MarkReceived(mux.Vars(r)["id"])
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "received"}})
	}
}
func CancelPO(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.PurchaseService{Logger: log, Config: cfg}
		svc.Cancel(mux.Vars(r)["id"])
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "cancelled"}})
	}
}
