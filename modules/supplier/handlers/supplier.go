package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	supplier_models "github.com/nutrixpos/pos/modules/supplier/models"
	"github.com/nutrixpos/pos/modules/supplier/services"
)

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

func GetAllSuppliers(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.SupplierService{Logger: log, Config: cfg}
		suppliers, err := svc.GetAllSuppliers()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: suppliers})
	}
}

func GetSupplier(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		svc := services.SupplierService{Logger: log, Config: cfg}
		supplier, err := svc.GetSupplier(id)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "supplier not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: supplier})
	}
}

func CreateSupplier(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var supplier supplier_models.Supplier
		if err := json.NewDecoder(r.Body).Decode(&supplier); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		svc := services.SupplierService{Logger: log, Config: cfg}
		if err := svc.CreateSupplier(&supplier); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to create supplier", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: supplier})
	}
}

func UpdateSupplier(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var supplier supplier_models.Supplier
		if err := json.NewDecoder(r.Body).Decode(&supplier); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		svc := services.SupplierService{Logger: log, Config: cfg}
		if err := svc.UpdateSupplier(id, &supplier); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to update supplier", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func DeleteSupplier(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		svc := services.SupplierService{Logger: log, Config: cfg}
		if err := svc.DeleteSupplier(id); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to delete supplier", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func GetSupplierOrders(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		supplierId := vars["supplier_id"]

		svc := services.SupplierService{Logger: log, Config: cfg}
		orders, err := svc.GetSupplierOrders(supplierId)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: orders})
	}
}

func CreateSupplierOrder(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var order supplier_models.SupplierOrder
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		svc := services.SupplierService{Logger: log, Config: cfg}
		if err := svc.CreateSupplierOrder(&order); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to create order", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSONApiOkResponse{Data: order})
	}
}

func UpdateSupplierOrderStatus(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var req struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		svc := services.SupplierService{Logger: log, Config: cfg}
		if err := svc.UpdateSupplierOrderStatus(id, req.Status); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to update status", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
