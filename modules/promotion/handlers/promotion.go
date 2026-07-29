package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/helpers"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/promotion/dto"
	"github.com/nutrixpos/pos/modules/promotion/services"
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

func GetPromotions(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageNumber, pageSize := helpers.ParsePagination(r, 50)
		activeOnly := r.URL.Query().Get("active_only") == "true"

		svc := services.PromotionService{Logger: logger, Config: config}

		promotions, totalRecords, err := svc.GetPromotions(services.GetPromotionsParams{
			PageNumber: pageNumber,
			PageSize:   pageSize,
			ActiveOnly: activeOnly,
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
			Data: promotions,
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

func ValidatePromotion(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		orderTotal := 0.0
		if ot := r.URL.Query().Get("order_total"); ot != "" {
			json.Unmarshal([]byte(ot), &orderTotal)
		}

		svc := services.PromotionService{Logger: logger, Config: config}

		promotion, err := svc.ValidatePromotionCode(code, orderTotal)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(JSONApiOkResponse{
			Data: promotion,
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

func CreatePromotion(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreatePromotionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		svc := services.PromotionService{Logger: logger, Config: config}

		promotion, err := svc.CreatePromotion(req)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		jsonResponse, err := json.Marshal(JSONApiOkResponse{
			Data: promotion,
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

func UpdatePromotion(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var req dto.UpdatePromotionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		svc := services.PromotionService{Logger: logger, Config: config}

		promotion, err := svc.UpdatePromotion(id, req)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(JSONApiOkResponse{
			Data: promotion,
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

func DeletePromotion(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		svc := services.PromotionService{Logger: logger, Config: config}

		err := svc.DeletePromotion(id)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
