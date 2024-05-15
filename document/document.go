package document

import (
	"encoding/xml"
	"fmt"
	"github.com/lance-free/pain-xml/order"
	"time"
)

// InitiatingParty represents the party initiating the order.
type InitiatingParty struct {
	Name string `document:"Nm"`
}

// GroupHeader represents the header information of a group in a Direct Debit document.
type GroupHeader struct {
	MessageID            string `document:"MsgId"`
	CreationDateTime     string `document:"CreDtTm"`
	NumberOfTransactions string `document:"NbOfTxs"`
	ControlSum           string `document:"CtrlSum"`
	InitiatingParty      `document:"InitgPty"`
}

// ServiceLevel represents the service level for a payment.
type ServiceLevel struct {
	Code string `document:"Cd"`
}

// LocalInstrument represents a local instrument code used in the payment type information.
type LocalInstrument struct {
	Code string `document:"Cd"`
}

// PaymentTypeInformation represents the type information for a payment.
type PaymentTypeInformation struct {
	ServiceLevel    `document:"SvcLvl"`
	LocalInstrument `document:"LclInstrm"`
	SequenceType    string `document:"SeqTp"`
}

// Creditor represents a creditor in a payment order.
type Creditor struct {
	Name          string `document:"Nm"`
	PostalAddress `document:"PstlAdr"`
}

// IBAN represents an International Bank Account Number.
type IBAN struct {
	IBAN string `document:"IBAN"`
}

// CreditorAccount represents the creditor's bank account information.
type CreditorAccount struct {
	IBAN `document:"Id"`
}

// CreditorAgent represents the creditor agent information in a direct debit order.
type CreditorAgent struct {
	FinancialInstitutionIdentification `document:"FinInstnId"`
}

type SchemeName struct {
	Proprietary string `document:"Prtry"`
}

type Other struct {
	ID         string `document:"Id"`
	SchemeName `document:"SchmeNm"`
}

// PrivateIdentification is a type that represents private identification information for a specific scheme.
type PrivateIdentification struct {
	Other `document:"Othr"`
}

// ID represents an identification element in XML, containing the PrivateIdentification element.
type ID struct {
	PrivateIdentification `document:"PrvtId"`
}

// CreditorSchemeIdentification represents the identification of the creditor's scheme.
type CreditorSchemeIdentification struct {
	ID `document:"Id"`
}

// InstigatedAmount represents the amount of an instigated order, including the currency and the textual representation.
type InstigatedAmount struct {
	Currency string `document:"Ccy,attr"`
	Text     string `document:",chardata"`
}

// PaymentID represents a payment identifier.
type PaymentID struct {
	EndToEndID string `document:"EndToEndId"`
}

// MandateRelatedInformation is a type that represents information related to a mandate.
type MandateRelatedInformation struct {
	MandateID       string `document:"MndtId"`
	DateOfSignature string `document:"DtOfSgntr"`
}

// DirectDebitTransaction represents a direct debit order in an XML document.
type DirectDebitTransaction struct {
	MandateRelatedInformation `document:"MndtRltdInf"`
}

// FinancialInstitutionIdentification represents the identification of a financial institution.
type FinancialInstitutionIdentification struct {
	BICFI string `document:"BICFI"`
}

// DebtorAgent represents the debtor's financial institution identification.
type DebtorAgent struct {
	FinancialInstitutionIdentification `document:"FinInstnId"`
}

// PostalAddress represents a postal address.
type PostalAddress struct {
	TownName   string `document:"TwnNm"`
	Country    string `document:"Ctry"`
	StreetName string `document:"StrtNm"`
	PostalCode string `document:"PstCd"`
}

// Debtor represents a debtor in a Direct Debit order.
type Debtor struct {
	Name          string `document:"Nm"`
	PostalAddress `document:"PstlAdr"`
}

// DebtorAccount represents a debtor's account information, including the IBAN.
type DebtorAccount struct {
	IBAN `document:"Id"`
}

// RemittanceInformation represents the remittance information for a payment order.
type RemittanceInformation struct {
	Unstructured string `document:"Ustrd"`
}

