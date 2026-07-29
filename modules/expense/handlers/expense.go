package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	expense_models "github.com/nutrixpos/pos/modules/expense/models"
	"github.com/nutrixpos/pos/modules/expense/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

func GetAllExpenses(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.ExpenseService{Logger: log, Config: cfg}
		expenses, err := svc.GetAllExpenses()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: expenses})
	}
}

func CreateExpense(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var expense expense_models.Expense
		if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		svc := services.ExpenseService{Logger: log, Config: cfg}
		if err := svc.CreateExpense(&expense); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to create expense", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: expense})
	}
}

func UpdateExpense(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var expense expense_models.Expense
		if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		svc := services.ExpenseService{Logger: log, Config: cfg}
		if err := svc.UpdateExpense(id, &expense); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to update expense", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func DeleteExpense(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		svc := services.ExpenseService{Logger: log, Config: cfg}
		if err := svc.DeleteExpense(id); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to delete expense", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func GetExpenseSummary(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startDate := r.URL.Query().Get("start_date")
		endDate := r.URL.Query().Get("end_date")

		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			start = time.Now().AddDate(0, 0, -30)
		}

		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			end = time.Now()
		}

		svc := services.ExpenseService{Logger: log, Config: cfg}
		summary, err := svc.GetExpenseSummary(start, end)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "summary query failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: summary})
	}
}
