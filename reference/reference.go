package reference

import (
	"strings"

	"github.com/lance-free/pain-xml/internal/iso7064"
)

// Issuer is the issuer of a creditor reference scheme, as used in the CdtrRefInf/Tp/Issr element.
type Issuer string

const (
	IssuerCurrence Issuer = "CUR" // IssuerCurrence issues the Dutch payment reference ("betalingskenmerk").
	IssuerISO      Issuer = "ISO" // IssuerISO issues the ISO 11649 RF creditor reference.
)

const (
	dutchLength           = 16
	dutchModulus          = 11
	dutchCheckDigitTen    = 10
	dutchCheckDigitEleven = 11
	rfPrefix              = "RF"
	rfMinLength           = 5
	rfMaxLength           = 25
)

// dutchWeights are the weights of the digits, counted from the rightmost digit to the left. They repeat for longer values.
var dutchWeights = [...]int{2, 4, 8, 5, 10, 9, 7, 3, 6, 1}

// CreditorReference is a creditor reference with a valid check digit, in its electronic format.
type CreditorReference struct {
	Value  string // Value is the reference without spaces and in upper case.
	Issuer Issuer // Issuer is the issuer of the reference scheme.
}

// Parse returns the creditor reference in raw when raw is a 16-digit Dutch payment reference or an ISO 11649 RF
// creditor reference with a valid check digit. Spaces are removed and letters are converted to upper case first.
// Shorter Dutch payment references are not recognised, because banks do not accept all of them as structured.
func Parse(raw string) (CreditorReference, bool) {
	value := strings.ToUpper(strings.ReplaceAll(raw, " ", ""))

	if isDutch(value) {
		return CreditorReference{Value: value, Issuer: IssuerCurrence}, true
	}

	if isRF(value) {
		return CreditorReference{Value: value, Issuer: IssuerISO}, true
	}

	return CreditorReference{}, false
}

func isDutch(value string) bool {
	if len(value) != dutchLength || !onlyDigits(value) {
		return false
	}

	sum := 0
	for position := 1; position < len(value); position++ {
		weight := dutchWeights[(len(value)-1-position)%len(dutchWeights)]
		sum += int(value[position]-'0') * weight
	}

	return dutchCheckDigit(sum) == int(value[0]-'0')
}

func dutchCheckDigit(sum int) int {
	digit := dutchModulus - sum%dutchModulus
	switch digit {
	case dutchCheckDigitTen:
		return 1
	case dutchCheckDigitEleven:
		return 0
	}

	return digit
}

func isRF(value string) bool {
	hasValidLength := len(value) >= rfMinLength && len(value) <= rfMaxLength
	if !hasValidLength || !strings.HasPrefix(value, rfPrefix) {
		return false
	}

	return iso7064.Valid(value)
}

func onlyDigits(value string) bool {
	for _, char := range value {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}
