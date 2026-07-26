package services

import (
	"fmt"
	"strings"
)

// QRCodeData represents the data for a QR code to be rendered.
type QRCodeData struct {
	ZOI           string
	QRData        string
	EOR           string
	InvoiceNumber string
}

// FormatReceiptText creates a human-readable receipt text.
func FormatReceiptText(data *QRCodeData, items []ReceiptItem, total float64) string {
	var sb strings.Builder

	sb.WriteString("================================\n")
	sb.WriteString("         NUTRIX POS\n")
	sb.WriteString("================================\n\n")

	for _, item := range items {
		sb.WriteString(fmt.Sprintf("%-20s %6.2f\n", item.Name, item.Price))
		if item.Quantity > 1 {
			sb.WriteString(fmt.Sprintf("  %.0f x %.2f\n", item.Quantity, item.UnitPrice))
		}
	}

	sb.WriteString("--------------------------------\n")
	sb.WriteString(fmt.Sprintf("%-20s %6.2f\n", "TOTAL", total))
	sb.WriteString("================================\n\n")

	sb.WriteString(fmt.Sprintf("Invoice: %s\n", data.InvoiceNumber))
	sb.WriteString(fmt.Sprintf("EOR: %s\n", data.EOR))
	sb.WriteString(fmt.Sprintf("ZOI: %s\n", data.ZOI))
	sb.WriteString("\n")
	sb.WriteString("Fiskalizirano po FURS\n")
	sb.WriteString("================================\n")

	return sb.String()
}

// ReceiptItem represents a single line on a receipt.
type ReceiptItem struct {
	Name     string
	Quantity float64
	UnitPrice float64
	Price    float64
}

// QRDataToBase64Image converts QR data string to a placeholder.
// In production, use github.com/skip2/go-qrcode for real QR codes.
func QRDataToBase64Image(qrData string) (string, error) {
	if qrData == "" {
		return "", fmt.Errorf("empty QR data")
	}
	// Placeholder — real implementation generates QR PNG
	return "", fmt.Errorf("QR code generation requires go-qrcode library (install: go get github.com/skip2/go-qrcode)")
}
