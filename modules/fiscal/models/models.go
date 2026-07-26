package models

import "time"

// FiscalSettings stores FURS ZAPOS configuration for fiscal receipt generation.
type FiscalSettings struct {
	Enabled            bool   `bson:"enabled" json:"enabled" mapstructure:"enabled"`
	TaxNumber          int    `bson:"tax_number" json:"tax_number" mapstructure:"tax_number"`
	BusinessPremiseID  string `bson:"business_premise_id" json:"business_premise_id" mapstructure:"business_premise_id"`
	ElectronicDeviceID string `bson:"electronic_device_id" json:"electronic_device_id" mapstructure:"electronic_device_id"`
	SoftwareSupplierTN int    `bson:"software_supplier_tax_number" json:"software_supplier_tax_number" mapstructure:"software_supplier_tax_number"`
	// CertificatePath is the path to the PKCS#12 (.p12) certificate file.
	CertificatePath string `bson:"certificate_path" json:"certificate_path" mapstructure:"certificate_path"`
	// CertificatePassword is the password for the PKCS#12 certificate.
	CertificatePassword string `bson:"certificate_password" json:"certificate_password" mapstructure:"certificate_password"`
	// Environment: "test" or "production"
	Environment string `bson:"environment" json:"environment" mapstructure:"environment"`
	// BusinessPremiseType: "real_estate" or "movable"
	BusinessPremiseType string `bson:"business_premise_type" json:"business_premise_type" mapstructure:"business_premise_type"`
	// Address fields for real estate business premise
	Street              string `bson:"street" json:"street" mapstructure:"street"`
	HouseNumber         string `bson:"house_number" json:"house_number" mapstructure:"house_number"`
	HouseNumberAdditional string `bson:"house_number_additional" json:"house_number_additional" mapstructure:"house_number_additional"`
	City                string `bson:"city" json:"city" mapstructure:"city"`
	Community           string `bson:"community" json:"community" mapstructure:"community"`
	PostalCode          string `bson:"postal_code" json:"postal_code" mapstructure:"postal_code"`
	CadastralNumber     int    `bson:"cadastral_number" json:"cadastral_number" mapstructure:"cadastral_number"`
	BuildingNumber      int    `bson:"building_number" json:"building_number" mapstructure:"building_number"`
	BuildingSectionNumber int  `bson:"building_section_number" json:"building_section_number" mapstructure:"building_section_number"`
	// InvoiceNumber tracks the next invoice number for this device
	InvoiceNumber       int    `bson:"invoice_number" json:"invoice_number" mapstructure:"invoice_number"`
}

// InvoiceItem represents a single line item on a fiscal invoice.
type InvoiceItem struct {
	Name         string  `json:"name"`
	Quantity     float64 `json:"quantity"`
	UnitPrice    float64 `json:"unit_price"`
	TaxRate      float64 `json:"tax_rate"`
	TaxableAmount float64 `json:"taxable_amount"`
	TaxAmount    float64 `json:"tax_amount"`
}

// InvoiceRequest is the full fiscal invoice request sent to FURS.
type InvoiceRequest struct {
	Header struct {
		MessageID string `json:"MessageID"`
		DateTime  string `json:"DateTime"`
	} `json:"InvoiceRequestHeader"`
	Invoice struct {
		TaxNumber         int    `json:"TaxNumber"`
		IssueDateTime     string `json:"IssueDateTime"`
		NumberingStructure string `json:"NumberingStructure"`
		InvoiceIdentifier struct {
			BusinessPremiseID  string `json:"BusinessPremiseID"`
			ElectronicDeviceID string `json:"ElectronicDeviceID"`
			InvoiceNumber      string `json:"InvoiceNumber"`
		} `json:"InvoiceIdentifier"`
		InvoiceAmount  float64 `json:"InvoiceAmount"`
		PaymentAmount  float64 `json:"PaymentAmount"`
		ProtectedID    string  `json:"ProtectedID"`
		TaxesPerSeller []struct {
			VAT []struct {
				TaxRate       float64 `json:"TaxRate"`
				TaxableAmount float64 `json:"TaxableAmount"`
				TaxAmount     float64 `json:"TaxAmount"`
			} `json:"VAT"`
		} `json:"TaxesPerSeller"`
		OperatorTaxNumber   int    `json:"OperatorTaxNumber"`
		ForeignOperator     bool   `json:"ForeignOperator"`
		SubsequentSubmit    bool   `json:"SubsequentSubmit"`
		SpecialNotes        string `json:"SpecialNotes"`
	} `json:"Invoice"`
}

// InvoiceResponse represents the FURS response containing the Unique Invoice ID (EOR).
type InvoiceResponse struct {
	UniqueInvoiceID string `json:"UniqueInvoiceID"`
}

// FiscalReceipt is the stored record of a fiscalized receipt.
type FiscalReceipt struct {
	ID              string    `json:"id" bson:"id,omitempty"`
	OrderID         string    `json:"order_id" bson:"order_id"`
	EOR             string    `json:"eor" bson:"eor"`
	ZOI             string    `json:"zoi" bson:"zoi"`
	QRData          string    `json:"qr_data" bson:"qr_data"`
	InvoiceNumber   string    `json:"invoice_number" bson:"invoice_number"`
	InvoiceAmount   float64   `json:"invoice_amount" bson:"invoice_amount"`
	IssuedAt        time.Time `json:"issued_at" bson:"issued_at"`
	PendingOffline  bool      `json:"pending_offline" bson:"pending_offline"`
	RetryCount      int       `json:"retry_count" bson:"retry_count"`
}
