package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	waste_models "github.com/nutrixpos/pos/modules/waste/models"
	"github.com/nutrixpos/pos/modules/waste/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

func GetAllWaste(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.WasteService{Logger: log, Config: cfg}
		entries, err := svc.GetAllWaste()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: entries})
	}
}

func CreateWaste(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var entry waste_models.WasteEntry
		if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		svc := services.WasteService{Logger: log, Config: cfg}
		if err := svc.CreateWaste(&entry); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to create waste entry", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: entry})
	}
}

func DeleteWaste(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		svc := services.WasteService{Logger: log, Config: cfg}
		if err := svc.DeleteWaste(id); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to delete waste entry", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func GetWasteSummary(cfg config.Config, log logger.ILogger) http.HandlerFunc {
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

		svc := services.WasteService{Logger: log, Config: cfg}
		summary, err := svc.GetSummary(start, end)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "summary query failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: summary})
	}
}
