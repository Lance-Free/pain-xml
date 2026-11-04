package ct

import (
	"encoding/xml"
	"slices"

	"github.com/lance-free/pain-xml/pain"
	"github.com/lance-free/pain-xml/reference"
	"github.com/shopspring/decimal"
)

const (
	Namespace                         = "urn:iso:std:iso:20022:tech:xsd:pain.001.001.03" // Namespace is the XML namespace of pain.001.001.03.
	PaymentMethodTransfer             = "TRF"                                            // PaymentMethodTransfer is the only payment method SEPA allows.
	ChargeBearerFollowingServiceLevel = "SLEV"                                           // ChargeBearerFollowingServiceLevel is the only charge bearer SEPA allows.
	CreditorReferenceTypeCode         = "SCOR"                                           // CreditorReferenceTypeCode is the only creditor reference type SEPA allows.
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
	PaymentInfoID              string                          `xml:"PmtInfId,omitempty"`    // PaymentInfoID is a unique identifier for the payment info block.
	PaymentMethod              string                          `xml:"PmtMtd,omitempty"`      // PaymentMethod defines the payment method (e.g., TRF for Transfer).
	BatchBooking               *bool                           `xml:"BtchBookg,omitempty"`   // BatchBooking requests one booking for all transactions (true) or one booking per transaction (false). If nil, the bank default applies.
	NumberOfTransactions       int                             `xml:"NbOfTxs,omitempty"`     // NumberOfTransactions is the total number of transactions. NewDocument sets it.
	ControlSum                 decimal.Decimal                 `xml:"CtrlSum,omitempty"`     // ControlSum is the total amount for the transaction batch. NewDocument sets it.
	PaymentTypeInfo            pain.PaymentTypeInfo            `xml:"PmtTpInf,omitempty"`    // PaymentTypeInfo provides additional payment type information.
	RequestedExecutionDate     string                          `xml:"ReqdExctnDt,omitempty"` // RequestedExecutionDate is the date on which the payment should be executed.
	Debtor                     pain.Party                      `xml:"Dbtr,omitempty"`        // Debtor contains information about the party sending the payment.
	DebtorAccount              pain.Account                    `xml:"DbtrAcct,omitempty"`    // DebtorAccount represents the account of the debtor.
	DebtorAgent                FinancialInstitution            `xml:"DbtrAgt,omitempty"`     // DebtorAgent represents the financial institution of the debtor.
	UltimateDebtor             *pain.Party                     `xml:"UltmtDbtr,omitempty"`   // UltimateDebtor is the party on whose behalf all transactions are paid. Set it here or per transaction, not both.
	ChargeBearer               string                          `xml:"ChrgBr,omitempty"`      // ChargeBearer specifies who bears the charges. SEPA allows only SLEV.
	CreditTransferTransactions []CreditTransferTransactionInfo `xml:"CdtTrfTxInf,omitempty"` // CreditTransferTransactions contains the credit transfers of the payment info block.
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
	PaymentID       pain.PaymentID        `xml:"PmtId,omitempty"`     // PaymentID contains payment identifiers like instruction ID.
	PaymentTypeInfo *pain.PaymentTypeInfo `xml:"PmtTpInf,omitempty"`  // PaymentTypeInfo is optional payment type information.
	Amount          Amount                `xml:"Amt,omitempty"`       // Amount represents the instructed amount for the transfer.
	ChargeBearer    string                `xml:"ChrgBr,omitempty"`    // ChargeBearer specifies who bears the charges. SEPA allows only SLEV.
	UltimateDebtor  *pain.Party           `xml:"UltmtDbtr,omitempty"` // UltimateDebtor is the party on whose behalf this transaction is paid.
	CreditorAgent   *FinancialInstitution `xml:"CdtrAgt,omitempty"`   // CreditorAgent is the financial institution of the creditor.
	Creditor        pain.Party            `xml:"Cdtr,omitempty"`      // Creditor contains information about the party receiving the payment.
	CreditorAccount pain.Account          `xml:"CdtrAcct,omitempty"`  // CreditorAccount represents the account of the creditor.
	RemittanceInfo  RemittanceInformation `xml:"RmtInf,omitempty"`    // RemittanceInfo contains the remittance details.
}

