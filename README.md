# PainXML

The `painxml` package facilitates the generation and parsing of XML documents conforming to the ISO 20022 standard for
payment initiation, specifically focusing on Direct Debit and Credit Transfer documents. 

## Credit transfers (pain.001.001.03)

Package `sct` builds a SEPA Credit Transfer from plain payment data. It sets the values that SEPA fixes (payment
method `TRF`, service level `SEPA`, currency `EUR`, charge bearer `SLEV`, the number of transactions and the control
sums), removes spaces from IBANs, and validates the message against the SEPA rules before it returns the document.

```go
document, err := sct.Message{
	ID:     "MSG-1",
	Debtor: sct.Account{Name: "PayToday B.V.", IBAN: "NL91 ABNA 0417 1643 00", BIC: "ABNANL2A"},
	Transfers: []sct.Transfer{{
		EndToEndID: "E2E-1",
		Amount:     decimal.RequireFromString("1210.50"),
		Creditor:   sct.Account{Name: "Klant B.V.", IBAN: "DE89370400440532013000"},
		UltimateDebtor: &sct.Party{
			Name:           "Debiteur B.V.",
			Identification: sct.CompanyRegistration("12345678", "KVK"),
		},
		Reference:   "5000 0567 8901 2345",
		Description: "Factuur 2026-0042",
	}},
}.Document()
if err != nil {
	return err
}

paymentID, err := payment.Initiate(ctx, document, keyPair) // ABN AMRO SDK
```

- `Reference` is sent as a structured creditor reference (`Strd/CdtrRefInf`, type `SCOR`) when it is a 16-digit Dutch
  payment reference ("betalingskenmerk", issuer `CUR`) or an ISO 11649 RF creditor reference (issuer `ISO`) with a valid
  check digit. Otherwise `Description` is sent as unstructured remittance information, or `Reference` when
  `Description` is empty. SEPA does not allow both, so `Description` is not sent next to a structured reference.
- `UltimateDebtor` names the party on whose behalf the payment is made. Set it on the message or on the transfer, not
  on both.
- If `CreatedAt` is zero, the current time is used. If `ExecutionDate` is zero, the date of `CreatedAt` is used. If
  `EndToEndID` is empty, `NOTPROVIDED` is sent.
- The error contains all violations with their XML path, for example
  `PmtInf[0].CdtTrfTxInf[0].CdtrAcct.Id.IBAN: must be a valid IBAN in upper case without spaces`.

### Low-level documents

Package `ct` maps the pain.001.001.03 schema to Go structs, for documents that `sct` cannot express, such as more than
one payment info block. `ct.NewDocument` sets the number of transactions and the control sums,
`ct.NewRemittanceInformation` selects structured or unstructured remittance information, and `Document.Validate`
checks the document against the SEPA Credit Transfer customer-to-bank rules (EPC132-08) and the Dutch usage rules.
Package `reference` recognises Dutch payment references and RF creditor references.

## Upgrading from v0.3

- `PaymentInfo.CreditTransferTransaction` is now the slice `PaymentInfo.CreditTransferTransactions`.
- `PaymentInfo.BatchBooking` is now a `*bool`, so `false` can be sent. If it is nil, the bank default applies.
- `NewDocument` overwrites `NumberOfTransactions` and `ControlSum`. You do not need to set them.
- `ct.InstructedAmount` is removed. Use `pain.InstigatedAmount`.
- `dd.FinancialInstitutionIdentification.BICFI` is now `BIC`, as pain.008.001.02 requires, and package `dd` compiles again.

## License

This package is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
