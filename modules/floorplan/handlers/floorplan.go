package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	fp_models "github.com/nutrixpos/pos/modules/floorplan/models"
	"github.com/nutrixpos/pos/modules/floorplan/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
}

func GetTables(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.FloorplanService{Logger: log, Config: cfg}
		d, _ := svc.GetTables()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: d})
	}
}
func SaveTable(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t fp_models.FloorTable
		json.NewDecoder(r.Body).Decode(&t)
		svc := services.FloorplanService{Logger: log, Config: cfg}
		svc.SaveTable(&t)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: t})
	}
}
func DeleteTable(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.FloorplanService{Logger: log, Config: cfg}
		svc.DeleteTable(mux.Vars(r)["id"])
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "deleted"}})
	}
}
func GetZones(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.FloorplanService{Logger: log, Config: cfg}
		d, _ := svc.GetZones()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: d})
	}
}
func SaveZone(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var z fp_models.FloorZone
		json.NewDecoder(r.Body).Decode(&z)
		svc := services.FloorplanService{Logger: log, Config: cfg}
		svc.SaveZone(&z)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: z})
	}
}
