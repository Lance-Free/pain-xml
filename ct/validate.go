package ct

import (
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/lance-free/pain-xml/internal/iso7064"
	"github.com/lance-free/pain-xml/pain"
	"github.com/lance-free/pain-xml/reference"
	"github.com/shopspring/decimal"
)

const (
	currencyEUR             = "EUR"
	dateLayout              = "2006-01-02"
	maxIdentifierLength     = 35
	maxNameLength           = 70
	maxAddressLineLength    = 70
	maxAddressLines         = 2
	maxUnstructuredLength   = 140
	maxStructuredLength     = 140
	maxAmountDecimals       = 2
	identifierSeparator     = "/"
	identifierDoubleSlashes = "//"
)

var (
	minAmount         = decimal.RequireFromString("0.01")
	maxAmount         = decimal.RequireFromString("999999999.99")
	latinCharacters   = regexp.MustCompile(`^[a-zA-Z0-9/\-?:().,'+ ]*$`)
	ibanPattern       = regexp.MustCompile(`^[A-Z]{2}[0-9]{2}[A-Z0-9]{11,30}$`)
	bicPattern        = regexp.MustCompile(`^[A-Z]{6}[A-Z2-9][A-NP-Z0-9]([A-Z0-9]{3})?$`)
	countryPattern    = regexp.MustCompile(`^[A-Z]{2}$`)
	structuredTagSize = len("<Strd></Strd>")
)

type violations []error

func (v *violations) add(path, format string, args ...any) {
	*v = append(*v, fmt.Errorf("%s: %s", path, fmt.Sprintf(format, args...)))
}

// Validate checks the document against the SEPA Credit Transfer customer-to-bank rules (EPC132-08) and the Dutch usage
// rules. It returns all violations joined into one error, or nil.
func (document Document) Validate() error {
	var v violations
	initiation := document.CustomerCreditTransferInitiation

	validateGroupHeader(&v, initiation)
	for index, info := range initiation.PaymentInfo {
		validatePaymentInfo(&v, fmt.Sprintf("PmtInf[%d]", index), info)
	}

	return errors.Join(v...)
}

func validateGroupHeader(v *violations, initiation CustomerCreditTransferInitiation) {
	header := initiation.GroupHeader
	validateIdentifier(v, "GrpHdr.MsgId", header.MessageID, true)
	validateText(v, "GrpHdr.InitgPty.Nm", header.InitiatingParty.Name, maxNameLength, false)

	if len(initiation.PaymentInfo) == 0 {
		v.add("PmtInf", "at least one payment info block is required")
	}

	transactions := 0
	sum := decimal.Zero
	for _, info := range initiation.PaymentInfo {
		transactions += len(info.CreditTransferTransactions)
		sum = sum.Add(controlSum(info.CreditTransferTransactions))
	}

	if header.NumberOfTransactions != transactions {
		v.add("GrpHdr.NbOfTxs", "is %d, but the document has %d transactions", header.NumberOfTransactions, transactions)
	}

	if !header.ControlSum.Equal(sum) {
		v.add("GrpHdr.CtrlSum", "is %s, but the transactions sum to %s", header.ControlSum, sum)
	}
}

func validatePaymentInfo(v *violations, path string, info PaymentInfo) {
	validateIdentifier(v, path+".PmtInfId", info.PaymentInfoID, true)

	if info.PaymentMethod != PaymentMethodTransfer {
		v.add(path+".PmtMtd", "must be %s", PaymentMethodTransfer)
	}

	if len(info.CreditTransferTransactions) == 0 {
		v.add(path+".CdtTrfTxInf", "at least one transaction is required")
	}

	if info.NumberOfTransactions != len(info.CreditTransferTransactions) {
		v.add(path+".NbOfTxs", "is %d, but the block has %d transactions", info.NumberOfTransactions, len(info.CreditTransferTransactions))
	}

	if sum := controlSum(info.CreditTransferTransactions); !info.ControlSum.Equal(sum) {
		v.add(path+".CtrlSum", "is %s, but the transactions sum to %s", info.ControlSum, sum)
	}

	if _, err := time.Parse(dateLayout, info.RequestedExecutionDate); err != nil {
		v.add(path+".ReqdExctnDt", "must be a date in the format YYYY-MM-DD")
	}

	validateParty(v, path+".Dbtr", info.Debtor, true)
	validateIBAN(v, path+".DbtrAcct", info.DebtorAccount.ID.IBAN)
	validateBIC(v, path+".DbtrAgt", info.DebtorAgent.FinancialInstitutionID.BIC)
	validateChargeBearer(v, path+".ChrgBr", info.ChargeBearer)

	if info.UltimateDebtor != nil {
		validateParty(v, path+".UltmtDbtr", *info.UltimateDebtor, false)
	}

	for index, transaction := range info.CreditTransferTransactions {
		transactionPath := fmt.Sprintf("%s.CdtTrfTxInf[%d]", path, index)
		validateTransaction(v, transactionPath, transaction)

		if info.UltimateDebtor != nil && transaction.UltimateDebtor != nil {
			v.add(transactionPath+".UltmtDbtr", "is also set on the payment info block; set it on one level only")
		}
	}
}