// DirectDebitTransactionInformation represents information about a direct debit order.
type DirectDebitTransactionInformation struct {
	PaymentID              `document:"PmtId"`
	InstigatedAmount       `document:"InstdAmt"`
	DirectDebitTransaction `document:"DrctDbtTx"`
	DebtorAgent            `document:"DbtrAgt"`
	Debtor                 `document:"Dbtr"`
	DebtorAccount          `document:"DbtrAcct"`
	RemittanceInformation  `document:"RmtInf"`
}

// PaymentInformation represents the payment information for a direct debit order.
type PaymentInformation struct {
	PaymentInformationId              string `document:"PmtInfId"`
	PaymentMethod                     string `document:"PmtMtd"`
	NumberOfTransactions              string `document:"NbOfTxs"`
	ControlSum                        string `document:"CtrlSum"`
	PaymentTypeInformation            `document:"PmtTpInf"`
	RequestedCollectionDate           string `document:"ReqdColltnDt"`
	Creditor                          `document:"Cdtr"`
	CreditorAccount                   `document:"CdtrAcct"`
	CreditorAgent                     `document:"CdtrAgt"`
	CreditorSchemeIdentification      `document:"CdtrSchmeId"`
	DirectDebitTransactionInformation []DirectDebitTransactionInformation `document:"DrctDbtTxInf"`
}

type CustomerDirectDebitInitiation struct {
	GroupHeader        GroupHeader        `document:"GrpHdr"`
	PaymentInformation PaymentInformation `document:"PmtInf"`
}

// DirectDebit represents a direct debit document that is used in the context of payment initiation.
type DirectDebit struct {
	XMLName                       xml.Name                      `document:"Document"`
	Xmlns                         string                        `document:"xmlns,attr"`
	CustomerDirectDebitInitiation CustomerDirectDebitInitiation `document:"CstmrDrctDbtInitn"`
}

// ToOrder converts a DirectDebit to a order.Order.
// It extracts the necessary information from the DirectDebit
// and populates the Order struct with the relevant data.
func (document DirectDebit) ToOrder() (order.Order, error) {
	var transactions []order.Transaction
	for _, t := range document.CustomerDirectDebitInitiation.PaymentInformation.DirectDebitTransactionInformation {
		transactions = append(transactions, order.Transaction{
			Name:       t.PaymentID.EndToEndID,
			Street:     t.Debtor.PostalAddress.StreetName,
			PostalCode: t.Debtor.PostalAddress.PostalCode,
			Place:      t.Debtor.PostalAddress.TownName,
			Country:    t.Debtor.PostalAddress.Country,
			IBAN:       t.DebtorAccount.IBAN.IBAN,
			BIC:        t.DebtorAgent.FinancialInstitutionIdentification.BICFI,
			Currency:   t.InstigatedAmount.Currency,
			Amount:     0,
		})
	}

	executionDate, err := time.Parse("2006-01-02", document.CustomerDirectDebitInitiation.GroupHeader.CreationDateTime)
	if err != nil {
		return order.Order{}, fmt.Errorf("failed to parse execution date: %w", err)
	}
	return order.Order{
		ExecutionDate: executionDate,
		Transactions:  transactions,
		Creditor: order.Creditor{
			Name:       document.CustomerDirectDebitInitiation.PaymentInformation.Creditor.Name,
			Street:     document.CustomerDirectDebitInitiation.PaymentInformation.Creditor.PostalAddress.StreetName,
			PostalCode: document.CustomerDirectDebitInitiation.PaymentInformation.Creditor.PostalAddress.PostalCode,
			Place:      document.CustomerDirectDebitInitiation.PaymentInformation.Creditor.PostalAddress.TownName,
			Country:    document.CustomerDirectDebitInitiation.PaymentInformation.Creditor.PostalAddress.Country,
			IBAN:       document.CustomerDirectDebitInitiation.PaymentInformation.CreditorAccount.IBAN.IBAN,
			BIC:        document.CustomerDirectDebitInitiation.PaymentInformation.CreditorAgent.FinancialInstitutionIdentification.BICFI,
		},
	}, nil
}
