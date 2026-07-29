package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/helpers"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/scheduling/dto"
	"github.com/nutrixpos/pos/modules/scheduling/services"
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

func GetShifts(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageNumber, pageSize := helpers.ParsePagination(r, 50)
		branchId := r.URL.Query().Get("branch_id")
		startDate := r.URL.Query().Get("start_date")
		endDate := r.URL.Query().Get("end_date")

		svc := services.SchedulingService{Logger: logger, Config: config}

		shifts, totalRecords, err := svc.GetShifts(services.GetShiftsParams{
			PageNumber: pageNumber,
			PageSize:   pageSize,
			BranchId:   branchId,
			StartDate:  startDate,
			EndDate:    endDate,
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
			Data: shifts,
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

func CreateShift(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateShiftRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		svc := services.SchedulingService{Logger: logger, Config: config}

		shift, err := svc.CreateShift(req)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		jsonResponse, err := json.Marshal(JSONApiOkResponse{
			Data: shift,
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

func UpdateShift(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var req dto.UpdateShiftRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		svc := services.SchedulingService{Logger: logger, Config: config}

		shift, err := svc.UpdateShift(id, req)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(JSONApiOkResponse{
			Data: shift,
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

func DeleteShift(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		svc := services.SchedulingService{Logger: logger, Config: config}

		err := svc.DeleteShift(id)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
