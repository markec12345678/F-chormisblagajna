package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	chat_models "github.com/nutrixpos/pos/modules/chat/models"
	"github.com/nutrixpos/pos/modules/chat/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

func GetChannels(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.ChatService{Logger: log, Config: cfg}
		channels, err := svc.GetChannels()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to load channels", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: channels})
	}
}

func CreateChannel(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var channel chat_models.ChatChannel
		if err := json.NewDecoder(r.Body).Decode(&channel); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if channel.Name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		svc := services.ChatService{Logger: log, Config: cfg}
		if err := svc.CreateChannel(&channel); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to create channel", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: channel})
	}
}

func GetMessages(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		channel := r.URL.Query().Get("channel")

		svc := services.ChatService{Logger: log, Config: cfg}
		messages, err := svc.GetMessages(channel, 100)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to load messages", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: messages})
	}
}

func SendMessage(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var msg chat_models.ChatMessage
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if msg.Content == "" || msg.Sender == "" {
			http.Error(w, "content and sender are required", http.StatusBadRequest)
			return
		}

		svc := services.ChatService{Logger: log, Config: cfg}
		if err := svc.SendMessage(&msg); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to send message", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: msg})
	}
}
