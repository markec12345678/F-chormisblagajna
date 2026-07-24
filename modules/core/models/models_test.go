package models

import (
	"encoding/json"
	"math"
	"testing"
	"time"
)

func TestJSONFloat_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    JSONFloat
		expected string
	}{
		{"positive", JSONFloat(1.5), "1.5"},
		{"negative", JSONFloat(-2.5), "-2.5"},
		{"zero", JSONFloat(0), "0"},
		{"infinity", JSONFloat(math.Inf(1)), "-1"},
		{"neg_infinity", JSONFloat(math.Inf(-1)), "-1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("MarshalJSON failed: %v", err)
			}
			if string(result) != tt.expected {
				t.Errorf("MarshalJSON(%v) = %s, want %s", tt.input, string(result), tt.expected)
			}
		})
	}
}

func TestJSONFloat_Unmarshal_Number(t *testing.T) {
	input := `{"sale_price": 12.5}`
	var result struct {
		SalePrice JSONFloat `json:"sale_price"`
	}
	if err := json.Unmarshal([]byte(input), &result); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if float64(result.SalePrice) != 12.5 {
		t.Errorf("SalePrice = %v, want 12.5", result.SalePrice)
	}
}

func TestJSONFloat_MarshalZero(t *testing.T) {
	result, err := json.Marshal(JSONFloat(0))
	if err != nil {
		t.Fatalf("MarshalJSON(0) failed: %v", err)
	}
	if string(result) != "0" {
		t.Errorf("MarshalJSON(0) = %s, want \"0\"", result)
	}
}

func TestProduct_JSON(t *testing.T) {
	product := Product{
		Id:       "prod-123",
		Name:     "Test Product",
		Price:    9.99,
		Quantity: 10,
		Ready:    5,
		Unit:     "kg",
		Materials: []Material{
			{Id: "mat-1", Name: "Flour", Quantity: 2},
		},
		SubProducts: []Product{
			{Id: "sub-1", Name: "Sub", Price: 3.0},
		},
		Entries: []ProductEntry{
			{Id: "e-1", Company: "Acme", Quantity: 100},
		},
		EnableInventoryConsumption: true,
		EnableFixedCost:            false,
		FixedCost:                  2.50,
	}

	data, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded Product
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Name != product.Name {
		t.Errorf("Name = %v, want %v", decoded.Name, product.Name)
	}
	if decoded.Price != product.Price {
		t.Errorf("Price = %v, want %v", decoded.Price, product.Price)
	}
	if len(decoded.Materials) != 1 {
		t.Errorf("Materials length = %v, want 1", len(decoded.Materials))
	}
	if len(decoded.SubProducts) != 1 {
		t.Errorf("SubProducts length = %v, want 1", len(decoded.SubProducts))
	}
	if !decoded.EnableInventoryConsumption {
		t.Error("EnableInventoryConsumption = false, want true")
	}
}

func TestOrder_JSON(t *testing.T) {
	now := time.Now()
	order := Order{
		Id:          "order-123",
		DisplayId:   "A-1",
		State:       "submitted",
		SubmittedAt: now,
		IsPaid:      false,
		Discount:    5.0,
		Tips:        2.5,
		IsDelivery:  true,
		IsTakeAway:  false,
		IsDineIn:    false,
		Customer:    Customer{Id: "c-1", Name: "Janez"},
		Items: []OrderItem{
			{Id: "item-1", Product: Product{Name: "Burger"}, Price: 9.99, Quantity: 2, Comment: "no pickles"},
		},
		CustomData: map[string]string{"table": "5"},
	}

	data, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded Order
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Id != order.Id {
		t.Errorf("Id = %v, want %v", decoded.Id, order.Id)
	}
	if decoded.State != order.State {
		t.Errorf("State = %v, want %v", decoded.State, order.State)
	}
	if len(decoded.Items) != 1 {
		t.Errorf("Items length = %v, want 1", len(decoded.Items))
	}
	if decoded.Items[0].Comment != "no pickles" {
		t.Errorf("Item comment = %v, want 'no pickles'", decoded.Items[0].Comment)
	}
	if decoded.Customer.Name != "Janez" {
		t.Errorf("Customer.Name = %v, want 'Janez'", decoded.Customer.Name)
	}
	if decoded.Discount != 5.0 {
		t.Errorf("Discount = %v, want 5.0", decoded.Discount)
	}
	if decoded.Tips != 2.5 {
		t.Errorf("Tips = %v, want 2.5", decoded.Tips)
	}
	if !decoded.IsDelivery {
		t.Error("IsDelivery = false, want true")
	}
	if decoded.CustomData["table"] != "5" {
		t.Errorf("CustomData[table] = %v, want '5'", decoded.CustomData["table"])
	}
}

