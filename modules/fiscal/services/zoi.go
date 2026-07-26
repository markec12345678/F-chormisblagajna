package services

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"
	"time"
)

// CalculateZOI computes the ZAŠČITNA OZNAKA IZDAJATELJA (ZOI) for a fiscal invoice.
//
// The ZOI is a 32-char hex string computed as:
//  1. Concatenate: taxNumber + issueDateTime(dd-MM-yyyy HH:mm:ss) + invoiceNumber + bpID + edID + invoiceAmount
//  2. RSA-SHA256 PKCS#1 v1.5 sign with private key
//  3. MD5 hash the signature → 32-char hex
func CalculateZOI(privKey *rsa.PrivateKey, taxNumber int, issueDateTime time.Time,
	invoiceNumber, bpID, edID string, amount float64) (string, error) {

	dateStr := issueDateTime.Format("02-01-2006 15:04:05")
	content := fmt.Sprintf("%d%s%s%s%s%s", taxNumber, dateStr, invoiceNumber, bpID, edID, formatAmount(amount))

	hashed := sha256.Sum256([]byte(content))

	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("ZOI sign: %w", err)
	}

	md5Hash := md5.Sum(signature)
	return fmt.Sprintf("%x", md5Hash), nil
}

// CalculateQRData builds the 60-digit printable string for QR/PDF417.
//
// Format: ZOI decimal (39 digits) + TaxNumber (8 digits) + DateTime YYMMDDHHmmss (12 digits) + Check digit (1 digit)
func CalculateQRData(zoiHex string, taxNumber int, issueDateTime time.Time) (string, error) {
	// Convert hex ZOI to decimal, zero-pad to 39 digits
	zoiBig := new(big.Int)
	_, ok := zoiBig.SetString(zoiHex, 16)
	if !ok {
		return "", fmt.Errorf("invalid ZOI hex: %s", zoiHex)
	}
	zoiDecimal := fmt.Sprintf("%039s", zoiBig.String())

	// Tax number as 8 digits
	taxStr := fmt.Sprintf("%08d", taxNumber)

	// DateTime as YYMMDDHHmmss (12 digits)
	dateStr := issueDateTime.Format("060102150405")

	// Concatenate without check digit
	data := zoiDecimal + taxStr + dateStr

	// Calculate check digit: sum all digits mod 10
	sum := 0
	for _, ch := range data {
		sum += int(ch - '0')
	}
	checkDigit := sum % 10

	return fmt.Sprintf("%s%d", data, checkDigit), nil
}

// FormatAmount formats a float64 amount with exactly 2 decimal places.
func formatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}

// FormatInvoiceNumber formats the invoice number as XXXXX-XXXXX-XXXXX.
func FormatInvoiceNumber(number int) string {
	s := fmt.Sprintf("%015d", number)
	return strings.Join([]string{s[0:5], s[5:10], s[10:15]}, "-")
}

// FormatIssueDateTime formats time for the FURS API payload (ISO 8601).
func FormatIssueDateTime(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05Z")
}
