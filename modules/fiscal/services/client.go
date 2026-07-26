package services

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/nutrixpos/pos/modules/fiscal/models"
	"golang.org/x/crypto/pkcs12"
)

var (
	testBaseURL      = "https://blagajne-test.fu.gov.si:9002"
	productionBaseURL = "https://blagajne.fu.gov.si:9003"
)

// FURSClient handles communication with the FURS ZAPOS API.
type FURSClient struct {
	settings  *models.FiscalSettings
	privateKey *rsa.PrivateKey
	cert      *x509.Certificate
	httpClient *http.Client
	baseURL   string
}

// NewFURSClient creates a new FURS API client from fiscal settings.
func NewFURSClient(settings *models.FiscalSettings) (*FURSClient, error) {
	if settings.CertificatePath == "" {
		return nil, fmt.Errorf("certificate path is required")
	}

	p12Data, err := readCertificateFile(settings.CertificatePath)
	if err != nil {
		return nil, fmt.Errorf("read certificate: %w", err)
	}

	privateKey, certificate, err := pkcs12.Decode(p12Data, settings.CertificatePassword)
	if err != nil {
		return nil, fmt.Errorf("decode PKCS#12: %w", err)
	}

	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not RSA")
	}

	tlsCert := tls.Certificate{
		Certificate: [][]byte{certificate.Raw},
		PrivateKey:  rsaKey,
		Leaf:        certificate,
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		MinVersion:   tls.VersionTLS12,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	baseURL := testBaseURL
	if settings.Environment == "production" {
		baseURL = productionBaseURL
	}

	return &FURSClient{
		settings:   settings,
		privateKey: rsaKey,
		cert:       certificate,
		httpClient: httpClient,
		baseURL:    baseURL,
	}, nil
}

// NewFURSClientFromKey creates a client from a pre-loaded RSA key and certificate (for testing).
func NewFURSClientFromKey(settings *models.FiscalSettings, privateKey *rsa.PrivateKey, cert *x509.Certificate) *FURSClient {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12},
		},
		Timeout: 30 * time.Second,
	}

	baseURL := testBaseURL
	if settings.Environment == "production" {
		baseURL = productionBaseURL
	}

	return &FURSClient{
		settings:   settings,
		privateKey: privateKey,
		cert:       cert,
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

// FiscalizeInvoice sends an invoice to FURS and returns the EOR (Unique Invoice ID).
func (c *FURSClient) FiscalizeInvoice(ctx context.Context, orderID string, items []models.InvoiceItem,
	totalAmount float64, issueTime time.Time) (*models.InvoiceResponse, error) {

	invoiceNum := FormatInvoiceNumber(c.settings.InvoiceNumber)

	zoi, err := CalculateZOI(
		c.privateKey,
		c.settings.TaxNumber,
		issueTime,
		invoiceNum,
		c.settings.BusinessPremiseID,
		c.settings.ElectronicDeviceID,
		totalAmount,
	)
	if err != nil {
		return nil, fmt.Errorf("calculate ZOI: %w", err)
	}

	// Calculate tax breakdown
	taxesPerSeller := buildTaxesPerSeller(items)

	// Build the request payload
	payload := models.InvoiceRequest{}
	payload.Header.MessageID = fmt.Sprintf("order-%s-%d", orderID, issueTime.UnixMilli())
	payload.Header.DateTime = FormatIssueDateTime(issueTime)
	payload.Invoice.TaxNumber = c.settings.TaxNumber
	payload.Invoice.IssueDateTime = FormatIssueDateTime(issueTime)
	payload.Invoice.NumberingStructure = "B"
	payload.Invoice.InvoiceIdentifier.BusinessPremiseID = c.settings.BusinessPremiseID
	payload.Invoice.InvoiceIdentifier.ElectronicDeviceID = c.settings.ElectronicDeviceID
	payload.Invoice.InvoiceIdentifier.InvoiceNumber = invoiceNum
	payload.Invoice.InvoiceAmount = totalAmount
	payload.Invoice.PaymentAmount = totalAmount
	payload.Invoice.ProtectedID = zoi
	payload.Invoice.TaxesPerSeller = taxesPerSeller
	payload.Invoice.OperatorTaxNumber = c.settings.TaxNumber
	payload.Invoice.ForeignOperator = false
	payload.Invoice.SubsequentSubmit = false
	payload.Invoice.SpecialNotes = "/"

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal invoice request: %w", err)
	}

	// Sign the request with JWS
	signedBody, err := c.signRequest(body)
	if err != nil {
		return nil, fmt.Errorf("sign request: %w", err)
	}

	// Send to FURS
	resp, err := c.postRequest(ctx, "/v1/cash_registers/invoices", signedBody)
	if err != nil {
		return nil, fmt.Errorf("FURS request: %w", err)
	}

	// Parse the JWS response
	fursResponse, err := c.parseJWSResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("parse FURS response: %w", err)
	}

	// Increment invoice number on success
	c.settings.InvoiceNumber++

	return fursResponse, nil
}