func TestCustomer_JSON(t *testing.T) {
	customer := Customer{
		Id:      "c-1",
		Name:    "Janez Novak",
		Phone:   "+38640123456",
		Address: "Glavna 1, Ljubljana",
	}

	data, err := json.Marshal(customer)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded Customer
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Name != customer.Name {
		t.Errorf("Name = %v, want %v", decoded.Name, customer.Name)
	}
	if decoded.Phone != customer.Phone {
		t.Errorf("Phone = %v, want %v", decoded.Phone, customer.Phone)
	}
	if decoded.Address != customer.Address {
		t.Errorf("Address = %v, want %v", decoded.Address, customer.Address)
	}
}

func TestMaterial_JSON(t *testing.T) {
	mat := Material{
		Id:       "mat-1",
		Name:     "Flour",
		Quantity: 50,
		Unit:     "kg",
		Settings: MaterialSettings{StockAlertTreshold: 10},
		Entries: []MaterialEntry{
			{
				Id:               "e-1",
				PurchaseQuantity: 25,
				PurchasePrice:    5.99,
				Quantity:         25,
				Company:          "Acme Corp",
				SKU:              "FL-001",
			},
		},
	}

	data, err := json.Marshal(mat)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded Material
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Name != mat.Name {
		t.Errorf("Name = %v, want %v", decoded.Name, mat.Name)
	}
	if decoded.Settings.StockAlertTreshold != 10 {
		t.Errorf("StockAlertTreshold = %v, want 10", decoded.Settings.StockAlertTreshold)
	}
	if len(decoded.Entries) != 1 {
		t.Errorf("Entries length = %v, want 1", len(decoded.Entries))
	}
	if decoded.Entries[0].SKU != "FL-001" {
		t.Errorf("Entry SKU = %v, want 'FL-001'", decoded.Entries[0].SKU)
	}
}

func TestSettings_JSON(t *testing.T) {
	settings := Settings{
		Id:                    "settings-1",
		ShopMode:              "kitchen",
		AutoOpenCashDrawer:    true,
		Inventory:             MaterialSettings{StockAlertTreshold: 5},
		Language:              LanguageSettings{Code: "sl", Language: "Slovenščina"},
		ClientReceiptPrinter:  PrinterSettings{Host: "192.168.1.100"},
		KitchenReceiptPrinter: PrinterSettings{Host: "192.168.1.101"},
		PaymentSources:        []PaymentSource{{Name: "Cash"}, {Name: "Card"}},
		Orders: OrderSettings{
			Queues: []OrderQueueSettings{
				{Prefix: "A", Next: 1},
				{Prefix: "B", Next: 15},
			},
			DefaultCostCalculationMethod: "weighted_average",
		},
	}

	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded Settings
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Id != settings.Id {
		t.Errorf("Id = %v, want %v", decoded.Id, settings.Id)
	}
	if decoded.ShopMode != "kitchen" {
		t.Errorf("ShopMode = %v, want 'kitchen'", decoded.ShopMode)
	}
	if !decoded.AutoOpenCashDrawer {
		t.Error("AutoOpenCashDrawer = false, want true")
	}
	if len(decoded.PaymentSources) != 2 {
		t.Errorf("PaymentSources length = %v, want 2", len(decoded.PaymentSources))
	}
	if len(decoded.Orders.Queues) != 2 {
		t.Errorf("Queues length = %v, want 2", len(decoded.Orders.Queues))
	}
	if decoded.Language.Code != "sl" {
		t.Errorf("Language.Code = %v, want 'sl'", decoded.Language.Code)
	}
}

func TestCategory_JSON(t *testing.T) {
	category := Category{
		Id:   "cat-1",
		Name: "Drinks",
		Products: []Product{
			{Id: "p-1", Name: "Coffee", Price: 2.5},
			{Id: "p-2", Name: "Tea", Price: 2.0},
		},
	}

	data, err := json.Marshal(category)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded Category
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Name != category.Name {
		t.Errorf("Name = %v, want %v", decoded.Name, category.Name)
	}
	if len(decoded.Products) != 2 {
		t.Errorf("Products length = %v, want 2", len(decoded.Products))
	}
}

