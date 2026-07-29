package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TransferRequest struct {
	FromTableID string `json:"from_table_id"`
	ToTableID   string `json:"to_table_id"`
}

type MergeRequest struct {
	FromTableID string `json:"from_table_id"`
	ToTableID   string `json:"to_table_id"`
}

func TransferTable(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req TransferRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		client, err := common.GetDatabaseClient(logger, &config)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db := client.Database(config.Databases[0].Database)
		tablesCol := db.Collection("tables")
		ordersCol := db.Collection("orders")

		var fromTable struct {
			Id      string `bson:"id"`
			OrderId string `bson:"order_id"`
			Status  string `bson:"status"`
		}
		err = tablesCol.FindOne(ctx, bson.M{"id": req.FromTableID}).Decode(&fromTable)
		if err != nil {
			http.Error(w, "source table not found", http.StatusNotFound)
			return
		}

		var toTable struct {
			Id     string `bson:"id"`
			Status string `bson:"status"`
		}
		err = tablesCol.FindOne(ctx, bson.M{"id": req.ToTableID}).Decode(&toTable)
		if err != nil {
			http.Error(w, "target table not found", http.StatusNotFound)
			return
		}

		if fromTable.Status == "available" {
			http.Error(w, "source table has no active order", http.StatusBadRequest)
			return
		}

		now := time.Now()

		if fromTable.OrderId != "" {
			ordersCol.UpdateOne(ctx, bson.M{"_id": fromTable.OrderId}, bson.M{
				"$set": bson.M{"table_id": req.ToTableID, "updated_at": now},
			})
		}

		tablesCol.UpdateOne(ctx, bson.M{"id": req.FromTableID}, bson.M{
			"$set": bson.M{"status": "available", "order_id": "", "updated_at": now},
		})

		tablesCol.UpdateOne(ctx, bson.M{"id": req.ToTableID}, bson.M{
			"$set": bson.M{"status": "occupied", "order_id": fromTable.OrderId, "updated_at": now},
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{
			Data: map[string]string{"message": "table transferred successfully"},
			Meta: JSONAPIMeta{TotalRecords: 1},
		})
	}
}

func MergeTables(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MergeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		client, err := common.GetDatabaseClient(logger, &config)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db := client.Database(config.Databases[0].Database)
		tablesCol := db.Collection("tables")

		var fromTable struct {
			Id      string `bson:"id"`
			OrderId string `bson:"order_id"`
			Status  string `bson:"status"`
		}
		err = tablesCol.FindOne(ctx, bson.M{"id": req.FromTableID}).Decode(&fromTable)
		if err != nil {
			http.Error(w, "source table not found", http.StatusNotFound)
			return
		}

		var toTable struct {
			Id      string `bson:"id"`
			OrderId string `bson:"order_id"`
			Status  string `bson:"status"`
		}
		err = tablesCol.FindOne(ctx, bson.M{"id": req.ToTableID}).Decode(&toTable)
		if err != nil {
			http.Error(w, "target table not found", http.StatusNotFound)
			return
		}

		now := time.Now()

		if fromTable.OrderId != "" && toTable.OrderId != "" {
			ordersCol := db.Collection("orders")
			ordersCol.UpdateOne(ctx, bson.M{"_id": fromTable.OrderId}, bson.M{
				"$set": bson.M{"table_id": req.ToTableID, "updated_at": now},
			})
		}

		tablesCol.UpdateOne(ctx, bson.M{"id": req.FromTableID}, bson.M{
			"$set": bson.M{"status": "available", "order_id": "", "updated_at": now},
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONApiOkResponse{
			Data: map[string]string{"message": "tables merged successfully"},
			Meta: JSONAPIMeta{TotalRecords: 1},
		})
	}
}

func GetFloorPlan(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client, err := common.GetDatabaseClient(logger, &config)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db := client.Database(config.Databases[0].Database)

		branchID := r.URL.Query().Get("branch_id")
		filter := bson.M{}
		if branchID != "" {
			filter["branch_id"] = branchID
		}

		cursor, err := db.Collection("tables").Find(ctx, filter)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		var tables []bson.M
		if err := cursor.All(ctx, &tables); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		summary := map[string]interface{}{
			"tables":     tables,
			"total":      len(tables),
			"available":  0,
			"occupied":   0,
			"reserved":   0,
			"cleaning":   0,
		}

		for _, t := range tables {
			if status, ok := t["status"].(string); ok {
				switch status {
				case "available":
					summary["available"] = summary["available"].(int) + 1
				case "occupied":
					summary["occupied"] = summary["occupied"].(int) + 1
				case "reserved":
					summary["reserved"] = summary["reserved"].(int) + 1
				case "cleaning":
					summary["cleaning"] = summary["cleaning"].(int) + 1
				}
			}
		}

		jsonResponse, _ := json.Marshal(JSONApiOkResponse{
			Data: summary,
			Meta: JSONAPIMeta{TotalRecords: len(tables)},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}
