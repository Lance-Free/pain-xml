package painxml

import (
	_ "embed"
	"encoding/xml"
	"github.com/lance-free/pain-xml/document"
	"github.com/lance-free/pain-xml/order"
	"testing"
	"time"
)

//go:embed sample.xml
var sampleXML []byte

func TestToOrder(t *testing.T) {
	var d document.DirectDebit
	err := xml.Unmarshal(sampleXML, &d)
	if err != nil {
		t.Fatalf("failed to unmarshal XML: %v", err)
	}

	_, err = ToOrder(d)
	if err != nil {
		t.Fatalf("failed to convert to order: %v", err)
	}

	if d.CustomerDirectDebitInitiation.PaymentInformation.PaymentInformationId != "Incasso SDD123" {
		t.Errorf("payment information id does not match: %v", d.CustomerDirectDebitInitiation.PaymentInformation.PaymentInformationId)
	}
}

func TestToDocument(t *testing.T) {
	o := order.Order{
		ExecutionDate: time.Now(),
		Transactions: []order.Transaction{
			{
				Name:       "John Doe",
				Street:     "Main Street 1",
				PostalCode: "12345",
				Place:      "Small Town",
				Country:    "DE",
				IBAN:       "DE12345678901234567890",
				BIC:        "GENODEF1M01",
				Currency:   "EUR",
				Amount:     100.00,
			},
		},
		Creditor: order.Creditor{
			Name:       "Jane Doe",
			Street:     "Main Street 2",
			PostalCode: "54321",
			Place:      "Big City",
			Country:    "DE",
			IBAN:       "DE09876543210987654321",
			BIC:        "GENODEF1M02",
		},
	}

	result, err := ToDocument(o)
	if err != nil {
		t.Errorf("failed to convert to document: %v", err)
	}

	if result.CustomerDirectDebitInitiation.GroupHeader.ControlSum != "100.00" {
		t.Errorf("control sum does not match: %v", result.CustomerDirectDebitInitiation.GroupHeader.ControlSum)
	}

	if len(result.CustomerDirectDebitInitiation.PaymentInformation.DirectDebitTransactionInformation) != 1 {
		t.Errorf("transaction count does not match: %v", len(result.CustomerDirectDebitInitiation.PaymentInformation.DirectDebitTransactionInformation))
	}

	if result.CustomerDirectDebitInitiation.PaymentInformation.Creditor.Name != "Jane Doe" {
		t.Errorf("creditor name does not match: %v", result.CustomerDirectDebitInitiation.PaymentInformation.Creditor.Name)
	}
}
