package ct

import (
	"encoding/xml"
	"github.com/lance-free/pain-xml/pain"
	"github.com/shopspring/decimal"
)

// Document represents the root of the XML structure with the required namespaces.
type Document struct {
	XMLName                          xml.Name                         `xml:"Document,omitempty"`
	XMLNS                            string                           `xml:"xmlns,attr,omitempty"`
	CustomerCreditTransferInitiation CustomerCreditTransferInitiation `xml:"CstmrCdtTrfInitn,omitempty"`
}

// CustomerCreditTransferInitiation represents the root of the customer credit transfer initiation.
type CustomerCreditTransferInitiation struct {
	GroupHeader pain.GroupHeader `xml:"GrpHdr,omitempty"`
	PaymentInfo []PaymentInfo    `xml:"PmtInf,omitempty"`
}

// PaymentInfo represents payment-related information.
type PaymentInfo struct {
	PaymentInfoID             string                        `xml:"PmtInfId,omitempty"`    // PaymentInfoID is a unique identifier for the payment info block.
	PaymentMethod             string                        `xml:"PmtMtd,omitempty"`      // PaymentMethod defines the payment method (e.g., TRF for Transfer).
	NumberOfTransactions      int                           `xml:"NbOfTxs,omitempty"`     // NumberOfTransactions is the total number of transactions.
	ControlSum                decimal.Decimal               `xml:"CtrlSum,omitempty"`     // ControlSum is the total amount for the transaction batch.
	BatchBooking              bool                          `xml:"BtchBookg,omitempty"`   // BatchBooking defines whether batch booking is enabled.
	PaymentTypeInfo           pain.PaymentTypeInfo          `xml:"PmtTpInf,omitempty"`    // PaymentTypeInfo provides additional payment type information.
	RequestedExecutionDate    string                        `xml:"ReqdExctnDt,omitempty"` // RequestedExecutionDate is the date on which the payment should be executed.
	Debtor                    pain.Party                    `xml:"Dbtr,omitempty"`        // Debtor contains information about the party sending the payment.
	DebtorAccount             pain.Account                  `xml:"DbtrAcct,omitempty"`    // DebtorAccount represents the account of the debtor.
	DebtorAgent               FinancialInstitution          `xml:"DbtrAgt,omitempty"`     // DebtorAgent represents the financial institution of the debtor.
	CreditTransferTransaction CreditTransferTransactionInfo `xml:"CdtTrfTxInf,omitempty"` // CreditTransferTransaction contains details of the credit transfer.
}

// FinancialInstitution represents the financial institution.
type FinancialInstitution struct {
	FinancialInstitutionID FinancialInstitutionID `xml:"FinInstnId,omitempty"` // FinancialInstitutionID holds the BIC of the institution.
}

// FinancialInstitutionID represents the BIC code of a financial institution.
type FinancialInstitutionID struct {
	BIC string `xml:"BIC,omitempty"` // BIC is the Bank Identifier Code.
}

// CreditTransferTransactionInfo contains information about a single credit transfer transaction.
type CreditTransferTransactionInfo struct {
	PaymentID       pain.PaymentID        `xml:"PmtId,omitempty"`    // PaymentID contains payment identifiers like instruction ID.
	PaymentTypeInfo pain.PaymentTypeInfo  `xml:"PmtTpInf,omitempty"` // PaymentTypeInfo is optional payment type information.
	Amount          Amount                `xml:"Amt,omitempty"`      // Amount represents the instructed amount for the transfer.
	CreditorAgent   *FinancialInstitution `xml:"CdtrAgt,omitempty"`  // CreditorAgent is the financial institution of the creditor.
	Creditor        pain.Party            `xml:"Cdtr,omitempty"`     // Creditor contains information about the party receiving the payment.
	CreditorAccount pain.Account          `xml:"CdtrAcct,omitempty"` // CreditorAccount represents the account of the creditor.
	RemittanceInfo  RemittanceInformation `xml:"RmtInf,omitempty"`   // RemittanceInfo contains the remittance details.
}

// Amount represents a monetary value and its associated currency.
type Amount struct {
	InstigatedAmount pain.InstigatedAmount `xml:"InstdAmt,omitempty"` // InstructedAmount specifies a monetary amount and its currency in a transaction.
}

// InstructedAmount specifies a monetary amount and its currency in a transaction.
type InstructedAmount struct {
	Amount   decimal.Decimal `xml:",chardata"` // Amount represents a monetary value and its associated currency.
	Currency string          `xml:"Ccy,attr"`  // Currency represents the currency in which the monetary amount is specified.
}

// RemittanceInformation contains remittance details for the transaction.
type RemittanceInformation struct {
	Unstructured string `xml:"Ustrd,omitempty"` // Unstructured contains free-form remittance information.
}

func NewDocument(customerCreditTransferInitiation CustomerCreditTransferInitiation) Document {
	return Document{
		XMLNS:                            "urn:iso:std:iso:20022:tech:xsd:pain.001.001.03",
		CustomerCreditTransferInitiation: customerCreditTransferInitiation,
	}
}
