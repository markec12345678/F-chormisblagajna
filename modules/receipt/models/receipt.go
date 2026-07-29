package models

type ReceiptTemplate struct {
	Id              string            `json:"id" bson:"id" mapstructure:"id"`
	Name            string            `json:"name" bson:"name" mapstructure:"name"`
	BusinessName    string            `json:"business_name" bson:"business_name" mapstructure:"business_name"`
	BusinessAddress string            `json:"business_address" bson:"business_address" mapstructure:"business_address"`
	BusinessPhone   string            `json:"business_phone" bson:"business_phone" mapstructure:"business_phone"`
	BusinessTaxId   string            `json:"business_tax_id" bson:"business_tax_id" mapstructure:"business_tax_id"`
	Header          string            `json:"header" bson:"header" mapstructure:"header"`
	Footer          string            `json:"footer" bson:"footer" mapstructure:"footer"`
	ShowLogo        bool              `json:"show_logo" bson:"show_logo" mapstructure:"show_logo"`
	ShowTaxId       bool              `json:"show_tax_id" bson:"show_tax_id" mapstructure:"show_tax_id"`
	ShowQRCode      bool              `json:"show_qr_code" bson:"show_qr_code" mapstructure:"show_qr_code"`
	ShowServer      bool              `json:"show_server" bson:"show_server" mapstructure:"show_server"`
	ShowTable       bool              `json:"show_table" bson:"show_table" mapstructure:"show_table"`
	PaperWidth      int               `json:"paper_width" bson:"paper_width" mapstructure:"paper_width"` // 58 or 80
	CustomFields    []CustomField     `json:"custom_fields" bson:"custom_fields" mapstructure:"custom_fields"`
	IsDefault       bool              `json:"is_default" bson:"is_default" mapstructure:"is_default"`
}

type CustomField struct {
	Key   string `json:"key" bson:"key" mapstructure:"key"`
	Value string `json:"value" bson:"value" mapstructure:"value"`
}

type PrintSettings struct {
	Id              string `json:"id" bson:"id" mapstructure:"id"`
	PrinterName     string `json:"printer_name" bson:"printer_name" mapstructure:"printer_name"`
	PrinterIP       string `json:"printer_ip" bson:"printer_ip" mapstructure:"printer_ip"`
	AutoPrint       bool   `json:"auto_print" bson:"auto_print" mapstructure:"auto_print"`
	PrintCopies     int    `json:"print_copies" bson:"print_copies" mapstructure:"print_copies"`
	TemplateId      string `json:"template_id" bson:"template_id" mapstructure:"template_id"`
	Connected       bool   `json:"connected" bson:"connected" mapstructure:"connected"`
}
