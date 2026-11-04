package iso7064

const (
	modulus           = 97
	validRemainder    = 1
	decimalBase       = 10
	letterBase        = 100
	firstLetterNumber = 10
	rotatedLength     = 4
)

// Valid reports whether value passes the ISO 7064 MOD 97-10 check used by IBANs and RF creditor references.
// The first four characters are moved to the end and every letter counts as two digits (A = 10, ..., Z = 35).
// The value must be in upper case.
func Valid(value string) bool {
	if len(value) <= rotatedLength {
		return false
	}

	remainder, ok := remainderOf(value[rotatedLength:] + value[:rotatedLength])
	return ok && remainder == validRemainder
}

func remainderOf(value string) (int, bool) {
	remainder := 0
	for _, char := range value {
		isDigit := char >= '0' && char <= '9'
		isLetter := char >= 'A' && char <= 'Z'
		if !isDigit && !isLetter {
			return 0, false
		}

		if isDigit {
			remainder = (remainder*decimalBase + int(char-'0')) % modulus
			continue
		}

		remainder = (remainder*letterBase + int(char-'A') + firstLetterNumber) % modulus
	}

	return remainder, true
}