func TestDisposal_JSON(t *testing.T) {
	disposal := Disposal{
		Id:       "d-1",
		OrderId:  "o-1",
		Type:     TypeDisposalMaterial,
		Quantity: 3.5,
		Comment:  "expired",
	}

	data, err := json.Marshal(disposal)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded Disposal
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Type != TypeDisposalMaterial {
		t.Errorf("Type = %v, want %v", decoded.Type, TypeDisposalMaterial)
	}
	if decoded.Quantity != 3.5 {
		t.Errorf("Quantity = %v, want 3.5", decoded.Quantity)
	}
}

func TestMaterialDisposal_JSON(t *testing.T) {
	md := MaterialDisposal{
		Disposal:   Disposal{Id: "d-2", Type: TypeDisposalMaterial, Quantity: 1},
		MaterialId: "mat-5",
		EntryId:    "e-3",
	}

	data, err := json.Marshal(md)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded MaterialDisposal
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.MaterialId != "mat-5" {
		t.Errorf("MaterialId = %v, want 'mat-5'", decoded.MaterialId)
	}
	if decoded.EntryId != "e-3" {
		t.Errorf("EntryId = %v, want 'e-3'", decoded.EntryId)
	}
}

func TestDisposalConstants(t *testing.T) {
	if TypeDisposalMaterial != "disposal_material" {
		t.Errorf("TypeDisposalMaterial = %v, want 'disposal_material'", TypeDisposalMaterial)
	}
	if TypeDisposalProduct != "disposal_product" {
		t.Errorf("TypeDisposalProduct = %v, want 'disposal_product'", TypeDisposalProduct)
	}
}

func TestSalesPerDay_JSON(t *testing.T) {
	spd := SalesPerDay{
		Id:           "spd-1",
		Date:         "2026-01-15",
		TotalSales:   1500.50,
		Costs:        450.25,
		RefundsValue: 30.0,
		Orders: []SalesPerDayOrder{
			{
				Id:    "spo-1",
				Order: Order{Id: "o-1", DisplayId: "A-1"},
				Costs: []ItemCost{{ProductId: "p-1", Cost: 5.0, SalePrice: 9.99}},
			},
		},
		Refunds: []ItemRefund{
			{Id: "r-1", OrderId: "o-1", Amount: 9.99, Reason: "wrong order"},
		},
	}

	data, err := json.Marshal(spd)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded SalesPerDay
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.TotalSales != 1500.50 {
		t.Errorf("TotalSales = %v, want 1500.50", decoded.TotalSales)
	}
	if len(decoded.Orders) != 1 {
		t.Errorf("Orders length = %v, want 1", len(decoded.Orders))
	}
	if len(decoded.Refunds) != 1 {
		t.Errorf("Refunds length = %v, want 1", len(decoded.Refunds))
	}
	if decoded.Refunds[0].Reason != "wrong order" {
		t.Errorf("Refund Reason = %v, want 'wrong order'", decoded.Refunds[0].Reason)
	}
}

func TestItemCost_JSON(t *testing.T) {
	ic := ItemCost{
		ProductId:    "p-1",
		ItemId:       "i-1",
		ItemName:     "Burger",
		Cost:         4.50,
		SalePrice:    9.99,
		Quantity:     2,
		UseFixedCost: true,
		CostMethod:   "fixed",
		Components: []struct {
			ComponentName string  `json:"component_name" bson:"component_name" mapstructure:"component_name"`
			ComponentId   string  `json:"component_id" bson:"component_id" mapstructure:"component_id"`
			EntryId       string  `json:"entry_id" bson:"entry_id" mapstructure:"entry_id"`
			Quantity      float64 `json:"quantity" bson:"quantity" mapstructure:"quantity"`
			Cost          float64 `json:"cost" bson:"cost" mapstructure:"cost"`
		}{
			{ComponentName: "Bread", ComponentId: "c-1", EntryId: "e-1", Quantity: 2, Cost: 0.50},
			{ComponentName: "Meat", ComponentId: "c-2", EntryId: "e-2", Quantity: 0.2, Cost: 2.00},
		},
	}

	data, err := json.Marshal(ic)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded ItemCost
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.ItemName != "Burger" {
		t.Errorf("ItemName = %v, want 'Burger'", decoded.ItemName)
	}
	if !decoded.UseFixedCost {
		t.Error("UseFixedCost = false, want true")
	}
	if len(decoded.Components) != 2 {
		t.Errorf("Components length = %v, want 2", len(decoded.Components))
	}
}

