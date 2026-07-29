package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/training/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

func GetModules(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.TrainingService{Logger: log, Config: cfg}
		modules := svc.GetModules()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: modules})
	}
}

func GetSteps(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		module := vars["module"]

		svc := services.TrainingService{Logger: log, Config: cfg}
		steps := svc.GetSteps(module)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: steps})
	}
}

func StartSession(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			UserId string `json:"user_id"`
			Module string `json:"module"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Module == "" {
			http.Error(w, "module is required", http.StatusBadRequest)
			return
		}

		svc := services.TrainingService{Logger: log, Config: cfg}
		session, err := svc.StartSession(body.UserId, body.Module)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to start session", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: session})
	}
}

func CompleteStep(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sessionId := vars["id"]

		svc := services.TrainingService{Logger: log, Config: cfg}
		if err := svc.CompleteStep(sessionId); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to complete step", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "step_completed"}})
	}
}

func CompleteSession(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sessionId := vars["id"]

		svc := services.TrainingService{Logger: log, Config: cfg}
		if err := svc.CompleteSession(sessionId); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to complete session", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "completed"}})
	}
}

func GetUserProgress(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userId := vars["userId"]

		svc := services.TrainingService{Logger: log, Config: cfg}
		progress, err := svc.GetUserProgress(userId)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to load progress", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: progress})
	}
}

func GetUserSessions(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userId := vars["userId"]

		svc := services.TrainingService{Logger: log, Config: cfg}
		sessions, err := svc.GetUserSessions(userId)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to load sessions", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: sessions})
	}
}
