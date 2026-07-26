package models

type VoiceOrderRequest struct {
	AudioBase64 string `json:"audio_base64"`
	Language    string `json:"language"`
	BranchId    string `json:"branch_id"`
}

type VoiceOrderResponse struct {
	Transcript  string       `json:"transcript"`
	Items       []AIOrderItem `json:"items"`
	Confidence  float64      `json:"confidence"`
	Suggestions []string     `json:"suggestions"`
}

type AIOrderItem struct {
	ProductName string  `json:"product_name"`
	Quantity    float64 `json:"quantity"`
	Comment     string  `json:"comment"`
	Confidence  float64 `json:"confidence"`
	ProductId   string  `json:"product_id"`
}

type SmartSuggestion struct {
	ProductId   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Reason      string  `json:"reason"`
	Score       float64 `json:"score"`
}

type AISearchRequest struct {
	Query    string `json:"query"`
	BranchId string `json:"branch_id"`
	Language string `json:"language"`
	Limit    int    `json:"limit"`
}

type AISearchResponse struct {
	Results     []SearchResult    `json:"results"`
	Suggestions []SmartSuggestion `json:"suggestions"`
}

type SearchResult struct {
	ProductId   string  `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Score       float64 `json:"score"`
	Category    string  `json:"category"`
}
