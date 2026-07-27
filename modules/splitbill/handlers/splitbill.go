package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/helpers"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/splitbill/dto"
	"github.com/nutrixpos/pos/modules/splitbill/services"
)

type JSONAPIMeta struct {
	TotalRecords int `json:"total_records"`
	PageNumber   int `json:"page_number"`
	PageSize     int `json:"page_size"`
	PageCount    int `json:"page_count"`
}

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta JSONAPIMeta `json:"meta"`
}

func CreateSplitBill(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateSplitBillRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		svc := services.SplitBillService{Logger: logger, Config: config}

		orderTotal := 0.0
		if req.SplitType == "equal" || req.SplitType == "by_item" {
			if req.SplitCount < 2 {
				req.SplitCount = 2
			}
		}

		splitBill, err := svc.CreateSplitBill(req, orderTotal)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		jsonResponse, err := json.Marshal(JSONApiOkResponse{
			Data: splitBill,
			Meta: JSONAPIMeta{TotalRecords: 1},
		})
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Write(jsonResponse)
	}
}

func GetSplitBill(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		svc := services.SplitBillService{Logger: logger, Config: config}

		splitBill, err := svc.GetSplitBill(id)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(JSONApiOkResponse{
			Data: splitBill,
			Meta: JSONAPIMeta{TotalRecords: 1},
		})
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Write(jsonResponse)
	}
}

func GetSplitBillByOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		orderId := params["orderId"]

		svc := services.SplitBillService{Logger: logger, Config: config}

		splitBill, err := svc.GetSplitBillByOrder(orderId)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(JSONApiOkResponse{
			Data: splitBill,
			Meta: JSONAPIMeta{TotalRecords: 1},
		})
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Write(jsonResponse)
	}
}

func PaySplitPart(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var req dto.PaySplitPartRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		svc := services.SplitBillService{Logger: logger, Config: config}

		splitBill, err := svc.PaySplitPart(id, req)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(JSONApiOkResponse{
			Data: splitBill,
			Meta: JSONAPIMeta{TotalRecords: 1},
		})
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Write(jsonResponse)
	}
}

func GetSplitBills(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageNumber, pageSize := helpers.ParsePagination(r, 50)

		svc := services.SplitBillService{Logger: logger, Config: config}

		bills, totalRecords, err := svc.GetSplitBills(services.GetSplitBillsParams{
			PageNumber: pageNumber,
			PageSize:   pageSize,
		})
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		pageCount := totalRecords / pageSize
		if totalRecords%pageSize > 0 {
			pageCount++
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(JSONApiOkResponse{
			Data: bills,
			Meta: JSONAPIMeta{
				TotalRecords: totalRecords,
				PageNumber:   pageNumber,
				PageSize:     pageSize,
				PageCount:    pageCount,
			},
		})
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Write(jsonResponse)
	}
}
