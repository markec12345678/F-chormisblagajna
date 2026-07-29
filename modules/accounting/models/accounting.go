package models

type QuickBooksTransaction struct {
	Date        string  `json:"date"`
	Type        string  `json:"type"`
	Num         string  `json:"num"`
	Name        string  `json:"name"`
	Memo        string  `json:"memo"`
	Amount      float64 `json:"amount"`
	Account     string  `json:"account"`
	IncomeAcct  string  `json:"income_account"`
	TaxRate     float64 `json:"tax_rate"`
}

type XeroTransaction struct {
	Date         string  `json:"date"`
	Reference    string  `json:"reference"`
	Payee        string  `json:"payee"`
	Description  string  `json:"description"`
	LineAmount   float64 `json:"line_amount"`
	TaxAmount    float64 `json:"tax_amount"`
	AccountCode  string  `json:"account_code"`
	TaxType      string  `json:"tax_type"`
	Category     string  `json:"category"`
}
