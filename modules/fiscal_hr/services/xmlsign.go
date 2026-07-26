package services

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	dsigNS    = "http://www.w3.org/2000/09/xmldsig#"
	excC14NNS = "http://www.w3.org/2001/10/xml-exc-c14n#"
	c14NNS    = "http://www.w3.org/TR/2001/REC-xml-c14n-20010315"
)

// SignEnvelope signs the RacunZahtjev XML inside a SOAP envelope using enveloped XML-DSig.
// Returns the complete signed SOAP envelope ready to send to CIS.
func SignEnvelope(soapXML string, privKey *rsa.PrivateKey, cert *x509.Certificate) (string, error) {
	// Extract the root Id from RacunZahtjev
	rootID := extractRootID(soapXML)
	if rootID == "" {
		return "", fmt.Errorf("RacunZahtjev Id attribute not found")
	}

	// 1. Compute DigestValue of RacunZahtjev (serialized canonical form)
	racunXML := extractRacunElement(soapXML)
	if racunXML == "" {
		return "", fmt.Errorf("cannot extract RacunZahtjev element")
	}

	digestBytes := sha1.Sum([]byte(racunXML))
	digestValue := base64.StdEncoding.EncodeToString(digestBytes[:])

	// 2. Build SignedInfo XML string
	signedInfoXML := buildSignedInfoString(rootID, digestValue)

	// 3. Sign canonicalized SignedInfo
	hashed := sha1.Sum([]byte(signedInfoXML))
	signatureValue, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA1, hashed[:])
	if err != nil {
		return "", fmt.Errorf("sign SignedInfo: %w", err)
	}
	sigValueB64 := base64.StdEncoding.EncodeToString(signatureValue)

	// 4. Build complete Signature element
	certB64 := base64.StdEncoding.EncodeToString(cert.Raw)
	signatureXML := buildSignatureString(rootID, digestValue, sigValueB64, certB64)

	// 5. Insert Signature as first child of RacunZahtjev in the SOAP envelope
	result := insertSignatureIntoRacun(soapXML, signatureXML)

	return result, nil
}

// extractRootID extracts the Id attribute from the RacunZahtjev element.
func extractRootID(soapXML string) string {
	idx := strings.Index(soapXML, "RacunZahtjev")
	if idx == -1 {
		return ""
	}

	// Find the opening tag that contains RacunZahtjev
	// Look backwards for '<' from the match point
	tagStart := strings.LastIndex(soapXML[:idx], "<")
	if tagStart == -1 {
		return ""
	}

	// Find the end of the opening tag
	tagEnd := strings.Index(soapXML[tagStart:], ">")
	if tagEnd == -1 {
		return ""
	}

	tag := soapXML[tagStart : tagStart+tagEnd]

	// Find Id="..."
	idStart := strings.Index(tag, `Id="`)
	if idStart == -1 {
		return ""
	}
	idStart += 4
	idEnd := strings.Index(tag[idStart:], `"`)
	if idEnd == -1 {
		return ""
	}

	return tag[idStart : idStart+idEnd]
}

// extractRacunElement extracts the full RacunZahtjev element from the SOAP envelope.
func extractRacunElement(soapXML string) string {
	// Find the opening tag by looking for <...RacunZahtjev
	idx := strings.Index(soapXML, "RacunZahtjev")
	if idx == -1 {
		return ""
	}

	// Go back to find the '<'
	tagStart := -1
	for i := idx - 1; i >= 0; i-- {
		if soapXML[i] == '<' {
			tagStart = i
			break
		}
	}
	if tagStart == -1 {
		return ""
	}

	// Find the matching '>' for the opening tag
	openEnd := strings.Index(soapXML[tagStart:], ">")
	if openEnd == -1 {
		return ""
	}

	// Find closing tag - try both prefixed and non-prefixed
	closingTag := "</tns:RacunZahtjev>"
	endIdx := strings.LastIndex(soapXML, closingTag)
	if endIdx == -1 {
		closingTag = "</RacunZahtjev>"
		endIdx = strings.LastIndex(soapXML, closingTag)
	}
	if endIdx == -1 {
		return ""
	}

	return soapXML[tagStart : endIdx+len(closingTag)]
}

// buildSignedInfoString builds the SignedInfo element as a canonicalized string.
func buildSignedInfoString(rootID, digestValue string) string {
	return fmt.Sprintf(`<ds:SignedInfo xmlns:ds="%s"><ds:CanonicalizationMethod Algorithm="%s"></ds:CanonicalizationMethod><ds:SignatureMethod Algorithm="%s"></ds:SignatureMethod><ds:Reference URI="#%s"><ds:Transforms><ds:Transform Algorithm="%s"></ds:Transform><ds:Transform Algorithm="%s"></ds:Transform></ds:Transforms><ds:DigestMethod Algorithm="%s"></ds:DigestMethod><ds:DigestValue>%s</ds:DigestValue></ds:Reference></ds:SignedInfo>`,
		dsigNS,
		c14NNS,
		dsigNS+"#rsa-sha1",
		rootID,
		dsigNS+"#enveloped-signature",
		excC14NNS,
		dsigNS+"#sha1",
		digestValue,
	)
}

// buildSignatureString builds the complete ds:Signature element.
func buildSignatureString(rootID, digestValue, sigValue, certB64 string) string {
	return fmt.Sprintf(`<ds:Signature xmlns:ds="%s" Id="sig0"><ds:SignedInfo><ds:CanonicalizationMethod Algorithm="%s"></ds:CanonicalizationMethod><ds:SignatureMethod Algorithm="%s"></ds:SignatureMethod><ds:Reference URI="#%s"><ds:Transforms><ds:Transform Algorithm="%s"></ds:Transform><ds:Transform Algorithm="%s"></ds:Transform></ds:Transforms><ds:DigestMethod Algorithm="%s"></ds:DigestMethod><ds:DigestValue>%s</ds:DigestValue></ds:Reference></ds:SignedInfo><ds:SignatureValue>%s</ds:SignatureValue><ds:KeyInfo><ds:X509Data><ds:X509Certificate>%s</ds:X509Certificate></ds:X509Data></ds:KeyInfo></ds:Signature>`,
		dsigNS,
		c14NNS,
		dsigNS+"#rsa-sha1",
		rootID,
		dsigNS+"#enveloped-signature",
		excC14NNS,
		dsigNS+"#sha1",
		digestValue,
		sigValue,
		certB64,
	)
}

// insertSignatureIntoRacun inserts the Signature element as the first child of RacunZahtjev.
func insertSignatureIntoRacun(soapXML, signatureXML string) string {
	// Find the end of the RacunZahtjev opening tag
	idx := strings.Index(soapXML, "RacunZahtjev")
	if idx == -1 {
		return soapXML
	}

	tagStart := strings.LastIndex(soapXML[:idx], "<")
	if tagStart == -1 {
		return soapXML
	}

	// Find the > that closes the opening tag
	tagEnd := strings.Index(soapXML[tagStart:], ">")
	if tagEnd == -1 {
		return soapXML
	}

	insertPos := tagStart + tagEnd + 1

	return soapXML[:insertPos] + "\n" + signatureXML + soapXML[insertPos:]
}

// ParseCertChain extracts the leaf certificate from raw cert bytes.
func ParseCertChain(rawCerts [][]byte) (*x509.Certificate, error) {
	if len(rawCerts) == 0 {
		return nil, fmt.Errorf("no certificates in chain")
	}
	return x509.ParseCertificate(rawCerts[0])
}