// Echo tests the connection to FURS.
func (c *FURSClient) Echo(ctx context.Context) error {
	resp, err := c.postRequest(ctx, "/v1/cash_registers/echo", []byte(`{}`))
	if err != nil {
		return fmt.Errorf("echo request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("echo failed with status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

// signRequest creates a JWS compact serialization with RS256.
func (c *FURSClient) signRequest(payload []byte) ([]byte, error) {
	header := map[string]interface{}{
		"alg": "RS256",
		"x5c": []string{base64.StdEncoding.EncodeToString(c.cert.Raw)},
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return nil, err
	}

	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadB64 := base64.RawURLEncoding.EncodeToString(payload)

	signingInput := []byte(headerB64 + "." + payloadB64)
	hash := sha256.Sum256(signingInput)
	signature, err := rsa.SignPKCS1v15(rand.Reader, c.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, err
	}

	sigB64 := base64.RawURLEncoding.EncodeToString(signature)

	return []byte(headerB64 + "." + payloadB64 + "." + sigB64), nil
}

// postRequest sends an HTTP POST to the FURS API.
func (c *FURSClient) postRequest(ctx context.Context, path string, body []byte) (*http.Response, error) {
	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/jose")
	req.Header.Set("Accept", "application/jose")

	return c.httpClient.Do(req)
}

// parseJWSResponse extracts the InvoiceResponse from a JWS response body.
func (c *FURSClient) parseJWSResponse(resp *http.Response) (*models.InvoiceResponse, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FURS returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse JWS compact format: header.payload.signature
	parts := splitJWS(string(body))
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWS response format")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("decode JWS payload: %w", err)
	}

	var invoiceResp models.InvoiceResponse
	if err := json.Unmarshal(payloadBytes, &invoiceResp); err != nil {
		return nil, fmt.Errorf("unmarshal invoice response: %w", err)
	}

	return &invoiceResp, nil
}

func splitJWS(s string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}

func buildTaxesPerSeller(items []models.InvoiceItem) []struct {
	VAT []struct {
		TaxRate       float64 `json:"TaxRate"`
		TaxableAmount float64 `json:"TaxableAmount"`
		TaxAmount     float64 `json:"TaxAmount"`
	} `json:"VAT"`
} {
	type taxGroup struct {
		taxableAmount float64
		taxAmount     float64
	}

	groups := make(map[float64]*taxGroup)
	order := make([]float64, 0)

	for _, item := range items {
		if _, exists := groups[item.TaxRate]; !exists {
			groups[item.TaxRate] = &taxGroup{}
			order = append(order, item.TaxRate)
		}
		groups[item.TaxRate].taxableAmount += item.TaxableAmount
		groups[item.TaxRate].taxAmount += item.TaxAmount
	}

	vatArray := make([]struct {
		TaxRate       float64 `json:"TaxRate"`
		TaxableAmount float64 `json:"TaxableAmount"`
		TaxAmount     float64 `json:"TaxAmount"`
	}, 0, len(groups))

	for _, rate := range order {
		g := groups[rate]
		vatArray = append(vatArray, struct {
			TaxRate       float64 `json:"TaxRate"`
			TaxableAmount float64 `json:"TaxableAmount"`
			TaxAmount     float64 `json:"TaxAmount"`
		}{
			TaxRate:       rate,
			TaxableAmount: roundTo2(g.taxableAmount),
			TaxAmount:     roundTo2(g.taxAmount),
		})
	}

	return []struct {
		VAT []struct {
			TaxRate       float64 `json:"TaxRate"`
			TaxableAmount float64 `json:"TaxableAmount"`
			TaxAmount     float64 `json:"TaxAmount"`
		} `json:"VAT"`
	}{{VAT: vatArray}}
}

func roundTo2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}

func readCertificateFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Try PEM first
	if block, _ := pem.Decode(data); block != nil {
		return data, nil
	}

	// Assume DER/PKCS#12
	return data, nil
}
