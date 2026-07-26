package models

import "time"

// FiscalSettingsHR stores Croatian Fina eRačun configuration.
type FiscalSettingsHR struct {
	Enabled            bool   `bson:"enabled" json:"enabled" mapstructure:"enabled"`
	OIB                string `bson:"oib" json:"oib" mapstructure:"oib"`
	BusinessPremiseID  string `bson:"business_premise_id" json:"business_premise_id" mapstructure:"business_premise_id"`
	ElectronicDeviceID string `bson:"electronic_device_id" json:"electronic_device_id" mapstructure:"electronic_device_id"`
	CertificatePath    string `bson:"certificate_path" json:"certificate_path" mapstructure:"certificate_path"`
	CertificatePassword string `bson:"certificate_password" json:"certificate_password" mapstructure:"certificate_password"`
	// Environment: "test" or "production"
	Environment string `bson:"environment" json:"environment" mapstructure:"environment"`
	// InvoiceNumber tracks the next invoice number
	InvoiceNumber int `bson:"invoice_number" json:"invoice_number" mapstructure:"invoice_number"`
	// OperatorOIB is the operator's OIB (if different from business OIB)
	OperatorOIB string `bson:"operator_oib" json:"operator_oib" mapstructure:"operator_oib"`
}

// InvoiceItemHR represents a single line item on a Croatian fiscal invoice.
type InvoiceItemHR struct {
	Name          string  `json:"name"`
	Quantity      float64 `json:"quantity"`
	UnitPrice     float64 `json:"unit_price"`
	TaxRate       float64 `json:"tax_rate"`
	TaxableAmount float64 `json:"taxable_amount"`
	TaxAmount     float64 `json:"tax_amount"`
}

// InvoiceRequestHR is the full Croatian fiscal invoice request.
type InvoiceRequestHR struct {
	OrderID       string          `json:"order_id"`
	Items         []InvoiceItemHR `json:"items"`
	TotalAmount   float64         `json:"total_amount"`
	PaymentMethod string          `json:"payment_method"` // G=cash, K=card, T=transfer, O=other
	OperatorOIB   string          `json:"operator_oib"`
}

// InvoiceResponseHR represents the Croatian CIS response containing the JIR.
type InvoiceResponseHR struct {
	JIR    string `json:"jir"`
	Status string `json:"status"`
}

// FiscalReceiptHR is the stored record of a Croatian fiscalized receipt.
type FiscalReceiptHR struct {
	ID            string    `json:"id" bson:"id,omitempty"`
	OrderID       string    `json:"order_id" bson:"order_id"`
	JIR           string    `json:"jir" bson:"jir"`
	ZKI           string    `json:"zki" bson:"zki"`
	InvoiceNumber string    `json:"invoice_number" bson:"invoice_number"`
	InvoiceAmount float64   `json:"invoice_amount" bson:"invoice_amount"`
	IssuedAt      time.Time `json:"issued_at" bson:"issued_at"`
	PendingOffline bool     `json:"pending_offline" bson:"pending_offline"`
	RetryCount    int       `json:"retry_count" bson:"retry_count"`
}
