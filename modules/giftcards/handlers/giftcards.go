package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	gc_models "github.com/nutrixpos/pos/modules/giftcards/models"
	"github.com/nutrixpos/pos/modules/giftcards/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
}

func IssueCard(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var c gc_models.GiftCard
		if json.NewDecoder(r.Body).Decode(&c) != nil {
			http.Error(w, "invalid", http.StatusBadRequest)
			return
		}
		svc := services.GiftCardService{Logger: log, Config: cfg}
		if err := svc.Issue(&c); err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: c})
	}
}
func GetAllCards(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.GiftCardService{Logger: log, Config: cfg}
		d, _ := svc.GetAll()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: d})
	}
}
func RedeemCard(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct{ Code string `json:"code"`; Amount float64 `json:"amount"`; OrderId string `json:"order_id"` }
		if json.NewDecoder(r.Body).Decode(&body) != nil || body.Code == "" || body.Amount <= 0 {
			http.Error(w, "code+amount required", http.StatusBadRequest)
			return
		}
		svc := services.GiftCardService{Logger: log, Config: cfg}
		if err := svc.Redeem(body.Code, body.Amount, body.OrderId); err != nil {
			http.Error(w, "redeem failed", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "redeemed"}})
	}
}
func DeactivateCard(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.GiftCardService{Logger: log, Config: cfg}
		svc.Deactivate(mux.Vars(r)["id"])
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "deactivated"}})
	}
}
