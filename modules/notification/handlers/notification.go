package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/notification/services"
)

type JSONAPIMeta struct {
	TotalRecords int `json:"total_records"`
}

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta JSONAPIMeta `json:"meta"`
}

func GetNotifications(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		unreadOnly := r.URL.Query().Get("unread") == "true"

		svc := services.NotificationService{Logger: logger, Config: config}
		notifications, err := svc.GetNotifications(userID, unreadOnly)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		jsonResponse, _ := json.Marshal(JSONApiOkResponse{Data: notifications, Meta: JSONAPIMeta{TotalRecords: len(notifications)}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func MarkAsRead(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		svc := services.NotificationService{Logger: logger, Config: config}
		err := svc.MarkAsRead(id)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{
			Data: map[string]string{"message": "marked as read"},
			Meta: JSONAPIMeta{TotalRecords: 1},
		})
	}
}

func MarkAllAsRead(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")

		svc := services.NotificationService{Logger: logger, Config: config}
		err := svc.MarkAllAsRead(userID)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{
			Data: map[string]string{"message": "all marked as read"},
			Meta: JSONAPIMeta{TotalRecords: 1},
		})
	}
}

func GetUnreadCount(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")

		svc := services.NotificationService{Logger: logger, Config: config}
		count, err := svc.GetUnreadCount(userID)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{
			Data: map[string]int64{"count": count},
			Meta: JSONAPIMeta{TotalRecords: 1},
		})
	}
}