func validateTransaction(v *violations, path string, transaction CreditTransferTransactionInfo) {
	validateIdentifier(v, path+".PmtId.InstrId", transaction.PaymentID.InstructionID, false)
	validateIdentifier(v, path+".PmtId.EndToEndId", transaction.PaymentID.EndToEndID, true)
	validateAmount(v, path+".Amt.InstdAmt", transaction.Amount.InstigatedAmount)
	validateChargeBearer(v, path+".ChrgBr", transaction.ChargeBearer)
	validateParty(v, path+".Cdtr", transaction.Creditor, true)
	validateIBAN(v, path+".CdtrAcct", transaction.CreditorAccount.ID.IBAN)
	validateRemittance(v, path+".RmtInf", transaction.RemittanceInfo)

	if transaction.UltimateDebtor != nil {
		validateParty(v, path+".UltmtDbtr", *transaction.UltimateDebtor, false)
	}

	if transaction.CreditorAgent != nil {
		validateBIC(v, path+".CdtrAgt", transaction.CreditorAgent.FinancialInstitutionID.BIC)
	}
}

func validateAmount(v *violations, path string, amount pain.InstigatedAmount) {
	if amount.Currency != currencyEUR {
		v.add(path, "currency must be %s", currencyEUR)
	}

	isInRange := amount.Text.GreaterThanOrEqual(minAmount) && amount.Text.LessThanOrEqual(maxAmount)
	if !isInRange {
		v.add(path, "must be between %s and %s", minAmount, maxAmount)
	}

	if !amount.Text.Equal(amount.Text.Round(maxAmountDecimals)) {
		v.add(path, "must have at most %d decimals", maxAmountDecimals)
	}
}

func validateRemittance(v *violations, path string, remittance RemittanceInformation) {
	validateText(v, path+".Ustrd", remittance.Unstructured, maxUnstructuredLength, false)

	structured := remittance.Structured
	if structured == nil {
		return
	}

	if remittance.Unstructured != "" {
		v.add(path, "must contain Ustrd or Strd, not both")
	}

	if length := structuredContentLength(*structured); length > maxStructuredLength {
		v.add(path+".Strd", "content is %d characters, the maximum is %d", length, maxStructuredLength)
	}

	creditorReference := structured.CreditorReferenceInformation
	if creditorReference.Type.CodeOrProprietary.Code != CreditorReferenceTypeCode {
		v.add(path+".Strd.CdtrRefInf.Tp.CdOrPrtry.Cd", "must be %s", CreditorReferenceTypeCode)
	}

	validateIdentifier(v, path+".Strd.CdtrRefInf.Ref", creditorReference.Reference, true)
	validateCreditorReference(v, path+".Strd.CdtrRefInf", creditorReference)
}

func validateCreditorReference(v *violations, path string, creditorReference CreditorReferenceInformation) {
	issuer := reference.Issuer(creditorReference.Type.Issuer)
	isCheckedIssuer := issuer == reference.IssuerCurrence || issuer == reference.IssuerISO
	if !isCheckedIssuer {
		return
	}

	parsed, ok := reference.Parse(creditorReference.Reference)
	isValid := ok && parsed.Issuer == issuer && parsed.Value == creditorReference.Reference
	if !isValid {
		v.add(path+".Ref", "is not a valid %s reference in electronic format", issuer)
	}
}

