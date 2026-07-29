package models

type DisplayConfig struct {
	Id                string   `json:"id" bson:"id" mapstructure:"id"`
	DisplayName       string   `json:"display_name" bson:"display_name" mapstructure:"display_name"`
	ShowPromotions    bool     `json:"show_promotions" bson:"show_promotions" mapstructure:"show_promotions"`
	ShowMenu          bool     `json:"show_menu" bson:"show_menu" mapstructure:"show_menu"`
	ShowOrderStatus   bool     `json:"show_order_status" bson:"show_order_status" mapstructure:"show_order_status"`
	ShowWaitTime      bool     `json:"show_wait_time" bson:"show_wait_time" mapstructure:"show_wait_time"`
	AutoSlideInterval int      `json:"auto_slide_interval" bson:"auto_slide_interval" mapstructure:"auto_slide_interval"`
	PromotionIds      []string `json:"promotion_ids" bson:"promotion_ids" mapstructure:"promotion_ids"`
	MenuCategoryIds   []string `json:"menu_category_ids" bson:"menu_category_ids" mapstructure:"menu_category_ids"`
	Theme             string   `json:"theme" bson:"theme" mapstructure:"theme"`
	WelcomeMessage    string   `json:"welcome_message" bson:"welcome_message" mapstructure:"welcome_message"`
	Active            bool     `json:"active" bson:"active" mapstructure:"active"`
}

type DisplayItem struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

type DisplayContent struct {
	Items      []DisplayItem `json:"items"`
	Interval   int           `json:"interval"`
	WelcomeMsg string        `json:"welcome_message"`
	Theme      string        `json:"theme"`
}
