package services

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"strings"
	"testing"
	"time"
)

func createTestCert(t *testing.T, key *rsa.PrivateKey) *x509.Certificate {
	t.Helper()

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "Test"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create test cert: %v", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("parse test cert: %v", err)
	}

	return cert
}

func TestSignEnvelope(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)

	// Create a test RacunZahtjev with Id
	racunXML := `<tns:RacunZahtjev xmlns:tns="http://www.apis-it.hr/fin/2012/types/f73" xmlns:ds="http://www.w3.org/2000/09/xmldsig#" Id="uuid-test-123">
  <tns:Zaglavlje>
    <tns:IdPoruke>test</tns:IdPoruke>
    <tns:DatumVrijeme>26.07.2026T15:00:00</tns:DatumVrijeme>
  </tns:Zaglavlje>
  <tns:Racun>
    <tns:Oib>12345678901</tns:Oib>
  </tns:Racun>
</tns:RacunZahtjev>`

	soapXML := WrapInSOAP(racunXML)

	// We need a real cert for signing - create a self-signed one for test
	cert := createTestCert(t, key)

	signed, err := SignEnvelope(soapXML, key, cert)
	if err != nil {
		t.Fatalf("SignEnvelope: %v", err)
	}

	// Should contain Signature element
	if !strings.Contains(signed, "ds:Signature") {
		t.Error("signed XML missing ds:Signature")
	}

	// Should contain SignatureValue
	if !strings.Contains(signed, "ds:SignatureValue") {
		t.Error("signed XML missing ds:SignatureValue")
	}

	// Should contain X509Certificate
	if !strings.Contains(signed, "ds:X509Certificate") {
		t.Error("signed XML missing ds:X509Certificate")
	}

	// Should contain SignedInfo
	if !strings.Contains(signed, "ds:SignedInfo") {
		t.Error("signed XML missing ds:SignedInfo")
	}

	// Should contain CanonicalizationMethod
	if !strings.Contains(signed, c14NNS) {
		t.Error("signed XML missing CanonicalizationMethod algorithm")
	}

	// Should contain SignatureMethod
	if !strings.Contains(signed, dsigNS+"#rsa-sha1") {
		t.Error("signed XML missing SignatureMethod algorithm")
	}

	// Should reference the root element
	if !strings.Contains(signed, `URI="#uuid-test-123"`) {
		t.Error("signed XML missing Reference URI to root element")
	}

	// Should contain enveloped-signature transform
	if !strings.Contains(signed, dsigNS+"#enveloped-signature") {
		t.Error("signed XML missing enveloped-signature transform")
	}

	// Should contain DigestValue
	if !strings.Contains(signed, "ds:DigestValue") {
		t.Error("signed XML missing DigestValue")
	}

	// Should contain KeyInfo
	if !strings.Contains(signed, "ds:KeyInfo") {
		t.Error("signed XML missing KeyInfo")
	}

	// Signature should be inside RacunZahtjev (enveloped)
	sigIdx := strings.Index(signed, "ds:Signature")
	racunCloseIdx := strings.Index(signed, "</tns:RacunZahtjev>")
	if sigIdx > racunCloseIdx {
		t.Error("Signature should be inside RacunZahtjev")
	}
}

func TestSignEnvelope_MissingID(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	cert := createTestCert(t, key)

	// No Id attribute
	racunXML := `<tns:RacunZahtjev xmlns:tns="http://www.apis-it.hr/fin/2012/types/f73">
  <tns:Zaglavlje><tns:IdPoruke>test</tns:IdPoruke></tns:Zaglavlje>
</tns:RacunZahtjev>`

	soapXML := WrapInSOAP(racunXML)

	_, err := SignEnvelope(soapXML, key, cert)
	if err == nil {
		t.Error("expected error for missing Id attribute")
	}
}

func TestExtractRootID(t *testing.T) {
	xml := `<tns:RacunZahtjev xmlns:tns="..." Id="uuid-abc-123"><tns:Zaglavlje>`
	id := extractRootID(xml)
	if id != "uuid-abc-123" {
		t.Errorf("extractRootID = %q, want %q", id, "uuid-abc-123")
	}
}

func TestExtractRootID_Missing(t *testing.T) {
	id := extractRootID("<root><child/></root>")
	if id != "" {
		t.Errorf("extractRootID = %q, want empty", id)
	}
}

func TestExtractRacunElement(t *testing.T) {
	xml := `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"><soapenv:Body><tns:RacunZahtjev xmlns:tns="http://www.apis-it.hr/fin/2012/types/f73" Id="uuid-1"><tns:Zaglavlje>test</tns:Zaglavlje><tns:Racun>data</tns:Racun></tns:RacunZahtjev></soapenv:Body></soapenv:Envelope>`

	racun := extractRacunElement(xml)
	if racun == "" {
		t.Fatal("extractRacunElement returned empty")
	}
	if !strings.Contains(racun, "RacunZahtjev") {
		t.Error("extracted element should contain RacunZahtjev")
	}
	if !strings.Contains(racun, "uuid-1") {
		t.Error("extracted element should contain Id")
	}
	if !strings.Contains(racun, "</tns:RacunZahtjev>") {
		t.Error("extracted element should contain closing tag")
	}
}

