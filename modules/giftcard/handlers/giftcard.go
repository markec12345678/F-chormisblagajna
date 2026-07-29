package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/giftcard/dto"
	"github.com/nutrixpos/pos/modules/giftcard/services"
)

type JSONAPIMeta struct {
	TotalRecords int `json:"total_records"`
}

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta JSONAPIMeta `json:"meta"`
}

func CreateGiftCard(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateGiftCardRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		svc := services.GiftCardService{Logger: logger, Config: config}
		card, err := svc.CreateGiftCard(req)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		jsonResponse, _ := json.Marshal(JSONApiOkResponse{Data: card, Meta: JSONAPIMeta{TotalRecords: 1}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func GetAllGiftCards(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.GiftCardService{Logger: logger, Config: config}
		cards, err := svc.GetAllGiftCards()
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		jsonResponse, _ := json.Marshal(JSONApiOkResponse{Data: cards, Meta: JSONAPIMeta{TotalRecords: len(cards)}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func GetGiftCard(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		code := vars["code"]

		svc := services.GiftCardService{Logger: logger, Config: config}
		card, err := svc.GetGiftCardByCode(code)
		if err != nil {
			http.Error(w, "gift card not found", http.StatusNotFound)
			return
		}

		jsonResponse, _ := json.Marshal(JSONApiOkResponse{Data: card, Meta: JSONAPIMeta{TotalRecords: 1}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func RedeemGiftCard(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.RedeemGiftCardRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		svc := services.GiftCardService{Logger: logger, Config: config}
		card, err := svc.RedeemGiftCard(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		jsonResponse, _ := json.Marshal(JSONApiOkResponse{Data: card, Meta: JSONAPIMeta{TotalRecords: 1}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func GetTransactions(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		svc := services.GiftCardService{Logger: logger, Config: config}
		txs, err := svc.GetTransactions(id)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		jsonResponse, _ := json.Marshal(JSONApiOkResponse{Data: txs, Meta: JSONAPIMeta{TotalRecords: len(txs)}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}
