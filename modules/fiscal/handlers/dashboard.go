package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type JSONAPIMeta struct {
	TotalRecords int `json:"total_records"`
}

type JSONApiOkResponse struct {
	Data interface{} `json:"data"`
	Meta JSONAPIMeta `json:"meta"`
}

func GetFiscalReceipts(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client, err := common.GetDatabaseClient(logger, &config)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db := client.Database(config.Databases[0].Database)

		startDate := r.URL.Query().Get("start_date")
		endDate := r.URL.Query().Get("end_date")

		filter := bson.M{}
		if startDate != "" || endDate != "" {
			dateFilter := bson.M{}
			if startDate != "" {
				t, err := time.Parse("2006-01-02", startDate)
				if err == nil {
					dateFilter["$gte"] = t
				}
			}
			if endDate != "" {
				t, err := time.Parse("2006-01-02", endDate)
				if err == nil {
					t = t.Add(24*time.Hour - time.Second)
					dateFilter["$lte"] = t
				}
			}
			filter["issued_at"] = dateFilter
		}

		opts := options.Find().SetSort(bson.M{"issued_at": -1}).SetLimit(100)
		cursor, err := db.Collection("fiscal_receipts").Find(ctx, filter, opts)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		var receipts []bson.M
		if err := cursor.All(ctx, &receipts); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		jsonResponse, _ := json.Marshal(JSONApiOkResponse{
			Data: receipts,
			Meta: JSONAPIMeta{TotalRecords: len(receipts)},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func GetFiscalDailySummary(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client, err := common.GetDatabaseClient(logger, &config)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db := client.Database(config.Databases[0].Database)

		dateStr := r.URL.Query().Get("date")
		if dateStr == "" {
			dateStr = time.Now().Format("2006-01-02")
		}

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			http.Error(w, "invalid date format", http.StatusBadRequest)
			return
		}

		dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
		dayEnd := dayStart.Add(24*time.Hour - time.Second)

		filter := bson.M{
			"issued_at": bson.M{"$gte": dayStart, "$lte": dayEnd},
		}

		count, err := db.Collection("fiscal_receipts").CountDocuments(ctx, filter)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		pipeline := []bson.M{
			{"$match": filter},
			{"$group": bson.M{
				"_id":          nil,
				"total_amount": bson.M{"$sum": "$invoice_amount"},
				"count":        bson.M{"$sum": 1},
			}},
		}

		cursor, err := db.Collection("fiscal_receipts").Aggregate(ctx, pipeline)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		var results []bson.M
		summary := map[string]interface{}{
			"date":        dateStr,
			"total_count": count,
			"total_amount": 0.0,
		}

		if cursor.All(ctx, &results) == nil && len(results) > 0 {
			if v, ok := results[0]["total_amount"].(float64); ok {
				summary["total_amount"] = v
			}
		}

		jsonResponse, _ := json.Marshal(JSONApiOkResponse{
			Data: summary,
			Meta: JSONAPIMeta{TotalRecords: 1},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func ExportFiscalReceiptsCSV(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client, err := common.GetDatabaseClient(logger, &config)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db := client.Database(config.Databases[0].Database)

		startDate := r.URL.Query().Get("start_date")
		endDate := r.URL.Query().Get("end_date")

		filter := bson.M{}
		if startDate != "" || endDate != "" {
			dateFilter := bson.M{}
			if startDate != "" {
				t, err := time.Parse("2006-01-02", startDate)
				if err == nil {
					dateFilter["$gte"] = t
				}
			}
			if endDate != "" {
				t, err := time.Parse("2006-01-02", endDate)
				if err == nil {
					t = t.Add(24*time.Hour - time.Second)
					dateFilter["$lte"] = t
				}
			}
			filter["issued_at"] = dateFilter
		}

		opts := options.Find().SetSort(bson.M{"issued_at": -1})
		cursor, err := db.Collection("fiscal_receipts").Find(ctx, filter, opts)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		var receipts []bson.M
		if err := cursor.All(ctx, &receipts); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment;filename=fiscal_receipts.csv")

		fmt.Fprintf(w, "Invoice Number,EOR,ZOI,Amount,Date\n")
		for _, r := range receipts {
			invoiceNum := ""
			if v, ok := r["invoice_number"].(string); ok {
				invoiceNum = v
			}
			eor := ""
			if v, ok := r["eor"].(string); ok {
				eor = v
			}
			zoi := ""
			if v, ok := r["zoi"].(string); ok {
				zoi = v
			}
			amount := 0.0
			if v, ok := r["invoice_amount"].(float64); ok {
				amount = v
			}
			issuedAt := ""
			if v, ok := r["issued_at"].(time.Time); ok {
				issuedAt = v.Format("2006-01-02 15:04:05")
			}
			fmt.Fprintf(w, "%s,%s,%s,%.2f,%s\n", invoiceNum, eor, zoi, amount, issuedAt)
		}
	}
}
