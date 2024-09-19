package ct

import (
	"encoding/xml"
	"github.com/lance-free/pain-xml/pain"
	"github.com/shopspring/decimal"
)

// Document represents the root of the XML structure with the required namespaces.
type Document struct {
	XMLName                          xml.Name                         `xml:"Document"`
	XMLNS                            string                           `xml:"xmlns,attr"`
	CustomerCreditTransferInitiation CustomerCreditTransferInitiation `xml:"CstmrCdtTrfInitn"`
}

// CustomerCreditTransferInitiation represents the root of the customer credit transfer initiation.
type CustomerCreditTransferInitiation struct {
	GroupHeader pain.GroupHeader `xml:"GrpHdr"`
	PaymentInfo []PaymentInfo    `xml:"PmtInf"`
}

// PaymentInfo represents payment-related information.
type PaymentInfo struct {
	PaymentInfoID             string                        `xml:"PmtInfId"`    // PaymentInfoID is a unique identifier for the payment info block.
	PaymentMethod             string                        `xml:"PmtMtd"`      // PaymentMethod defines the payment method (e.g., TRF for Transfer).
	BatchBooking              bool                          `xml:"BtchBookg"`   // BatchBooking defines whether batch booking is enabled.
	PaymentTypeInfo           pain.PaymentTypeInfo          `xml:"PmtTpInf"`    // PaymentTypeInfo provides additional payment type information.
	RequestedExecutionDate    string                        `xml:"ReqdExctnDt"` // RequestedExecutionDate is the date on which the payment should be executed.
	Debtor                    pain.Party                    `xml:"Dbtr"`        // Debtor contains information about the party sending the payment.
	DebtorAccount             pain.Account                  `xml:"DbtrAcct"`    // DebtorAccount represents the account of the debtor.
	DebtorAgent               FinancialInstitution          `xml:"DbtrAgt"`     // DebtorAgent represents the financial institution of the debtor.
	CreditTransferTransaction CreditTransferTransactionInfo `xml:"CdtTrfTxInf"` // CreditTransferTransaction contains details of the credit transfer.
}

// FinancialInstitution represents the financial institution.
type FinancialInstitution struct {
	FinancialInstitutionID FinancialInstitutionID `xml:"FinInstnId"` // FinancialInstitutionID holds the BIC of the institution.
}

// FinancialInstitutionID represents the BIC code of a financial institution.
type FinancialInstitutionID struct {
	BIC string `xml:"BIC"` // BIC is the Bank Identifier Code.
}

// CreditTransferTransactionInfo contains information about a single credit transfer transaction.
type CreditTransferTransactionInfo struct {
	PaymentID       pain.PaymentID        `xml:"PmtId"`              // PaymentID contains payment identifiers like instruction ID.
	PaymentTypeInfo pain.PaymentTypeInfo  `xml:"PmtTpInf,omitempty"` // PaymentTypeInfo is optional payment type information.
	Amount          Amount                `xml:"Amt"`                // Amount represents the instructed amount for the transfer.
	CreditorAgent   FinancialInstitution  `xml:"CdtrAgt"`            // CreditorAgent is the financial institution of the creditor.
	Creditor        pain.Party            `xml:"Cdtr"`               // Creditor contains information about the party receiving the payment.
	CreditorAccount pain.Account          `xml:"CdtrAcct"`           // CreditorAccount represents the account of the creditor.
	RemittanceInfo  RemittanceInformation `xml:"RmtInf"`             // RemittanceInfo contains the remittance details.
}

// Amount represents the amount of money in a transaction.
type Amount struct {
	InstructedAmount decimal.Decimal `xml:",chardata"` // InstructedAmount is the amount of money to be transferred.
	Currency         string          `xml:"Ccy,attr"`  // Currency defines the currency of the instructed amount.
}

// RemittanceInformation contains remittance details for the transaction.
type RemittanceInformation struct {
	Unstructured string `xml:"Ustrd"` // Unstructured contains free-form remittance information.
}

func NewDocument(customerCreditTransferInitiation CustomerCreditTransferInitiation) Document {
	return Document{
		XMLNS:                            "urn:iso:std:iso:20022:tech:xsd:pain.001.001.03",
		CustomerCreditTransferInitiation: customerCreditTransferInitiation,
	}
}
