package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/tips/dto"
	"github.com/nutrixpos/pos/modules/tips/services"
)

type JSONAPIMeta struct {
	TotalRecords int `json:"total_records"`
}

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta JSONAPIMeta `json:"meta"`
}

func RecordTip(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.RecordTipRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		svc := services.TipsService{Logger: logger, Config: config}
		tip, err := svc.RecordTip(req)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		jsonResponse, _ := json.Marshal(JSONApiOkResponse{Data: tip, Meta: JSONAPIMeta{TotalRecords: 1}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func GetTipsByEmployee(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startDate := r.URL.Query().Get("start_date")
		endDate := r.URL.Query().Get("end_date")

		svc := services.TipsService{Logger: logger, Config: config}
		summaries, err := svc.GetTipsByEmployee(startDate, endDate)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		jsonResponse, _ := json.Marshal(JSONApiOkResponse{Data: summaries, Meta: JSONAPIMeta{TotalRecords: len(summaries)}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func PayoutTips(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.PayoutTipsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		svc := services.TipsService{Logger: logger, Config: config}
		payout, err := svc.PayoutTips(req)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		jsonResponse, _ := json.Marshal(JSONApiOkResponse{Data: payout, Meta: JSONAPIMeta{TotalRecords: 1}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func GetPayouts(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.TipsService{Logger: logger, Config: config}
		payouts, err := svc.GetPayouts()
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		jsonResponse, _ := json.Marshal(JSONApiOkResponse{Data: payouts, Meta: JSONAPIMeta{TotalRecords: len(payouts)}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}
