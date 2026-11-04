package sct

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"

	"github.com/lance-free/pain-xml/internal/schematest"
	"github.com/shopspring/decimal"
)

func validMessage() Message {
	batchBooking := false

	return Message{
		ID:            "MSG-20261104-150405-0a1b2c3d",
		CreatedAt:     time.Date(2026, 11, 4, 15, 4, 5, 0, time.UTC),
		Debtor:        Account{Name: "PayToday B.V.", IBAN: "nl91 abna 0417 1643 00", BIC: "abnanl2a"},
		ExecutionDate: time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC),
		BatchBooking:  &batchBooking,
		Priority:      PriorityNormal,
		Transfers: []Transfer{
			{
				EndToEndID: "5000056789012345",
				Amount:     decimal.RequireFromString("1210.50"),
				Creditor:   Account{Name: "Klant B.V.", IBAN: "DE89 3704 0044 0532 0130 00"},
				UltimateDebtor: &Party{
					Name:           "Debiteur B.V.",
					Address:        &Address{Street: "Damrak", BuildingNumber: "1", PostCode: "1012 LG", Town: "Amsterdam", Country: "NL"},
					Identification: CompanyRegistration("12345678", "KVK"),
				},
				Reference:   "5000 0567 8901 2345",
				Description: "Doorbetaling factuur 2026-0042",
			},
			{
				Amount:      decimal.RequireFromString("99.99"),
				Creditor:    Account{Name: "Andere Klant B.V.", IBAN: "NL91ABNA0417164300", BIC: "DEUTDEFF500", Identification: OrganisationBIC("DEUTDEFF")},
				Description: "Doorbetaling factuur 2026-0043",
			},
		},
	}
}

func TestDocumentMatchesSchema(t *testing.T) {
	document, err := validMessage().Document()
	if err != nil {
		t.Fatalf("Document() error = %v", err)
	}

	schematest.Validate(t, document)
}

func TestDocumentSetsSEPAValues(t *testing.T) {
	document, err := validMessage().Document()
	if err != nil {
		t.Fatalf("Document() error = %v", err)
	}

	encoded, err := xml.Marshal(document)
	if err != nil {
		t.Fatalf("marshalling document: %v", err)
	}

	for _, want := range []string{
		"<NbOfTxs>2</NbOfTxs><CtrlSum>1310.49</CtrlSum><InitgPty><Nm>PayToday B.V.</Nm></InitgPty>",
		"<PmtInfId>MSG-20261104-150405-0a1b2c3d</PmtInfId><PmtMtd>TRF</PmtMtd><BtchBookg>false</BtchBookg>",
		"<PmtTpInf><InstrPrty>NORM</InstrPrty><SvcLvl><Cd>SEPA</Cd></SvcLvl></PmtTpInf><ReqdExctnDt>2026-11-05</ReqdExctnDt>",
		"<DbtrAcct><Id><IBAN>NL91ABNA0417164300</IBAN></Id></DbtrAcct><DbtrAgt><FinInstnId><BIC>ABNANL2A</BIC></FinInstnId></DbtrAgt><ChrgBr>SLEV</ChrgBr>",
		`<InstdAmt Ccy="EUR">1210.5</InstdAmt>`,
		"<Id><OrgId><Othr><Id>12345678</Id><SchmeNm><Cd>COID</Cd></SchmeNm><Issr>KVK</Issr></Othr></OrgId></Id>",
		"<IBAN>DE89370400440532013000</IBAN>",
		"<Ref>5000056789012345</Ref>",
		"<EndToEndId>NOTPROVIDED</EndToEndId>",
		"<CdtrAgt><FinInstnId><BIC>DEUTDEFF500</BIC></FinInstnId></CdtrAgt>",
		"<Cdtr><Nm>Andere Klant B.V.</Nm><Id><OrgId><BICOrBEI>DEUTDEFF</BICOrBEI></OrgId></Id></Cdtr>",
		"<Ustrd>Doorbetaling factuur 2026-0043</Ustrd>",
	} {
		if !strings.Contains(string(encoded), want) {
			t.Errorf("document does not contain %s", want)
		}
	}

	if strings.Contains(string(encoded), "Doorbetaling factuur 2026-0042") {
		t.Error("document sends the description next to the structured reference")
	}
}

func TestDocumentUsesDefaults(t *testing.T) {
	message := validMessage()
	message.CreatedAt = time.Time{}
	message.ExecutionDate = time.Time{}
	message.InitiatingParty = ""

	before := time.Now()
	document, err := message.Document()
	if err != nil {
		t.Fatalf("Document() error = %v", err)
	}

	header := document.CustomerCreditTransferInitiation.GroupHeader
	if time.Time(header.CreationDateTime).Before(before) {
		t.Errorf("CreationDateTime = %v; want the current time", time.Time(header.CreationDateTime))
	}

	if header.InitiatingParty.Name != message.Debtor.Name {
		t.Errorf("InitiatingParty = %q; want the debtor name %q", header.InitiatingParty.Name, message.Debtor.Name)
	}

	wantDate := time.Time(header.CreationDateTime).Format(executionDateLayout)
	if got := document.CustomerCreditTransferInitiation.PaymentInfo[0].RequestedExecutionDate; got != wantDate {
		t.Errorf("RequestedExecutionDate = %q; want %q", got, wantDate)
	}
}

func TestDocumentReturnsViolations(t *testing.T) {
	tests := []struct {
		name   string
		change func(*Message)
		want   string
	}{
		{name: "no transfers", change: func(m *Message) { m.Transfers = nil }, want: "at least one transaction is required"},
		{name: "no message id", change: func(m *Message) { m.ID = "" }, want: "GrpHdr.MsgId: is required"},
		{name: "invalid creditor IBAN", change: func(m *Message) { m.Transfers[0].Creditor.IBAN = "NL00ABNA0417164300" }, want: "CdtrAcct.Id.IBAN"},
		{name: "zero amount", change: func(m *Message) { m.Transfers[1].Amount = decimal.Zero }, want: "must be between 0.01 and 999999999.99"},
		{name: "description too long", change: func(m *Message) { m.Transfers[1].Description = strings.Repeat("A", 141) }, want: "RmtInf.Ustrd: is 141 characters"},
		{name: "ultimate debtor on two levels", change: func(m *Message) { m.UltimateDebtor = &Party{Name: "Debiteur B.V."} }, want: "set it on one level only"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			message := validMessage()
			test.change(&message)

			_, err := message.Document()
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Document() error = %v; want an error that contains %q", err, test.want)
			}
		})
	}
}
