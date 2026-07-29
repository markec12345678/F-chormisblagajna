package services

import (
	"context"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	receipt_models "github.com/nutrixpos/pos/modules/receipt/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ReceiptService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *ReceiptService) GetTemplates() ([]receipt_models.ReceiptTemplate, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("receipt_templates")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	templates := make([]receipt_models.ReceiptTemplate, 0)
	for cursor.Next(ctx) {
		var tpl receipt_models.ReceiptTemplate
		if err := cursor.Decode(&tpl); err != nil {
			continue
		}
		templates = append(templates, tpl)
	}

	return templates, nil
}

func (s *ReceiptService) GetTemplate(id string) (*receipt_models.ReceiptTemplate, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("receipt_templates")

	var tpl receipt_models.ReceiptTemplate
	err = collection.FindOne(ctx, bson.M{"id": id}).Decode(&tpl)
	if err != nil {
		return nil, err
	}

	return &tpl, nil
}

func (s *ReceiptService) SaveTemplate(tpl *receipt_models.ReceiptTemplate) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("receipt_templates")

	if tpl.Id == "" {
		tpl.Id = bson.NewObjectID().Hex()
		_, err = collection.InsertOne(ctx, tpl)
	} else {
		_, err = collection.UpdateOne(ctx,
			bson.M{"id": tpl.Id},
			bson.M{"$set": tpl},
		)
	}

	return err
}

func (s *ReceiptService) DeleteTemplate(id string) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("receipt_templates")

	_, err = collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (s *ReceiptService) GetPrintSettings() (*receipt_models.PrintSettings, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("print_settings")

	var settings receipt_models.PrintSettings
	err = collection.FindOne(ctx, bson.M{}).Decode(&settings)
	if err != nil {
		return &receipt_models.PrintSettings{
			PrintCopies: 1,
			AutoPrint:   true,
		}, nil
	}

	return &settings, nil
}

func (s *ReceiptService) SavePrintSettings(settings *receipt_models.PrintSettings) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("print_settings")

	if settings.Id == "" {
		settings.Id = bson.NewObjectID().Hex()
		_, err = collection.InsertOne(ctx, settings)
	} else {
		_, err = collection.UpdateOne(ctx,
			bson.M{"id": settings.Id},
			bson.M{"$set": settings},
		)
	}

	return err
}

func (s *ReceiptService) PreviewReceipt(tpl *receipt_models.ReceiptTemplate, orderData map[string]interface{}) string {
	preview := ""

	if tpl.Header != "" {
		preview += tpl.Header + "\n"
		preview += "========================\n"
	}

	preview += tpl.BusinessName + "\n"
	if tpl.ShowTaxId && tpl.BusinessTaxId != "" {
		preview += "Tax ID: " + tpl.BusinessTaxId + "\n"
	}
	if tpl.BusinessAddress != "" {
		preview += tpl.BusinessAddress + "\n"
	}
	if tpl.BusinessPhone != "" {
		preview += tpl.BusinessPhone + "\n"
	}

	preview += "========================\n"
	preview += "\n"

	if items, ok := orderData["items"].([]map[string]interface{}); ok {
		for _, item := range items {
			name := ""
			qty := float64(0)
			price := float64(0)
			if n, ok := item["name"].(string); ok {
				name = n
			}
			if q, ok := item["quantity"].(float64); ok {
				qty = q
			}
			if p, ok := item["price"].(float64); ok {
				price = p
			}
			lineTotal := qty * price
			preview += name + "\n"
			preview += "  " + "  " + "  " + "\n"
			_ = lineTotal
		}
	}

	if total, ok := orderData["total"].(float64); ok {
		preview += "========================\n"
		preview += "TOTAL: " + "\n"
		_ = total
	}

	preview += "\n"

	if tpl.Footer != "" {
		preview += tpl.Footer + "\n"
	}

	return preview
}
