package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/inventorytransfer/dto"
	"github.com/nutrixpos/pos/modules/inventorytransfer/services"
)

type JSONAPIMeta struct {
	TotalRecords int `json:"total_records"`
}

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta JSONAPIMeta `json:"meta"`
}

func CreateTransfer(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateTransferRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		svc := services.InventoryTransferService{Logger: logger, Config: config}
		transfer, err := svc.CreateTransfer(req)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		jsonResponse, _ := json.Marshal(JSONApiOkResponse{Data: transfer, Meta: JSONAPIMeta{TotalRecords: 1}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func GetAllTransfers(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.InventoryTransferService{Logger: logger, Config: config}
		transfers, err := svc.GetAllTransfers()
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		jsonResponse, _ := json.Marshal(JSONApiOkResponse{Data: transfers, Meta: JSONAPIMeta{TotalRecords: len(transfers)}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func UpdateTransferStatus(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var req dto.UpdateTransferStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		svc := services.InventoryTransferService{Logger: logger, Config: config}
		transfer, err := svc.UpdateTransferStatus(id, req.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		jsonResponse, _ := json.Marshal(JSONApiOkResponse{Data: transfer, Meta: JSONAPIMeta{TotalRecords: 1}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}