// Amount represents a monetary value and its associated currency.
type Amount struct {
	InstigatedAmount pain.InstigatedAmount `xml:"InstdAmt,omitempty"` // InstructedAmount specifies a monetary amount and its currency in a transaction.
}

// RemittanceInformation contains remittance details for the transaction. SEPA allows Unstructured or Structured, not both.
type RemittanceInformation struct {
	Unstructured string                           `xml:"Ustrd,omitempty"` // Unstructured contains free-form remittance information.
	Structured   *StructuredRemittanceInformation `xml:"Strd,omitempty"`  // Structured contains the creditor reference.
}

// StructuredRemittanceInformation contains structured remittance details.
type StructuredRemittanceInformation struct {
	XMLName                      xml.Name                     `xml:"Strd"`
	CreditorReferenceInformation CreditorReferenceInformation `xml:"CdtrRefInf"` // CreditorReferenceInformation contains the reference the creditor asked to be paid with.
}

// CreditorReferenceInformation contains the reference the creditor asked to be paid with.
type CreditorReferenceInformation struct {
	Type      CreditorReferenceType `xml:"Tp"`  // Type is the type and the issuer of the reference.
	Reference string                `xml:"Ref"` // Reference is the creditor reference.
}

// CreditorReferenceType is the type and the issuer of a creditor reference.
type CreditorReferenceType struct {
	CodeOrProprietary CreditorReferenceTypeChoice `xml:"CdOrPrtry"`      // CodeOrProprietary is the type of the reference.
	Issuer            string                      `xml:"Issr,omitempty"` // Issuer is the issuer of the reference scheme (e.g., CUR or ISO).
}

// CreditorReferenceTypeChoice is the code of a creditor reference type.
type CreditorReferenceTypeChoice struct {
	Code string `xml:"Cd"` // Code is the type code. SEPA allows only SCOR.
}

// NewRemittanceInformation returns structured remittance information with ref as the creditor reference if ref is a
// valid Dutch payment reference or RF creditor reference (see reference.Parse). SEPA does not allow both, so the
// description is then not sent. Otherwise it returns unstructured remittance information with description, or with ref
// if description is empty.
func NewRemittanceInformation(ref, description string) RemittanceInformation {
	creditorReference, ok := reference.Parse(ref)
	if !ok && description == "" {
		return RemittanceInformation{Unstructured: ref}
	}

	if !ok {
		return RemittanceInformation{Unstructured: description}
	}

	return RemittanceInformation{
		Structured: &StructuredRemittanceInformation{
			CreditorReferenceInformation: CreditorReferenceInformation{
				Type: CreditorReferenceType{
					CodeOrProprietary: CreditorReferenceTypeChoice{Code: CreditorReferenceTypeCode},
					Issuer:            string(creditorReference.Issuer),
				},
				Reference: creditorReference.Value,
			},
		},
	}
}

// NewDocument returns the document of the initiation. It sets the number of transactions and the control sum of the
// group header and of each payment info block from the transactions.
func NewDocument(customerCreditTransferInitiation CustomerCreditTransferInitiation) Document {
	paymentInfo := slices.Clone(customerCreditTransferInitiation.PaymentInfo)
	groupHeader := customerCreditTransferInitiation.GroupHeader
	groupHeader.NumberOfTransactions = 0
	groupHeader.ControlSum = decimal.Zero

	for index := range paymentInfo {
		info := &paymentInfo[index]
		info.NumberOfTransactions = len(info.CreditTransferTransactions)
		info.ControlSum = controlSum(info.CreditTransferTransactions)
		groupHeader.NumberOfTransactions += info.NumberOfTransactions
		groupHeader.ControlSum = groupHeader.ControlSum.Add(info.ControlSum)
	}

	return Document{
		XMLNS: Namespace,
		CustomerCreditTransferInitiation: CustomerCreditTransferInitiation{
			GroupHeader: groupHeader,
			PaymentInfo: paymentInfo,
		},
	}
}

func controlSum(transactions []CreditTransferTransactionInfo) decimal.Decimal {
	sum := decimal.Zero
	for _, transaction := range transactions {
		sum = sum.Add(transaction.Amount.InstigatedAmount.Text)
	}

	return sum
}
