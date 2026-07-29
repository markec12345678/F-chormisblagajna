package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/kitchen/services"
)

type JSONAPIMeta struct {
	TotalRecords int `json:"total_records"`
}

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta JSONAPIMeta `json:"meta"`
}

func GetStations(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.KitchenService{Logger: logger, Config: config}
		stations, err := svc.GetStations()
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		jsonResponse, _ := json.Marshal(JSONApiOkResponse{Data: stations, Meta: JSONAPIMeta{TotalRecords: len(stations)}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func CreateStation(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name     string `json:"name"`
			BranchID string `json:"branch_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		svc := services.KitchenService{Logger: logger, Config: config}
		station, err := svc.CreateStation(req.Name, req.BranchID)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		jsonResponse, _ := json.Marshal(JSONApiOkResponse{Data: station, Meta: JSONAPIMeta{TotalRecords: 1}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func UpdateItemStatus(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderID := vars["order_id"]

		var req struct {
			ItemIndex int    `json:"item_index"`
			Status    string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		svc := services.KitchenService{Logger: logger, Config: config}
		err := svc.UpdateItemStatus(orderID, req.ItemIndex, req.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{
			Data: map[string]string{"message": "item status updated"},
			Meta: JSONAPIMeta{TotalRecords: 1},
		})
	}
}

func GetOrdersByStation(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		station := r.URL.Query().Get("station")

		svc := services.KitchenService{Logger: logger, Config: config}
		orders, err := svc.GetOrdersByStation(station)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		jsonResponse, _ := json.Marshal(JSONApiOkResponse{Data: orders, Meta: JSONAPIMeta{TotalRecords: len(orders)}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}