func TestLogConstants(t *testing.T) {
	consts := map[string]string{
		LogTypeDisposalAdd:             "disposal_add",
		LogTypeMaterialInventoryReturn: "material_inventory_return",
		LogTypeOrderItemRefunded:       "order_item_refunded",
		LogTypeOrderStart:              "order_Start",
		LogTypeOrderFinish:             "order_finish",
		LogTypeMaterialConsume:         "component_consume",
		LogTypeMaterialAdd:             "component_add",
		LogTypeMaterialWaste:           "material_waste",
		LogTypeProductIncrease:         "product_increase",
		LogTypeSalesPerDayOrder:        "sales_per_day_order",
		LogTypeSalesPerDayRefund:       "sales_per_day_refund",
	}

	for got, want := range consts {
		if got != want {
			t.Errorf("Log constant = %q, want %q", got, want)
		}
	}
	if len(consts) != 11 {
		t.Errorf("Log constant count = %d, want 11", len(consts))
	}
}

func TestLogRefundOrder_JSON(t *testing.T) {
	lr := LogRefundOrder{
		Log:     Log{Id: "l-1", Type: LogTypeOrderItemRefunded, Date: time.Now(), UserId: "u-1"},
		Reason:  "wrong item",
		OrderId: "o-123",
	}

	data, err := json.Marshal(lr)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded LogRefundOrder
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Reason != "wrong item" {
		t.Errorf("Reason = %v, want 'wrong item'", decoded.Reason)
	}
	if decoded.OrderId != "o-123" {
		t.Errorf("OrderId = %v, want 'o-123'", decoded.OrderId)
	}
	if decoded.Type != LogTypeOrderItemRefunded {
		t.Errorf("Type = %v, want %v", decoded.Type, LogTypeOrderItemRefunded)
	}
}

func TestOrderDeliveryInfo_JSON(t *testing.T) {
	di := OrderDeliveryInfo{
		ReceiverName: "Janez",
		Address:      "Glavna 1, 1000 Ljubljana",
		PhoneNumber:  "+38640123456",
	}

	data, err := json.Marshal(di)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded OrderDeliveryInfo
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.ReceiverName != "Janez" {
		t.Errorf("ReceiverName = %v, want 'Janez'", decoded.ReceiverName)
	}
	if decoded.PhoneNumber != "+38640123456" {
		t.Errorf("PhoneNumber = %v, want '+38640123456'", decoded.PhoneNumber)
	}
}

func TestSubmitOrderMeta_JSON(t *testing.T) {
	meta := SubmitOrderMeta{
		IsPrintClientReceipt:  true,
		IsPrintKitchenReceipt: false,
	}

	data, err := json.Marshal(meta)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded SubmitOrderMeta
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if !decoded.IsPrintClientReceipt {
		t.Error("IsPrintClientReceipt = false, want true")
	}
	if decoded.IsPrintKitchenReceipt {
		t.Error("IsPrintKitchenReceipt = true, want false")
	}
}

func TestComponentConsumeLogs_JSON(t *testing.T) {
	ccl := ComponentConsumeLogs{
		Id:             "ccl-1",
		Date:           time.Now(),
		Name:           "Flour",
		Quantity:       2.5,
		Company:        "Acme",
		ItemName:       "Bread",
		OrderItemIndex: 0,
		OrderId:        "o-1",
		Type:           "production",
	}

	data, err := json.Marshal(ccl)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded ComponentConsumeLogs
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Name != "Flour" {
		t.Errorf("Name = %v, want 'Flour'", decoded.Name)
	}
	if decoded.Quantity != 2.5 {
		t.Errorf("Quantity = %v, want 2.5", decoded.Quantity)
	}
}

func TestLanguageData_JSON(t *testing.T) {
	ld := LanguageData{
		ApiVersion:  "1.0",
		Code:        "sl",
		Language:    "Slovenščina",
		Orientation: "ltr",
		Pack:        map[string]interface{}{"hello": "živjo"},
	}

	data, err := json.Marshal(ld)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded LanguageData
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Code != "sl" {
		t.Errorf("Code = %v, want 'sl'", decoded.Code)
	}
	if decoded.Pack["hello"] != "živjo" {
		t.Errorf("Pack[hello] = %v, want 'živjo'", decoded.Pack["hello"])
	}
}
