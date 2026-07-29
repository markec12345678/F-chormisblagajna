package models

type KioskConfig struct {
	Id              string `json:"id" bson:"id"`
	Name            string `json:"name" bson:"name"`
	Location        string `json:"location" bson:"location"`
	ShowCategories  bool   `json:"show_categories" bson:"show_categories"`
	ShowImages      bool   `json:"show_images" bson:"show_images"`
	AllowCustomize  bool   `json:"allow_customize" bson:"allow_customize"`
	AutoSendToKitchen bool `json:"auto_send_to_kitchen" bson:"auto_send_to_kitchen"`
	Theme           string `json:"theme" bson:"theme"`
	Active          bool   `json:"active" bson:"active"`
}

type KioskOrder struct {
	Id        string        `json:"id" bson:"id"`
	KioskId   string        `json:"kiosk_id" bson:"kiosk_id"`
	Items     []KioskOrderItem `json:"items" bson:"items"`
	Total     float64       `json:"total" bson:"total"`
	Status    string        `json:"status" bson:"status"`
	CreatedAt string        `json:"created_at" bson:"created_at"`
}

type KioskOrderItem struct {
	ProductId   string  `json:"product_id" bson:"product_id"`
	ProductName string  `json:"product_name" bson:"product_name"`
	Quantity    int     `json:"quantity" bson:"quantity"`
	UnitPrice   float64 `json:"unit_price" bson:"unit_price"`
	Notes       string  `json:"notes,omitempty" bson:"notes,omitempty"`
}
