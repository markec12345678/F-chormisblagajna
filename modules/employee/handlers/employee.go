package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/employee/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

func GetEmployeePerformance(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startDate := r.URL.Query().Get("start_date")
		endDate := r.URL.Query().Get("end_date")

		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			start = time.Now().AddDate(0, 0, -30)
		}

		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			end = time.Now()
		}

		svc := services.EmployeeService{Logger: log, Config: cfg}
		summary, err := svc.GetPerformance(start, end)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "performance query failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(JSONApiOkResponse{Data: summary})
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Write(jsonResponse)
	}
}
