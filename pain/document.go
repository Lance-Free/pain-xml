package pain

import (
	"encoding/xml"
	"github.com/shopspring/decimal"
	"time"
)

type CreationDateTime time.Time

func (c CreationDateTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(time.Time(c).Format("2006-01-02T15:04:05"), start)
}

func (c *CreationDateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02T15:04:05", v)
	if err != nil {
		return err
	}
	*c = CreationDateTime(t)
	return nil
}

// GroupHeader contains details about the group of transactions.
type GroupHeader struct {
	MessageID            string           `xml:"MsgId"`    // MessageID is a unique identifier for the message.
	CreationDateTime     CreationDateTime `xml:"CreDtTm"`  // CreationDateTime indicates when the message was created.
	NumberOfTransactions int              `xml:"NbOfTxs"`  // NumberOfTransactions is the total number of transactions.
	ControlSum           decimal.Decimal  `xml:"CtrlSum"`  // ControlSum is the total amount for the transaction batch.
	InitiatingParty      InitiatingParty  `xml:"InitgPty"` // InitiatingParty provides details of the party initiating the transaction.
}

// InitiatingParty holds information about the party initiating the transaction.
type InitiatingParty struct {
	Name string `xml:"Nm"` // Name is the name of the initiating party.
}

// Account represents the account details of the debtor or creditor.
type Account struct {
	ID       AccountID `xml:"Id"`  // ID contains account identification (e.g., IBAN).
	Currency string    `xml:"Ccy"` // Currency represents the account's currency (e.g., EUR).
}

// AccountID contains the IBAN information.
type AccountID struct {
	IBAN string `xml:"IBAN"` // IBAN represents the International Bank Account Number.
}

// Party represents the party that sends or receives the payment.
type Party struct {
	Name          string        `xml:"Nm"`      // Name is the name of the debtor.
	PostalAddress PostalAddress `xml:"PstlAdr"` // PostalAddress contains the address of the debtor.
}

// PostalAddress represents the postal address of the debtor or creditor.
type PostalAddress struct {
	Country     string   `xml:"Ctry"`    // Country is the country code (e.g., NL for Netherlands).
	AddressLine []string `xml:"AdrLine"` // AddressLine represents the lines of the address.
}

// PaymentTypeInfo contains information about the priority and service level.
type PaymentTypeInfo struct {
	InstructionPriority string       `xml:"InstrPrty"` // InstructionPriority indicates the priority of the instruction (e.g., NORM).
	ServiceLevel        ServiceLevel `xml:"SvcLvl"`    // ServiceLevel defines the service level (e.g., SEPA).
}

// ServiceLevel represents the service level code (e.g., SEPA).
type ServiceLevel struct {
	Code string `xml:"Cd"` // Code is the code for the service level.
}

// PaymentID represents the identifiers for a payment.
type PaymentID struct {
	InstructionID string `xml:"InstrId"`    // InstructionID is the instruction identifier.
	EndToEndID    string `xml:"EndToEndId"` // EndToEndID is the unique end-to-end transaction identifier.
}
