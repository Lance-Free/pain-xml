package dd

import (
	"encoding/xml"
	"github.com/lance-free/pain-xml/pain"
	"github.com/shopspring/decimal"
)

// LocalInstrument represents a local instrument code used in the payment type information.
type LocalInstrument struct {
	Code string `xml:"Cd,omitempty"`
}

// Creditor represents a creditor in a payment order.
type Creditor struct {
	Name          string        `xml:"Nm,omitempty"`
	PostalAddress PostalAddress `xml:"PstlAdr,omitempty"`
}

// CreditorAgent represents the creditor agent information in a direct debit order.
type CreditorAgent struct {
	FinancialInstitutionIdentification FinancialInstitutionIdentification `xml:"FinInstnId,omitempty"`
}

type SchemeName struct {
	Proprietary string `xml:"Prtry,omitempty"`
}

type Other struct {
	ID         string     `xml:"Id,omitempty"`
	SchemeName SchemeName `xml:"SchmeNm,omitempty"`
}

// PrivateIdentification is a type that represents private identification information for a specific scheme.
type PrivateIdentification struct {
	Other Other `xml:"Othr,omitempty"`
}

// ID represents an identification element in XML, containing the PrivateIdentification element.
type ID struct {
	PrivateIdentification PrivateIdentification `xml:"PrvtId,omitempty"`
}

// CreditorSchemeIdentification represents the identification of the creditor's scheme.
type CreditorSchemeIdentification struct {
	ID ID `xml:"Id,omitempty"`
}

// MandateRelatedInformation is a type that represents information related to a mandate.
type MandateRelatedInformation struct {
	MandateID       string `xml:"MndtId,omitempty"`
	DateOfSignature string `xml:"DtOfSgntr,omitempty"`
}

// DirectDebitTransaction represents a direct debit order in an XML document.
type DirectDebitTransaction struct {
	MandateRelatedInformation MandateRelatedInformation `xml:"MndtRltdInf,omitempty"`
}

// FinancialInstitutionIdentification represents the identification of a financial institution.
type FinancialInstitutionIdentification struct {
	BICFI string `xml:"BICFI,omitempty"`
}

// DebtorAgent represents the debtor's financial institution identification.
type DebtorAgent struct {
	FinancialInstitutionIdentification FinancialInstitutionIdentification `xml:"FinInstnId,omitempty"`
}

// PostalAddress represents a postal address.
type PostalAddress struct {
	TownName   string `xml:"TwnNm,omitempty"`
	Country    string `xml:"Ctry,omitempty"`
	StreetName string `xml:"StrtNm,omitempty"`
	PostalCode string `xml:"PstCd,omitempty"`
}

// Debtor represents a debtor in a Direct Debit order.
type Debtor struct {
	Name          string        `xml:"Nm,omitempty"`
	PostalAddress PostalAddress `xml:"PstlAdr,omitempty"`
}

// RemittanceInformation represents the remittance information for a payment order.
type RemittanceInformation struct {
	Unstructured string `xml:"Ustrd,omitempty"`
}

// DirectDebitTransactionInformation represents information about a direct debit order.
type DirectDebitTransactionInformation struct {
	PaymentID              pain.PaymentID         `xml:"PmtId,omitempty"`
	InstigatedAmount       InstigatedAmount       `xml:"InstdAmt,omitempty"`
	DirectDebitTransaction DirectDebitTransaction `xml:"DrctDbtTx,omitempty"`
	DebtorAgent            DebtorAgent            `xml:"DbtrAgt,omitempty"`
	Debtor                 Debtor                 `xml:"Dbtr,omitempty"`
	DebtorAccount          pain.Account           `xml:"DbtrAcct,omitempty"`
	RemittanceInformation  RemittanceInformation  `xml:"RmtInf,omitempty"`
}

// PaymentInformation represents the payment information for a direct debit order.
type PaymentInformation struct {
	PaymentInformationId              string                              `xml:"PmtInfId,omitempty"`
	PaymentMethod                     string                              `xml:"PmtMtd,omitempty"`
	NumberOfTransactions              int                                 `xml:"NbOfTxs,omitempty"`
	ControlSum                        decimal.Decimal                     `xml:"CtrlSum,omitempty"`
	PaymentTypeInformation            pain.PaymentTypeInfo                `xml:"PmtTpInf,omitempty"`
	RequestedCollectionDate           string                              `xml:"ReqdColltnDt,omitempty"`
	Creditor                          Creditor                            `xml:"Cdtr,omitempty"`
	CreditorAccount                   pain.Account                        `xml:"CdtrAcct,omitempty"`
	CreditorAgent                     CreditorAgent                       `xml:"CdtrAgt,omitempty"`
	CreditorSchemeIdentification      CreditorSchemeIdentification        `xml:"CdtrSchmeId,omitempty"`
	DirectDebitTransactionInformation []DirectDebitTransactionInformation `xml:"DrctDbtTxInf,omitempty"`
}

type CustomerDirectDebitInitiation struct {
	GroupHeader        pain.GroupHeader   `xml:"GrpHdr,omitempty"`
	PaymentInformation PaymentInformation `xml:"PmtInf,omitempty"`
}

// DirectDebit represents a direct debit document that is used in the context of payment initiation.
type DirectDebit struct {
	XMLName                       xml.Name                      `xml:"Document,omitempty"`
	Xmlns                         string                        `xml:"xmlns,attr,omitempty"`
	CustomerDirectDebitInitiation CustomerDirectDebitInitiation `xml:"CstmrDrctDbtInitn,omitempty"`
}

func NewDocument(customerDirectDebitInitiation CustomerDirectDebitInitiation) DirectDebit {
	return DirectDebit{
		Xmlns:                         "urn:iso:std:iso:20022:tech:xsd:pain.008.001.02",
		CustomerDirectDebitInitiation: customerDirectDebitInitiation,
	}
}
