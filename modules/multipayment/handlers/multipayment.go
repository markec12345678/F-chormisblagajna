package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	mp_models "github.com/nutrixpos/pos/modules/multipayment/models"
	"github.com/nutrixpos/pos/modules/multipayment/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

func AddPayment(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payment mp_models.PaymentPart
		if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if payment.OrderId == "" || payment.Amount <= 0 {
			http.Error(w, "order_id and positive amount are required", http.StatusBadRequest)
			return
		}

		svc := services.MultiPaymentService{Logger: log, Config: cfg}
		if err := svc.AddPayment(&payment); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to add payment", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: payment})
	}
}

func GetPaymentsByOrder(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderId := vars["orderId"]

		svc := services.MultiPaymentService{Logger: log, Config: cfg}
		payments, err := svc.GetPaymentsByOrder(orderId)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to load payments", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: payments})
	}
}

func GetPaymentSummary(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderId := vars["orderId"]

		totalDue := 0.0
		if t := r.URL.Query().Get("total"); t != "" {
			json.Unmarshal([]byte(t), &totalDue)
		}

		svc := services.MultiPaymentService{Logger: log, Config: cfg}
		summary, err := svc.GetPaymentSummary(orderId, totalDue)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "summary query failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: summary})
	}
}

func GetDailyPayments(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startDate := r.URL.Query().Get("start_date")
		endDate := r.URL.Query().Get("end_date")

		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			start = time.Now().AddDate(0, 0, -7)
		}

		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			end = time.Now()
		}

		svc := services.MultiPaymentService{Logger: log, Config: cfg}
		payments, err := svc.GetDailyPayments(start, end)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to load daily payments", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: payments})
	}
}

func SettleOrder(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderId := vars["orderId"]

		var body struct{ TotalDue float64 `json:"total_due"` }
		json.NewDecoder(r.Body).Decode(&body)

		svc := services.MultiPaymentService{Logger: log, Config: cfg}
		if err := svc.SettleOrder(orderId, body.TotalDue); err != nil {
			log.Error(err.Error())
			http.Error(w, "settle failed", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "settled", "order_id": orderId}})
	}
}

func RefundPayment(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paymentId := vars["id"]

		svc := services.MultiPaymentService{Logger: log, Config: cfg}
		if err := svc.RefundPayment(paymentId); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to refund payment", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: map[string]string{"status": "refunded"}})
	}
}
