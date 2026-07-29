package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	ly_models "github.com/nutrixpos/pos/modules/loyalty/models"
	"github.com/nutrixpos/pos/modules/loyalty/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

func GetAllCards(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.LoyaltyService{Logger: log, Config: cfg}
		cards, err := svc.GetAllCards()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to load cards", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: cards})
	}
}

func CreateCard(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var card ly_models.LoyaltyCard
		if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		svc := services.LoyaltyService{Logger: log, Config: cfg}
		if err := svc.CreateCard(&card); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to create card", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: card})
	}
}

func AddPoints(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		cardId := vars["cardId"]
		var body struct {
			Amount float64 `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Amount <= 0 {
			http.Error(w, "positive amount is required", http.StatusBadRequest)
			return
		}
		svc := services.LoyaltyService{Logger: log, Config: cfg}
		if err := svc.AddPoints(cardId, body.Amount); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to add points", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "points_added"}})
	}
}

func RedeemPoints(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		cardId := vars["cardId"]
		var body struct {
			RewardId string `json:"reward_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RewardId == "" {
			http.Error(w, "reward_id is required", http.StatusBadRequest)
			return
		}
		svc := services.LoyaltyService{Logger: log, Config: cfg}
		redemption, err := svc.RedeemPoints(cardId, body.RewardId)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to redeem", http.StatusInternalServerError)
			return
		}
		if redemption == nil {
			http.Error(w, "insufficient points", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: redemption})
	}
}

func GetAllRewards(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.LoyaltyService{Logger: log, Config: cfg}
		rewards, err := svc.GetAllRewards()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to load rewards", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: rewards})
	}
}

func CreateReward(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reward ly_models.Reward
		if err := json.NewDecoder(r.Body).Decode(&reward); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		svc := services.LoyaltyService{Logger: log, Config: cfg}
		if err := svc.CreateReward(&reward); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to create reward", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: reward})
	}
}

func GetSettings(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.LoyaltyService{Logger: log, Config: cfg}
		settings := svc.GetDefaultSettings()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: settings})
	}
}
