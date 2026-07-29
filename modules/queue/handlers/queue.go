package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	q_models "github.com/nutrixpos/pos/modules/queue/models"
	"github.com/nutrixpos/pos/modules/queue/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
}

func AddToQueue(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var e q_models.QueueEntry
		if json.NewDecoder(r.Body).Decode(&e) != nil || e.CustomerName == "" {
			http.Error(w, "customer_name required", http.StatusBadRequest)
			return
		}
		svc := services.QueueService{Logger: log, Config: cfg}
		if err := svc.Add(&e); err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: e})
	}
}
func GetQueue(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.QueueService{Logger: log, Config: cfg}
		d, _ := svc.GetAll()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: d})
	}
}
func UpdateQueueStatus(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct{ Status string `json:"status"` }
		json.NewDecoder(r.Body).Decode(&body)
		svc := services.QueueService{Logger: log, Config: cfg}
		svc.UpdateStatus(mux.Vars(r)["id"], body.Status)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "updated"}})
	}
}
func RemoveFromQueue(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.QueueService{Logger: log, Config: cfg}
		svc.Remove(mux.Vars(r)["id"])
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "removed"}})
	}
}
