package document

import (
	"encoding/xml"
)

// InitiatingParty represents the party initiating the order.
type InitiatingParty struct {
	Name string `xml:"Nm"`
}

// GroupHeader represents the header information of a group in a Direct Debit document.
type GroupHeader struct {
	MessageID            string          `xml:"MsgId"`
	CreationDateTime     string          `xml:"CreDtTm"`
	NumberOfTransactions string          `xml:"NbOfTxs"`
	ControlSum           string          `xml:"CtrlSum"`
	InitiatingParty      InitiatingParty `xml:"InitgPty"`
}

// ServiceLevel represents the service level for a payment.
type ServiceLevel struct {
	Code string `xml:"Cd"`
}

// LocalInstrument represents a local instrument code used in the payment type information.
type LocalInstrument struct {
	Code string `xml:"Cd"`
}

// PaymentTypeInformation represents the type information for a payment.
type PaymentTypeInformation struct {
	ServiceLevel    ServiceLevel    `xml:"SvcLvl"`
	LocalInstrument LocalInstrument `xml:"LclInstrm"`
	SequenceType    string          `xml:"SeqTp"`
}

// Creditor represents a creditor in a payment order.
type Creditor struct {
	Name          string        `xml:"Nm"`
	PostalAddress PostalAddress `xml:"PstlAdr"`
}

// IBAN represents an International Bank Account Number.
type IBAN struct {
	IBAN string `xml:"IBAN"`
}

// CreditorAccount represents the creditor's bank account information.
type CreditorAccount struct {
	ID IBAN `xml:"Id"`
}

// CreditorAgent represents the creditor agent information in a direct debit order.
type CreditorAgent struct {
	FinancialInstitutionIdentification FinancialInstitutionIdentification `xml:"FinInstnId"`
}

type SchemeName struct {
	Proprietary string `xml:"Prtry"`
}

type Other struct {
	ID         string     `xml:"Id"`
	SchemeName SchemeName `xml:"SchmeNm"`
}

// PrivateIdentification is a type that represents private identification information for a specific scheme.
type PrivateIdentification struct {
	Other Other `xml:"Othr"`
}

// ID represents an identification element in XML, containing the PrivateIdentification element.
type ID struct {
	PrivateIdentification PrivateIdentification `xml:"PrvtId"`
}

// CreditorSchemeIdentification represents the identification of the creditor's scheme.
type CreditorSchemeIdentification struct {
	ID ID `xml:"Id"`
}

// InstigatedAmount represents the amount of an instigated order, including the currency and the textual representation.
type InstigatedAmount struct {
	Currency string `xml:"Ccy,attr"`
	Text     string `xml:",chardata"`
}

// PaymentID represents a payment identifier.
type PaymentID struct {
	EndToEndID string `xml:"EndToEndId"`
}

// MandateRelatedInformation is a type that represents information related to a mandate.
type MandateRelatedInformation struct {
	MandateID       string `xml:"MndtId"`
	DateOfSignature string `xml:"DtOfSgntr"`
}

// DirectDebitTransaction represents a direct debit order in an XML document.
type DirectDebitTransaction struct {
	MandateRelatedInformation MandateRelatedInformation `xml:"MndtRltdInf"`
}

// FinancialInstitutionIdentification represents the identification of a financial institution.
type FinancialInstitutionIdentification struct {
	BICFI string `xml:"BICFI"`
}

// DebtorAgent represents the debtor's financial institution identification.
type DebtorAgent struct {
	FinancialInstitutionIdentification FinancialInstitutionIdentification `xml:"FinInstnId"`
}

// PostalAddress represents a postal address.
type PostalAddress struct {
	TownName   string `xml:"TwnNm"`
	Country    string `xml:"Ctry"`
	StreetName string `xml:"StrtNm"`
	PostalCode string `xml:"PstCd"`
}

// Debtor represents a debtor in a Direct Debit order.
type Debtor struct {
	Name          string        `xml:"Nm"`
	PostalAddress PostalAddress `xml:"PstlAdr"`
}

// DebtorAccount represents a debtor's account information, including the IBAN.
type DebtorAccount struct {
	ID IBAN `xml:"Id"`
}

// RemittanceInformation represents the remittance information for a payment order.
type RemittanceInformation struct {
	Unstructured string `xml:"Ustrd"`
}

// DirectDebitTransactionInformation represents information about a direct debit order.
type DirectDebitTransactionInformation struct {
	PaymentID              PaymentID              `xml:"PmtId"`
	InstigatedAmount       InstigatedAmount       `xml:"InstdAmt"`
	DirectDebitTransaction DirectDebitTransaction `xml:"DrctDbtTx"`
	DebtorAgent            DebtorAgent            `xml:"DbtrAgt"`
	Debtor                 Debtor                 `xml:"Dbtr"`
	DebtorAccount          DebtorAccount          `xml:"DbtrAcct"`
	RemittanceInformation  RemittanceInformation  `xml:"RmtInf"`
}

// PaymentInformation represents the payment information for a direct debit order.
type PaymentInformation struct {
	PaymentInformationId              string                              `xml:"PmtInfId"`
	PaymentMethod                     string                              `xml:"PmtMtd"`
	NumberOfTransactions              string                              `xml:"NbOfTxs"`
	ControlSum                        string                              `xml:"CtrlSum"`
	PaymentTypeInformation            PaymentTypeInformation              `xml:"PmtTpInf"`
	RequestedCollectionDate           string                              `xml:"ReqdColltnDt"`
	Creditor                          Creditor                            `xml:"Cdtr"`
	CreditorAccount                   CreditorAccount                     `xml:"CdtrAcct"`
	CreditorAgent                     CreditorAgent                       `xml:"CdtrAgt"`
	CreditorSchemeIdentification      CreditorSchemeIdentification        `xml:"CdtrSchmeId"`
	DirectDebitTransactionInformation []DirectDebitTransactionInformation `xml:"DrctDbtTxInf"`
}

type CustomerDirectDebitInitiation struct {
	GroupHeader        GroupHeader        `xml:"GrpHdr"`
	PaymentInformation PaymentInformation `xml:"PmtInf"`
}

// DirectDebit represents a direct debit document that is used in the context of payment initiation.
type DirectDebit struct {
	XMLName                       xml.Name                      `xml:"Document"`
	Xmlns                         string                        `xml:"xmlns,attr"`
	CustomerDirectDebitInitiation CustomerDirectDebitInitiation `xml:"CstmrDrctDbtInitn"`
}