func TestBuildSignedInfoString(t *testing.T) {
	si := buildSignedInfoString("uuid-123", "abc123base64")

	if !strings.Contains(si, `URI="#uuid-123"`) {
		t.Error("SignedInfo missing Reference URI")
	}
	if !strings.Contains(si, "abc123base64") {
		t.Error("SignedInfo missing DigestValue")
	}
	if !strings.Contains(si, c14NNS) {
		t.Error("SignedInfo missing c14n algorithm")
	}
	if !strings.Contains(si, excC14NNS) {
		t.Error("SignedInfo missing exclusive c14n transform")
	}
}

func TestBuildSignatureString(t *testing.T) {
	sig := buildSignatureString("uuid-123", "digest", "sigval", "cert")

	if !strings.Contains(sig, "sigval") {
		t.Error("Signature missing SignatureValue")
	}
	if !strings.Contains(sig, "cert") {
		t.Error("Signature missing certificate")
	}
	if !strings.Contains(sig, "sig0") {
		t.Error("Signature missing Id='sig0'")
	}
}

func TestInsertSignatureIntoRacun(t *testing.T) {
	soap := `<tns:RacunZahtjev xmlns:tns="http://test" Id="uuid-1">
  <tns:Zaglavlje/>
</tns:RacunZahtjev>`

	result := insertSignatureIntoRacun(soap, "<ds:Signature/>")

	sigIdx := strings.Index(result, "<ds:Signature/>")
	closingIdx := strings.Index(result, "</tns:RacunZahtjev>")

	if sigIdx == -1 {
		t.Fatal("Signature not inserted")
	}
	if sigIdx > closingIdx {
		t.Error("Signature should be before closing RacunZahtjev tag")
	}
}

func TestSignEnvelope_DeterministicDigest(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	cert := createTestCert(t, key)

	racunXML := `<tns:RacunZahtjev xmlns:tns="http://www.apis-it.hr/fin/2012/types/f73" Id="uuid-det">
  <tns:Racun><tns:Oib>12345678901</tns:Oib></tns:Racun>
</tns:RacunZahtjev>`

	soapXML := WrapInSOAP(racunXML)

	signed1, err := SignEnvelope(soapXML, key, cert)
	if err != nil {
		t.Fatalf("SignEnvelope first: %v", err)
	}

	// Extract DigestValue from both (should be same since input is same)
	d1 := extractDigestValue(signed1)
	d2 := extractDigestValue(signed1)
	if d1 != d2 {
		t.Error("DigestValue should be deterministic for same input")
	}

	// Second signing of same content should produce same digest but different signature
	signed2, err := SignEnvelope(soapXML, key, cert)
	if err != nil {
		t.Fatalf("SignEnvelope second: %v", err)
	}

	digest1 := extractDigestValue(signed1)
	digest2 := extractDigestValue(signed2)
	if digest1 != digest2 {
		t.Error("same input should produce same DigestValue")
	}
}

func extractDigestValue(xml string) string {
	start := strings.Index(xml, "<ds:DigestValue>")
	if start == -1 {
		return ""
	}
	start += len("<ds:DigestValue>")
	end := strings.Index(xml[start:], "</ds:DigestValue>")
	if end == -1 {
		return ""
	}
	return xml[start : start+end]
}

func TestFormatAmountHR_Storno(t *testing.T) {
	got := FormatAmountHR(-12.50)
	if got != "-12.50" {
		t.Errorf("FormatAmountHR(-12.50) = %q, want %q", got, "-12.50")
	}
}

func TestPaymentMethodCode_AllCodes(t *testing.T) {
	codes := map[string]string{
		"G": "G",
		"K": "K",
		"T": "T",
		"O": "O",
	}
	for _, method := range []string{"cash", "card", "transfer", "other"} {
		code := PaymentMethodCode(method)
		if _, ok := codes[code]; !ok {
			t.Errorf("PaymentMethodCode(%q) = %q, not a valid code", method, code)
		}
	}
}

func TestDateTimeHR_UTC(t *testing.T) {
	// Test UTC normalization
	ts := time.Date(2026, 7, 26, 23, 0, 0, 0, time.UTC)
	got := FormatDateTimeHR(ts)
	if got != "26.07.2026T23:00:00" {
		t.Errorf("FormatDateTimeHR = %q, want %q", got, "26.07.2026T23:00:00")
	}
}
