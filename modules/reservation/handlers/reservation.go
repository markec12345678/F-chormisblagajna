package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/helpers"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/reservation/dto"
	"github.com/nutrixpos/pos/modules/reservation/services"
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

func GetReservations(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageNumber, pageSize := helpers.ParsePagination(r, 50)
		branchId := r.URL.Query().Get("branch_id")
		date := r.URL.Query().Get("date")
		status := r.URL.Query().Get("status")

		svc := services.ReservationService{Logger: logger, Config: config}

		reservations, totalRecords, err := svc.GetReservations(services.GetReservationsParams{
			PageNumber: pageNumber,
			PageSize:   pageSize,
			BranchId:   branchId,
			Date:       date,
			Status:     status,
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
			Data: reservations,
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

func GetReservation(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		svc := services.ReservationService{Logger: logger, Config: config}

		reservation, err := svc.GetReservation(id)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(JSONApiOkResponse{
			Data: reservation,
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

func CreateReservation(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateReservationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		svc := services.ReservationService{Logger: logger, Config: config}

		reservation, err := svc.CreateReservation(req)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		jsonResponse, err := json.Marshal(JSONApiOkResponse{
			Data: reservation,
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

func UpdateReservation(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var req dto.UpdateReservationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		svc := services.ReservationService{Logger: logger, Config: config}

		reservation, err := svc.UpdateReservation(id, req)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(JSONApiOkResponse{
			Data: reservation,
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

func DeleteReservation(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		svc := services.ReservationService{Logger: logger, Config: config}

		err := svc.DeleteReservation(id)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
