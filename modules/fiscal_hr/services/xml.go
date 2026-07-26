package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nutrixpos/pos/modules/fiscal_hr/models"
)

const (
	fiscalNamespace = "http://www.apis-it.hr/fin/2012/types/f73"
	dsigNamespace   = "http://www.w3.org/2000/09/xmldsig#"
	soapNamespace   = "http://schemas.xmlsoap.org/soap/envelope/"
)

// BuildRacunXML creates the RacunZahtjev XML for Croatian fiscalization.
func BuildRacunXML(oib, datVrijeme, brOznRac, oznPosPr, oznNapUr string,
	pdv []PDVEntry, iznosUkupno float64, nacinPlac, oibOper, zastKod string) string {

	var pdvXML strings.Builder
	for _, p := range pdv {
		pdvXML.WriteString(fmt.Sprintf(`      <tns:Porez>
        <tns:Stopa>%.2f</tns:Stopa>
        <tns:Osnovica>%.2f</tns:Osnovica>
        <tns:Iznos>%.2f</tns:Iznos>
      </tns:Porez>
`, p.TaxRate, p.TaxableAmount, p.TaxAmount))
	}

	idPoruke := uuid.New().String()
	rootID := "uuid-" + idPoruke

	xml := fmt.Sprintf(`<tns:RacunZahtjev xmlns:tns="%s" xmlns:ds="%s" Id="%s">
  <tns:Zaglavlje>
    <tns:IdPoruke>%s</tns:IdPoruke>
    <tns:DatumVrijeme>%s</tns:DatumVrijeme>
  </tns:Zaglavlje>
  <tns:Racun>
    <tns:Oib>%s</tns:Oib>
    <tns:USustPdv>1</tns:USustPdv>
    <tns:DatVrijeme>%s</tns:DatVrijeme>
    <tns:OznSlijed>P</tns:OznSlijed>
    <tns:BrRac>
      <tns:BrOznRac>%s</tns:BrOznRac>
      <tns:OznPosPr>%s</tns:OznPosPr>
      <tns:OznNapUr>%s</tns:OznNapUr>
    </tns:BrRac>
    <tns:Pdv>
%s    </tns:Pdv>
    <tns:IznosOslobPdv>0.00</tns:IznosOslobPdv>
    <tns:IznosMarza>0.00</tns:IznosMarza>
    <tns:IznosNePodlOpor>0.00</tns:IznosNePodlOpor>
    <tns:IznosUkupno>%.2f</tns:IznosUkupno>
    <tns:NacinPlac>%s</tns:NacinPlac>
    <tns:OibOper>%s</tns:OibOper>
    <tns:ZastKod>%s</tns:ZastKod>
    <tns:NakDost>false</tns:NakDost>
  </tns:Racun>
</tns:RacunZahtjev>`,
		fiscalNamespace, dsigNamespace, rootID,
		idPoruke, datVrijeme,
		oib, datVrijeme,
		brOznRac, oznPosPr, oznNapUr,
		pdvXML.String(),
		iznosUkupno, nacinPlac, oibOper, zastKod,
	)

	return xml
}

// WrapInSOAP wraps the RacunZahtjev XML in a SOAP 1.1 envelope.
func WrapInSOAP(body string) string {
	return fmt.Sprintf(`<soapenv:Envelope xmlns:soapenv="%s" xmlns:tns="%s">
  <soapenv:Body>
    %s
  </soapenv:Body>
</soapenv:Envelope>`, soapNamespace, fiscalNamespace, body)
}

// BuildEchoXML creates a SOAP echo request.
func BuildEchoXML() string {
	return fmt.Sprintf(`<soapenv:Envelope xmlns:soapenv="%s" xmlns:tns="%s">
  <soapenv:Body>
    <tns:echo/>
  </soapenv:Body>
</soapenv:Envelope>`, soapNamespace, fiscalNamespace)
}

// PDVEntry represents a single VAT rate group.
type PDVEntry struct {
	TaxRate       float64
	TaxableAmount float64
	TaxAmount     float64
}

// BuildPDVEntries groups invoice items by tax rate.
func BuildPDVEntries(items []models.InvoiceItemHR) []PDVEntry {
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

	entries := make([]PDVEntry, 0, len(groups))
	for _, rate := range order {
		g := groups[rate]
		entries = append(entries, PDVEntry{
			TaxRate:       rate,
			TaxableAmount: roundTo2(g.taxableAmount),
			TaxAmount:     roundTo2(g.taxAmount),
		})
	}

	return entries
}

// FormatTimestampHR formats the current time for the Zaglavlje DatumVrijeme field.
func FormatTimestampHR(t time.Time) string {
	return FormatDateTimeHR(t)
}

func roundTo2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}
