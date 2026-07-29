package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/multilocation/services"
)

type JSONAPIMeta struct {
	TotalRecords int `json:"total_records"`
}

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta JSONAPIMeta `json:"meta"`
}

func GetDashboard(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.MultiLocationService{Logger: logger, Config: config}
		dashboard, err := svc.GetDashboard()
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		jsonResponse, _ := json.Marshal(JSONApiOkResponse{
			Data: dashboard,
			Meta: JSONAPIMeta{TotalRecords: 1},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func GetComparison(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.MultiLocationService{Logger: logger, Config: config}
		comparison, err := svc.GetComparison()
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		jsonResponse, _ := json.Marshal(JSONApiOkResponse{
			Data: comparison,
			Meta: JSONAPIMeta{TotalRecords: len(comparison)},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}