func structuredContentLength(structured StructuredRemittanceInformation) int {
	encoded, err := xml.Marshal(structured)
	if err != nil {
		return 0
	}

	return utf8.RuneCount(encoded) - structuredTagSize
}

func validateParty(v *violations, path string, party pain.Party, nameRequired bool) {
	validateText(v, path+".Nm", party.Name, maxNameLength, nameRequired)

	if party.PostalAddress != nil {
		validateAddress(v, path+".PstlAdr", *party.PostalAddress)
	}

	if party.Identification != nil {
		validatePartyIdentification(v, path+".Id", *party.Identification)
	}
}

func validateAddress(v *violations, path string, address pain.PostalAddress) {
	if address.Country != "" && !countryPattern.MatchString(address.Country) {
		v.add(path+".Ctry", "must be an ISO 3166 alpha-2 country code")
	}

	if len(address.AddressLine) > maxAddressLines {
		v.add(path+".AdrLine", "at most %d lines are allowed", maxAddressLines)
	}

	for index, line := range address.AddressLine {
		validateText(v, fmt.Sprintf("%s.AdrLine[%d]", path, index), line, maxAddressLineLength, false)
	}
}

func validatePartyIdentification(v *violations, path string, identification pain.PartyIdentification) {
	hasOrganisation := identification.Organisation != nil
	hasPrivate := identification.Private != nil
	if hasOrganisation == hasPrivate {
		v.add(path, "must contain OrgId or PrvtId, not both or neither")
		return
	}

	if hasPrivate {
		validateOtherIdentification(v, path+".PrvtId.Othr", identification.Private.Other)
		return
	}

	organisation := identification.Organisation
	hasBICOrBEI := organisation.BICOrBEI != ""
	hasOther := organisation.Other != nil
	if hasBICOrBEI == hasOther {
		v.add(path+".OrgId", "must contain BICOrBEI or Othr, not both or neither")
		return
	}

	if hasBICOrBEI {
		validateBIC(v, path+".OrgId.BICOrBEI", organisation.BICOrBEI)
		return
	}

	validateOtherIdentification(v, path+".OrgId.Othr", organisation.Other)
}

func validateOtherIdentification(v *violations, path string, other *pain.GenericIdentification) {
	if other == nil {
		v.add(path, "is required")
		return
	}

	validateText(v, path+".Id", other.ID, maxIdentifierLength, true)

	if other.SchemeName == nil {
		return
	}

	hasCode := other.SchemeName.Code != ""
	hasProprietary := other.SchemeName.Proprietary != ""
	if hasCode == hasProprietary {
		v.add(path+".SchmeNm", "must contain Cd or Prtry, not both or neither")
	}
}

func validateIBAN(v *violations, path, iban string) {
	isValid := ibanPattern.MatchString(iban) && iso7064.Valid(iban)
	if !isValid {
		v.add(path+".Id.IBAN", "must be a valid IBAN in upper case without spaces")
	}
}

func validateBIC(v *violations, path, bic string) {
	if bic == "" || bicPattern.MatchString(bic) {
		return
	}

	v.add(path, "must be a valid BIC")
}

func validateChargeBearer(v *violations, path, chargeBearer string) {
	if chargeBearer == "" || chargeBearer == ChargeBearerFollowingServiceLevel {
		return
	}

	v.add(path, "must be %s", ChargeBearerFollowingServiceLevel)
}

func validateText(v *violations, path, text string, maxLength int, required bool) {
	length := utf8.RuneCountInString(text)
	if required && length == 0 {
		v.add(path, "is required")
		return
	}

	if length > maxLength {
		v.add(path, "is %d characters, the maximum is %d", length, maxLength)
	}
}

func validateIdentifier(v *violations, path, identifier string, required bool) {
	validateText(v, path, identifier, maxIdentifierLength, required)

	if !latinCharacters.MatchString(identifier) {
		v.add(path, "must contain only the SEPA Latin character set")
	}

	hasSeparatorAtEdge := strings.HasPrefix(identifier, identifierSeparator) || strings.HasSuffix(identifier, identifierSeparator)
	if hasSeparatorAtEdge || strings.Contains(identifier, identifierDoubleSlashes) {
		v.add(path, "must not start or end with %q or contain %q", identifierSeparator, identifierDoubleSlashes)
	}
}
