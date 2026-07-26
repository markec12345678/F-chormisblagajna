package services

import (
	"strings"
	"testing"
	"time"

	"github.com/nutrixpos/pos/modules/fiscal_hr/models"
)

func TestBuildRacunXML(t *testing.T) {
	items := []models.InvoiceItemHR{
		{Name: "Kava", Quantity: 2, UnitPrice: 2.50, TaxRate: 25, TaxableAmount: 5.00, TaxAmount: 1.25},
		{Name: "Sendvič", Quantity: 1, UnitPrice: 5.00, TaxRate: 25, TaxableAmount: 5.00, TaxAmount: 1.25},
	}

	pdv := BuildPDVEntries(items)

	xml := BuildRacunXML(
		"12345678901",
		"26.07.2026T15:00:00",
		"1",
		"1",
		"1",
		pdv,
		12.50,
		"G",
		"12345678901",
		"testzki123",
	)

	// Should contain required elements
	if !strings.Contains(xml, "<tns:Oib>12345678901</tns:Oib>") {
		t.Error("XML missing Oib")
	}
	if !strings.Contains(xml, "<tns:DatVrijeme>26.07.2026T15:00:00</tns:DatVrijeme>") {
		t.Error("XML missing DatVrijeme")
	}
	if !strings.Contains(xml, "<tns:BrOznRac>1</tns:BrOznRac>") {
		t.Error("XML missing BrOznRac")
	}
	if !strings.Contains(xml, "<tns:NacinPlac>G</tns:NacinPlac>") {
		t.Error("XML missing NacinPlac")
	}
	if !strings.Contains(xml, "<tns:ZastKod>testzki123</tns:ZastKod>") {
		t.Error("XML missing ZastKod")
	}
	if !strings.Contains(xml, "<tns:IznosUkupno>12.50</tns:IznosUkupno>") {
		t.Error("XML missing IznosUkupno")
	}
	if !strings.Contains(xml, "RacunZahtjev") {
		t.Error("XML missing RacunZahtjev root element")
	}
}

func TestWrapInSOAP(t *testing.T) {
	body := "<tns:echo/>"
	soap := WrapInSOAP(body)

	if !strings.Contains(soap, "soapenv:Envelope") {
		t.Error("SOAP missing Envelope")
	}
	if !strings.Contains(soap, "soapenv:Body") {
		t.Error("SOAP missing Body")
	}
	if !strings.Contains(soap, body) {
		t.Error("SOAP missing body content")
	}
}

func TestBuildEchoXML(t *testing.T) {
	xml := BuildEchoXML()

	if !strings.Contains(xml, "echo") {
		t.Error("Echo XML missing echo element")
	}
	if !strings.Contains(xml, "soapenv:Envelope") {
		t.Error("Echo XML missing SOAP envelope")
	}
}

func TestBuildPDVEntries_GroupsCorrectly(t *testing.T) {
	items := []models.InvoiceItemHR{
		{Name: "A", TaxRate: 25, TaxableAmount: 10, TaxAmount: 2.5},
		{Name: "B", TaxRate: 25, TaxableAmount: 20, TaxAmount: 5},
		{Name: "C", TaxRate: 13, TaxableAmount: 50, TaxAmount: 6.5},
	}

	entries := BuildPDVEntries(items)

	if len(entries) != 2 {
		t.Fatalf("expected 2 PDV entries, got %d", len(entries))
	}

	// First entry: 25% tax rate
	if entries[0].TaxRate != 25 {
		t.Errorf("entry 0 tax rate = %v, want 25", entries[0].TaxRate)
	}
	if entries[0].TaxableAmount != 30 {
		t.Errorf("entry 0 taxable = %v, want 30", entries[0].TaxableAmount)
	}
	if entries[0].TaxAmount != 7.5 {
		t.Errorf("entry 0 tax = %v, want 7.5", entries[0].TaxAmount)
	}

	// Second entry: 13% tax rate
	if entries[1].TaxRate != 13 {
		t.Errorf("entry 1 tax rate = %v, want 13", entries[1].TaxRate)
	}
}

func TestBuildPDVEntries_Empty(t *testing.T) {
	entries := BuildPDVEntries([]models.InvoiceItemHR{})
	if len(entries) != 0 {
		t.Errorf("expected 0 PDV entries, got %d", len(entries))
	}
}

func TestBuildPDVEntries_SingleRate(t *testing.T) {
	items := []models.InvoiceItemHR{
		{Name: "A", TaxRate: 25, TaxableAmount: 10, TaxAmount: 2.5},
		{Name: "B", TaxRate: 25, TaxableAmount: 20, TaxAmount: 5},
	}

	entries := BuildPDVEntries(items)
	if len(entries) != 1 {
		t.Fatalf("expected 1 PDV entry, got %d", len(entries))
	}
	if entries[0].TaxableAmount != 30 {
		t.Errorf("taxable = %v, want 30", entries[0].TaxableAmount)
	}
	if entries[0].TaxAmount != 7.5 {
		t.Errorf("tax = %v, want 7.5", entries[0].TaxAmount)
	}
}

func TestFormatTimestampHR(t *testing.T) {
	ts := time.Date(2026, 12, 25, 8, 15, 30, 0, time.UTC)
	got := FormatTimestampHR(ts)
	want := "25.12.2026T08:15:30"
	if got != want {
		t.Errorf("FormatTimestampHR = %q, want %q", got, want)
	}
}
