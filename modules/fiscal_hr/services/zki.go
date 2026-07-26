package services

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// CalculateZKI computes the Zaštitni Kod Izdavatelja (ZKI) for Croatian fiscalization.
//
// ZKI = MD5(RSA_SHA1_SIGN(OIB + DatVrijeme + BrOznRac + OznPosPr + OznNapUr + IznosUkupno))
func CalculateZKI(privKey *rsa.PrivateKey, oib, datVrijeme, brOznRac, oznPosPr, oznNapUr, iznosUkupno string) (string, error) {
	data := oib + datVrijeme + brOznRac + oznPosPr + oznNapUr + iznosUkupno

	hashed := sha1.Sum([]byte(data))

	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA1, hashed[:])
	if err != nil {
		return "", fmt.Errorf("ZKI sign: %w", err)
	}

	md5Hash := md5.Sum(signature)
	return hex.EncodeToString(md5Hash[:]), nil
}

// FormatDateTimeHR formats time for Croatian fiscalization: dd.mm.yyyyThh:mm:ss
func FormatDateTimeHR(t time.Time) string {
	return t.UTC().Format("02.01.2006T15:04:05")
}

// FormatDateHR formats time for Croatian fiscalization: dd.mm.yyyy
func FormatDateHR(t time.Time) string {
	return t.UTC().Format("02.01.2006")
}

// FormatAmountHR formats a float64 amount with exactly 2 decimal places.
func FormatAmountHR(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}

// ValidateOIB validates a Croatian OIB (11 digits with checksum).
func ValidateOIB(oib string) bool {
	if len(oib) != 11 {
		return false
	}

	for _, ch := range oib {
		if ch < '0' || ch > '9' {
			return false
		}
	}

	// OIB checksum algorithm
	sum := 0
	for i := 0; i < 10; i++ {
		digit := int(oib[i] - '0')
		if i%2 == 0 {
			sum += digit
		} else {
			sum += digit + 10
		}
	}
	checkDigit := 11 - (sum % 11)
	if checkDigit == 11 {
		checkDigit = 0
	}
	if checkDigit == 10 {
		checkDigit = 11
	}

	return checkDigit == int(oib[10]-'0')
}

// PaymentMethodCode converts common payment method strings to Croatian codes.
func PaymentMethodCode(method string) string {
	switch strings.ToLower(method) {
	case "cash", "gotovina":
		return "G"
	case "card", "kartica":
		return "K"
	case "transfer", "transakcija":
		return "T"
	default:
		return "O"
	}
}
