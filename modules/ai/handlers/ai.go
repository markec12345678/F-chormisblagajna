package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/ai/models"
	"github.com/nutrixpos/pos/modules/ai/services"
)

type JSONAPIMeta struct {
	TotalRecords int `json:"total_records"`
	PageNumber   int `json:"page_number"`
	PageSize     int `json:"page_size"`
	PageCount    int `json:"page_count"`
}

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta JSONAPIMeta  `json:"meta"`
}

func AISearch(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := models.AISearchRequest{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if request.Query == "" {
			http.Error(w, "query is required", http.StatusBadRequest)
			return
		}

		if request.Limit <= 0 {
			request.Limit = 20
		}

		aiSvc := services.AIService{
			Logger: logger,
			Config: config,
		}

		response, err := aiSvc.AISearch(request.Query, request.BranchId, request.Language, request.Limit)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(JSONApiOkResponse{
			Data: response,
			Meta: JSONAPIMeta{
				TotalRecords: len(response.Results),
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

func VoiceOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := models.VoiceOrderRequest{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if request.AudioBase64 == "" {
			http.Error(w, "audio_base64 is required", http.StatusBadRequest)
			return
		}

		if request.Language == "" {
			request.Language = "sl"
		}

		aiSvc := services.AIService{
			Logger: logger,
			Config: config,
		}

		response, err := aiSvc.ProcessVoiceOrder(request.AudioBase64, request.Language, request.BranchId)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(JSONApiOkResponse{
			Data: response,
			Meta: JSONAPIMeta{
				TotalRecords: len(response.Items),
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

func SmartSuggestions(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		branchId := r.URL.Query().Get("branch_id")
		orderId := r.URL.Query().Get("order_id")

		limit := 10
		limitStr := r.URL.Query().Get("limit")
		if limitStr != "" {
			if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
				limit = parsedLimit
			}
		}

		aiSvc := services.AIService{
			Logger: logger,
			Config: config,
		}

		suggestions, err := aiSvc.GetSmartSuggestions(branchId, orderId, limit)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(JSONApiOkResponse{
			Data: suggestions,
			Meta: JSONAPIMeta{
				TotalRecords: len(suggestions),
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
