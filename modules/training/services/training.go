package services

import (
	"context"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	tr_models "github.com/nutrixpos/pos/modules/training/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TrainingService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *TrainingService) getCollection() (*mongo.Collection, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database).Collection("training_sessions"), nil
}

func (s *TrainingService) GetModules() []tr_models.TrainingModule {
	return []tr_models.TrainingModule{
		{Key: "cashier", Name: "Cashier Training", Description: "Learn to process orders and payments", Icon: "pi pi-desktop", Steps: 5},
		{Key: "kitchen", Name: "Kitchen Training", Description: "Learn to manage orders and cooking workflow", Icon: "fa fa-kitchen-set", Steps: 4},
		{Key: "inventory", Name: "Inventory Training", Description: "Learn stock management and transfers", Icon: "pi pi-box", Steps: 4},
		{Key: "admin", Name: "Admin Training", Description: "Learn system configuration and management", Icon: "pi pi-cog", Steps: 5},
	}
}

func (s *TrainingService) GetSteps(module string) []tr_models.TrainingStep {
	allSteps := map[string][]tr_models.TrainingStep{
		"cashier": {
			{Id: "cs-1", Module: "cashier", Title: "Open a New Order", Description: "Click the 'New Order' button on the cashier screen", Action: "navigate_new_order", ExpectedOut: "Order creation screen opens"},
			{Id: "cs-2", Module: "cashier", Title: "Add Items", Description: "Search and add items to the order from the menu", Action: "add_items", ExpectedOut: "Items appear in the order list"},
			{Id: "cs-3", Module: "cashier", Title: "Apply Discount", Description: "Apply a discount to the current order", Action: "apply_discount", ExpectedOut: "Discounted total is shown"},
			{Id: "cs-4", Module: "cashier", Title: "Process Payment", Description: "Complete the order by processing the payment", Action: "process_payment", ExpectedOut: "Order is marked as paid"},
			{Id: "cs-5", Module: "cashier", Title: "Issue Receipt", Description: "Print or email the receipt to the customer", Action: "issue_receipt", ExpectedOut: "Receipt is generated"},
		},
		"kitchen": {
			{Id: "kt-1", Module: "kitchen", Title: "View Incoming Orders", Description: "Check the kitchen display for new orders", Action: "view_orders", ExpectedOut: "Orders appear in the queue"},
			{Id: "kt-2", Module: "kitchen", Title: "Update Order Status", Description: "Mark an order as 'preparing' when you start working on it", Action: "update_status", ExpectedOut: "Status changes to preparing"},
			{Id: "kt-3", Module: "kitchen", Title: "Complete an Order", Description: "Mark the order as 'ready' when done", Action: "complete_order", ExpectedOut: "Order is moved to ready"},
			{Id: "kt-4", Module: "kitchen", Title: "Flag Low Stock", Description: "Report a low-stock ingredient to the manager", Action: "flag_stock", ExpectedOut: "Alert is sent to inventory"},
		},
		"inventory": {
			{Id: "inv-1", Module: "inventory", Title: "View Stock Levels", Description: "Navigate to the inventory page and review stock", Action: "view_stock", ExpectedOut: "Current stock levels are displayed"},
			{Id: "inv-2", Module: "inventory", Title: "Receive Shipment", Description: "Add new stock by receiving a shipment", Action: "receive_shipment", ExpectedOut: "Stock quantities are updated"},
			{Id: "inv-3", Module: "inventory", Title: "Transfer Stock", Description: "Transfer items between storage locations", Action: "transfer_stock", ExpectedOut: "Transfer is recorded"},
			{Id: "inv-4", Module: "inventory", Title: "Perform Count", Description: "Do a manual inventory count and adjust quantities", Action: "perform_count", ExpectedOut: "Adjustment is logged"},
		},
		"admin": {
			{Id: "ad-1", Module: "admin", Title: "Add a User", Description: "Create a new staff user account", Action: "add_user", ExpectedOut: "User is created with defined role"},
			{Id: "ad-2", Module: "admin", Title: "Configure Settings", Description: "Change a system setting like tax rate or currency", Action: "configure_settings", ExpectedOut: "Setting is saved"},
			{Id: "ad-3", Module: "admin", Title: "View Reports", Description: "Access the sales or inventory reports", Action: "view_reports", ExpectedOut: "Report data is displayed"},
			{Id: "ad-4", Module: "admin", Title: "Manage Products", Description: "Add or edit a product in the product catalog", Action: "manage_products", ExpectedOut: "Product is updated"},
			{Id: "ad-5", Module: "admin", Title: "Review Audit Log", Description: "Check recent activity in the audit log", Action: "review_audit", ExpectedOut: "Audit entries are displayed"},
		},
	}

	if steps, ok := allSteps[module]; ok {
		return steps
	}
	return []tr_models.TrainingStep{}
}

