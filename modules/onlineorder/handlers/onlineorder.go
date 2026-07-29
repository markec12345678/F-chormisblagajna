package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	onlineorder_models "github.com/nutrixpos/pos/modules/onlineorder/models"
	"github.com/nutrixpos/pos/modules/onlineorder/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

func GetMenu(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.OnlineOrderService{Logger: log, Config: cfg}
		menu, err := svc.GetMenu()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to load menu", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: menu})
	}
}

func CreateOrder(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var order onlineorder_models.OnlineOrder
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if order.CustomerName == "" || len(order.Items) == 0 {
			http.Error(w, "customer_name and items are required", http.StatusBadRequest)
			return
		}

		svc := services.OnlineOrderService{Logger: log, Config: cfg}
		if err := svc.CreateOrder(&order); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to create order", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: order})
	}
}

func TrackOrder(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		displayId := vars["displayId"]

		svc := services.OnlineOrderService{Logger: log, Config: cfg}
		order, err := svc.TrackOrder(displayId)
		if err != nil {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: order})
	}
}

func GetAllOrders(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")

		svc := services.OnlineOrderService{Logger: log, Config: cfg}
		orders, err := svc.GetAllOrders(status)
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
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		svc := services.OnlineOrderService{Logger: log, Config: cfg}
		if err := svc.UpdateOrderStatus(orderId, body.Status); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to update order status", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "updated"}})
	}
}
