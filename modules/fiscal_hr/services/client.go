package services

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/nutrixpos/pos/modules/fiscal_hr/models"
	"golang.org/x/crypto/pkcs12"
)

var (
	testBaseURLHR      = "https://cistest.apis-it.hr:8449/FiskalizacijaServiceTest"
	productionBaseURLHR = "https://cis.porezna-uprava.hr:8449/FiskalizacijaService"
)

// CISClient handles communication with the Croatian CIS fiscalization service.
type CISClient struct {
	settings   *models.FiscalSettingsHR
	privateKey *rsa.PrivateKey
	cert       *x509.Certificate
	httpClient *http.Client
	baseURL    string
}

// NewCISClient creates a new Croatian CIS client from fiscal settings.
func NewCISClient(settings *models.FiscalSettingsHR) (*CISClient, error) {
	if settings.CertificatePath == "" {
		return nil, fmt.Errorf("certificate path is required")
	}

	p12Data, err := readP12File(settings.CertificatePath)
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

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: 30 * time.Second,
	}

	baseURL := testBaseURLHR
	if settings.Environment == "production" {
		baseURL = productionBaseURLHR
	}

	return &CISClient{
		settings:   settings,
		privateKey: rsaKey,
		cert:       certificate,
		httpClient: httpClient,
		baseURL:    baseURL,
	}, nil
}

// NewCISClientFromKey creates a client from a pre-loaded RSA key and certificate (for testing).
func NewCISClientFromKey(settings *models.FiscalSettingsHR, privateKey *rsa.PrivateKey, cert *x509.Certificate) *CISClient {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12},
		},
		Timeout: 30 * time.Second,
	}

	baseURL := testBaseURLHR
	if settings.Environment == "production" {
		baseURL = productionBaseURLHR
	}

	return &CISClient{
		settings:   settings,
		privateKey: privateKey,
		cert:       cert,
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

// FiscalizeInvoice sends an invoice to CIS and returns the JIR.
func (c *CISClient) FiscalizeInvoice(ctx context.Context, req *models.InvoiceRequestHR) (*models.InvoiceResponseHR, error) {
	oib := c.settings.OIB
	datVrijeme := FormatDateTimeHR(time.Now())
	brOznRac := fmt.Sprintf("%d", c.settings.InvoiceNumber)
	oznPosPr := c.settings.BusinessPremiseID
	oznNapUr := c.settings.ElectronicDeviceID
	nacinPlac := PaymentMethodCode(req.PaymentMethod)

	operatorOIB := c.settings.OperatorOIB
	if req.OperatorOIB != "" {
		operatorOIB = req.OperatorOIB
	}

	pdvEntries := BuildPDVEntries(req.Items)

	zki, err := CalculateZKI(
		c.privateKey,
		oib,
		datVrijeme,
		brOznRac,
		oznPosPr,
		oznNapUr,
		FormatAmountHR(req.TotalAmount),
	)
	if err != nil {
		return nil, fmt.Errorf("calculate ZKI: %w", err)
	}

	bodyXML := BuildRacunXML(
		oib, datVrijeme, brOznRac, oznPosPr, oznNapUr,
		pdvEntries, req.TotalAmount, nacinPlac, operatorOIB, zki,
	)

	// Sign the XML body and wrap in SOAP envelope
	signedBody, err := SignEnvelope(
		WrapInSOAP(bodyXML),
		c.privateKey,
		c.cert,
	)
	if err != nil {
		return nil, fmt.Errorf("sign XML: %w", err)
	}

	soapEnvelope := signedBody

	resp, err := c.postRequest(ctx, soapEnvelope)
	if err != nil {
		return nil, fmt.Errorf("CIS request: %w", err)
	}

	jir, err := ParseRacunOdgovor(resp)
	if err != nil {
		return nil, fmt.Errorf("parse CIS response: %w", err)
	}

	c.settings.InvoiceNumber++

	return &models.InvoiceResponseHR{
		JIR:    jir,
		Status: "OK",
	}, nil
}

// PrivateKey returns the client's RSA private key for signing.
func (c *CISClient) PrivateKey() *rsa.PrivateKey { return c.privateKey }

// Certificate returns the client's X.509 certificate.
func (c *CISClient) Certificate() *x509.Certificate { return c.cert }

// SendRaw sends a pre-built SOAP envelope to CIS and returns the HTTP response.
func (c *CISClient) SendRaw(ctx context.Context, body string) (*http.Response, error) {
	return c.postRequest(ctx, body)
}

// Echo tests the connection to CIS.
func (c *CISClient) Echo(ctx context.Context) error {
	soapXML := BuildEchoXML()
	resp, err := c.postRequest(ctx, soapXML)
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

// postRequest sends an HTTP POST to the CIS SOAP service.
func (c *CISClient) postRequest(ctx context.Context, body string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewReader([]byte(body)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "")

	return c.httpClient.Do(req)
}

// RacunOdgovorXML represents the CIS XML response.
type RacunOdgovorXML struct {
	XMLName xml.Name `xml:"RacunOdgovor"`
	JIR     string   `xml:"Jir"`
	Greske  *struct {
		Greska []struct {
			SifraGreske   string `xml:"SifraGreske"`
			PorukaGreske  string `xml:"PorukaGreske"`
		} `xml:"Greska"`
	} `xml:"Greske"`
}

// ParseRacunOdgovor parses a CIS SOAP response and extracts the JIR.
func ParseRacunOdgovor(resp *http.Response) (string, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("CIS returned status %d: %s", resp.StatusCode, string(body))
	}

	// Try to extract JIR from SOAP body
	var soapResp struct {
		Body struct {
			RacunOdgovor RacunOdgovorXML `xml:"RacunOdgovor"`
		} `xml:"Body"`
	}

	if err := xml.Unmarshal(body, &soapResp); err != nil {
		return "", fmt.Errorf("unmarshal SOAP response: %w", err)
	}

	if soapResp.Body.RacunOdgovor.Greske != nil {
		for _, g := range soapResp.Body.RacunOdgovor.Greske.Greska {
			return "", fmt.Errorf("CIS error %s: %s", g.SifraGreske, g.PorukaGreske)
		}
	}

	return soapResp.Body.RacunOdgovor.JIR, nil
}

func readP12File(path string) ([]byte, error) {
	return os.ReadFile(path)
}