func (s *TrainingService) StartSession(userId, module string) (*tr_models.TrainingSession, error) {
	coll, err := s.getCollection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session := &tr_models.TrainingSession{
		Id:         bson.NewObjectID().Hex(),
		UserId:     userId,
		Module:     module,
		StartedAt:  time.Now(),
		Score:      0,
		MaxScore:   100,
		Completed:  false,
		StepsDone:  0,
		TotalSteps: len(s.GetSteps(module)),
	}

	_, err = coll.InsertOne(ctx, session)
	return session, err
}

func (s *TrainingService) CompleteStep(sessionId string) error {
	coll, err := s.getCollection()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = coll.UpdateOne(ctx, bson.M{"id": sessionId, "completed": false}, bson.M{
		"$inc": bson.M{"steps_done": 1, "score": 20},
	})
	return err
}

func (s *TrainingService) CompleteSession(sessionId string) error {
	coll, err := s.getCollection()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	_, err = coll.UpdateOne(ctx, bson.M{"id": sessionId, "completed": false}, bson.M{
		"$set": bson.M{"completed": true, "ended_at": now, "steps_done": 0},
		"$inc": bson.M{"score": 0},
	})

	var session tr_models.TrainingSession
	err2 := coll.FindOne(ctx, bson.M{"id": sessionId}).Decode(&session)
	if err2 == nil {
		steps := len(s.GetSteps(session.Module))
		_, _ = coll.UpdateOne(ctx, bson.M{"id": sessionId}, bson.M{
			"$set": bson.M{"total_steps": steps, "steps_done": steps},
		})
	}

	return err
}

func (s *TrainingService) GetUserSessions(userId string) ([]tr_models.TrainingSession, error) {
	coll, err := s.getCollection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, bson.M{"user_id": userId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	sessions := make([]tr_models.TrainingSession, 0)
	for cursor.Next(ctx) {
		var sess tr_models.TrainingSession
		if err := cursor.Decode(&sess); err != nil {
			continue
		}
		sessions = append(sessions, sess)
	}

	return sessions, nil
}

func (s *TrainingService) GetUserProgress(userId string) ([]tr_models.TrainingProgress, error) {
	sessions, err := s.GetUserSessions(userId)
	if err != nil {
		return nil, err
	}

	progressMap := make(map[string]*tr_models.TrainingProgress)
	for _, sess := range sessions {
		if _, exists := progressMap[sess.Module]; !exists {
			pct := 0.0
			if sess.TotalSteps > 0 {
				pct = float64(sess.StepsDone) / float64(sess.TotalSteps) * 100
			}
			progressMap[sess.Module] = &tr_models.TrainingProgress{
				SessionId:     sess.Id,
				Module:        sess.Module,
				StartedAt:     sess.StartedAt.Format("2006-01-02 15:04"),
				StepsDone:     sess.StepsDone,
				TotalSteps:    sess.TotalSteps,
				Score:         sess.Score,
				MaxScore:      sess.MaxScore,
				CompletionPct: pct,
				Completed:     sess.Completed,
			}
		}
	}

	result := make([]tr_models.TrainingProgress, 0, len(progressMap))
	for _, p := range progressMap {
		result = append(result, *p)
	}

	return result, nil
}
