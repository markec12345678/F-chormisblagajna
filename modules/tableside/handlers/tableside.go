package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	ts_models "github.com/nutrixpos/pos/modules/tableside/models"
	"github.com/nutrixpos/pos/modules/tableside/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

func GetAllSessions(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.TablesideService{Logger: log, Config: cfg}
		sessions, err := svc.GetAllSessions()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to load sessions", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: sessions})
	}
}

func CreateSession(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var session ts_models.TableSession
		if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if session.TableLabel == "" {
			http.Error(w, "table_label is required", http.StatusBadRequest)
			return
		}

		svc := services.TablesideService{Logger: log, Config: cfg}
		if err := svc.CreateSession(&session); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to create session", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: session})
	}
}

func CloseSession(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		svc := services.TablesideService{Logger: log, Config: cfg}
		if err := svc.CloseSession(id); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to close session", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "closed"}})
	}
}

func PlaceOrder(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var order ts_models.TableOrder
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if order.SessionId == "" || len(order.Items) == 0 {
			http.Error(w, "session_id and items are required", http.StatusBadRequest)
			return
		}

		svc := services.TablesideService{Logger: log, Config: cfg}
		if err := svc.PlaceTableOrder(&order); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to place order", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: order})
	}
}

func GetOrdersBySession(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sessionId := vars["sessionId"]

		svc := services.TablesideService{Logger: log, Config: cfg}
		orders, err := svc.GetOrdersBySession(sessionId)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to load orders", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: orders})
	}
}

func UpdateOrderStatus(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderId := vars["id"]

		var body struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Status == "" {
			http.Error(w, "status is required", http.StatusBadRequest)
			return
		}

		svc := services.TablesideService{Logger: log, Config: cfg}
		if err := svc.UpdateOrderStatus(orderId, body.Status); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to update status", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "updated"}})
	}
}

func GetQrUrl(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sessionId := vars["id"]
		host := r.URL.Query().Get("host")
		if host == "" {
			host = r.Host
		}

		svc := services.TablesideService{Logger: log, Config: cfg}
		qr, err := svc.GetQrUrl(sessionId, host)
		if err != nil {
			http.Error(w, "session not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: qr})
	}
}
