package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	fb_models "github.com/nutrixpos/pos/modules/feedback/models"
	"github.com/nutrixpos/pos/modules/feedback/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

func SubmitFeedback(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fb fb_models.Feedback
		if err := json.NewDecoder(r.Body).Decode(&fb); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if fb.Rating < 1 || fb.Rating > 5 {
			http.Error(w, "rating must be between 1 and 5", http.StatusBadRequest)
			return
		}

		svc := services.FeedbackService{Logger: log, Config: cfg}
		if err := svc.SubmitFeedback(&fb); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to submit feedback", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: fb})
	}
}

func GetAllFeedbacks(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.FeedbackService{Logger: log, Config: cfg}
		feedbacks, err := svc.GetAllFeedbacks()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to load feedbacks", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: feedbacks})
	}
}

func GetFeedbackSummary(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.FeedbackService{Logger: log, Config: cfg}
		summary, err := svc.GetFeedbackSummary()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to load summary", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: summary})
	}
}

func RespondToFeedback(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var body struct {
			Response string `json:"response"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Response == "" {
			http.Error(w, "response is required", http.StatusBadRequest)
			return
		}

		svc := services.FeedbackService{Logger: log, Config: cfg}
		if err := svc.RespondToFeedback(id, body.Response); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to respond", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "responded"}})
	}
}

func DeleteFeedback(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		svc := services.FeedbackService{Logger: log, Config: cfg}
		if err := svc.DeleteFeedback(id); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to delete", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "deleted"}})
	}
}
