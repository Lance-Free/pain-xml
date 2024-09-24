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
	MessageID            string           `xml:"MsgId,omitempty"`    // MessageID is a unique identifier for the message.
	CreationDateTime     CreationDateTime `xml:"CreDtTm,omitempty"`  // CreationDateTime indicates when the message was created.
	NumberOfTransactions int              `xml:"NbOfTxs,omitempty"`  // NumberOfTransactions is the total number of transactions.
	ControlSum           decimal.Decimal  `xml:"CtrlSum,omitempty"`  // ControlSum is the total amount for the transaction batch.
	InitiatingParty      InitiatingParty  `xml:"InitgPty,omitempty"` // InitiatingParty provides details of the party initiating the transaction.
}

// InitiatingParty holds information about the party initiating the transaction.
type InitiatingParty struct {
	Name string `xml:"Nm,omitempty"` // Name is the name of the initiating party.
}

// Account represents the account details of the debtor or creditor.
type Account struct {
	ID       AccountID `xml:"Id,omitempty"`  // ID contains account identification (e.g., IBAN).
	Currency string    `xml:"Ccy,omitempty"` // Currency represents the account's currency (e.g., EUR).
}

// AccountID contains the IBAN information.
type AccountID struct {
	IBAN string `xml:"IBAN,omitempty"` // IBAN represents the International Bank Account Number.
}

// Party represents the party that sends or receives the payment.
type Party struct {
	Name          string        `xml:"Nm,omitempty"`      // Name is the name of the debtor.
	PostalAddress PostalAddress `xml:"PstlAdr,omitempty"` // PostalAddress contains the address of the debtor.
}

// PostalAddress represents the postal address of the debtor or creditor.
type PostalAddress struct {
	Country     string   `xml:"Ctry,omitempty"`    // Country is the country code (e.g., NL for Netherlands).
	AddressLine []string `xml:"AdrLine,omitempty"` // AddressLine represents the lines of the address.
}

// PaymentTypeInfo contains information about the priority and service level.
type PaymentTypeInfo struct {
	InstructionPriority string       `xml:"InstrPrty,omitempty"` // InstructionPriority indicates the priority of the instruction (e.g., NORM).
	ServiceLevel        ServiceLevel `xml:"SvcLvl,omitempty"`    // ServiceLevel defines the service level (e.g., SEPA).
}

// ServiceLevel represents the service level code (e.g., SEPA).
type ServiceLevel struct {
	Code string `xml:"Cd,omitempty"` // Code is the code for the service level.
}

// PaymentID represents the identifiers for a payment.
type PaymentID struct {
	InstructionID string `xml:"InstrId,omitempty"`    // InstructionID is the instruction identifier.
	EndToEndID    string `xml:"EndToEndId,omitempty"` // EndToEndID is the unique end-to-end transaction identifier.
}

// InstigatedAmount represents the amount of an instigated order, including the currency and the textual representation.
type InstigatedAmount struct {
	Currency string          `xml:"Ccy,attr,omitempty"`
	Text     decimal.Decimal `xml:",chardata"`
}
