package ct

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"

	"github.com/lance-free/pain-xml/internal/schematest"
	"github.com/lance-free/pain-xml/pain"
	"github.com/shopspring/decimal"
)

func amount(value string) Amount {
	return Amount{InstigatedAmount: pain.InstigatedAmount{Currency: currencyEUR, Text: decimal.RequireFromString(value)}}
}

func validDocument() Document {
	batchBooking := false

	return NewDocument(CustomerCreditTransferInitiation{
		GroupHeader: pain.GroupHeader{
			MessageID:        "MSG-20261104-150405-0a1b2c3d",
			CreationDateTime: pain.CreationDateTime(time.Date(2026, 11, 4, 15, 4, 5, 0, time.UTC)),
			InitiatingParty:  pain.InitiatingParty{Name: "PayToday B.V."},
		},
		PaymentInfo: []PaymentInfo{
			{
				PaymentInfoID:          "PMT-20261104-150405-0a1b2c3d",
				PaymentMethod:          PaymentMethodTransfer,
				BatchBooking:           &batchBooking,
				PaymentTypeInfo:        pain.PaymentTypeInfo{InstructionPriority: "NORM", ServiceLevel: pain.ServiceLevel{Code: "SEPA"}},
				RequestedExecutionDate: "2026-11-04",
				Debtor:                 pain.Party{Name: "PayToday B.V."},
				DebtorAccount:          pain.Account{ID: pain.AccountID{IBAN: "NL91ABNA0417164300"}},
				DebtorAgent:            FinancialInstitution{FinancialInstitutionID: FinancialInstitutionID{BIC: "ABNANL2A"}},
				ChargeBearer:           ChargeBearerFollowingServiceLevel,
				CreditTransferTransactions: []CreditTransferTransactionInfo{
					{
						PaymentID: pain.PaymentID{EndToEndID: "5000056789012345"},
						Amount:    amount("1210.50"),
						UltimateDebtor: &pain.Party{
							Name: "Debiteur B.V.",
							PostalAddress: &pain.PostalAddress{
								StreetName:     "Damrak",
								BuildingNumber: "1",
								PostCode:       "1012 LG",
								TownName:       "Amsterdam",
								Country:        "NL",
							},
							Identification: &pain.PartyIdentification{
								Organisation: &pain.OrganisationIdentification{
									Other: &pain.GenericIdentification{ID: "12345678", SchemeName: &pain.SchemeName{Proprietary: "KVK"}},
								},
							},
						},
						Creditor:        pain.Party{Name: "Klant B.V."},
						CreditorAccount: pain.Account{ID: pain.AccountID{IBAN: "DE89370400440532013000"}},
						RemittanceInfo:  NewRemittanceInformation("5000 0567 8901 2345", "Doorbetaling factuur 2026-0042"),
					},
					{
						PaymentID:       pain.PaymentID{InstructionID: "INSTR-2", EndToEndID: "NOTPROVIDED"},
						Amount:          amount("99.99"),
						CreditorAgent:   &FinancialInstitution{FinancialInstitutionID: FinancialInstitutionID{BIC: "DEUTDEFF500"}},
						Creditor:        pain.Party{Name: "Andere Klant B.V.", PostalAddress: &pain.PostalAddress{Country: "NL", AddressLine: []string{"Coolsingel 40", "3011 AD Rotterdam"}}},
						CreditorAccount: pain.Account{ID: pain.AccountID{IBAN: "NL91ABNA0417164300"}},
						RemittanceInfo:  NewRemittanceInformation("2026-0043", "Doorbetaling factuur 2026-0043"),
					},
				},
			},
		},
	})
}

func TestNewDocumentSetsTotals(t *testing.T) {
	document := validDocument()
	header := document.CustomerCreditTransferInitiation.GroupHeader
	info := document.CustomerCreditTransferInitiation.PaymentInfo[0]
	wantSum := decimal.RequireFromString("1310.49")

	if header.NumberOfTransactions != 2 || !header.ControlSum.Equal(wantSum) {
		t.Fatalf("group header totals = %d, %s; want 2, %s", header.NumberOfTransactions, header.ControlSum, wantSum)
	}

	if info.NumberOfTransactions != 2 || !info.ControlSum.Equal(wantSum) {
		t.Fatalf("payment info totals = %d, %s; want 2, %s", info.NumberOfTransactions, info.ControlSum, wantSum)
	}
}

func TestNewDocumentDoesNotChangeInput(t *testing.T) {
	paymentInfo := []PaymentInfo{{CreditTransferTransactions: []CreditTransferTransactionInfo{{Amount: amount("1.00")}}}}
	NewDocument(CustomerCreditTransferInitiation{PaymentInfo: paymentInfo})

	if paymentInfo[0].NumberOfTransactions != 0 {
		t.Fatal("NewDocument changed the payment info of the caller")
	}
}

func TestNewRemittanceInformation(t *testing.T) {
	tests := []struct {
		name        string
		reference   string
		description string
		wantIssuer  string
		wantRef     string
		wantText    string
	}{
		{name: "dutch payment reference", reference: "5000 0567 8901 2345", description: "Factuur 2026-0042", wantIssuer: "CUR", wantRef: "5000056789012345"},
		{name: "rf creditor reference", reference: "RF18 5390 0754 7034", description: "Factuur 2026-0042", wantIssuer: "ISO", wantRef: "RF18539007547034"},
		{name: "other reference", reference: "2026-0042", description: "Factuur 2026-0042", wantText: "Factuur 2026-0042"},
		{name: "other reference without description", reference: "2026-0042", wantText: "2026-0042"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewRemittanceInformation(test.reference, test.description)
			if got.Unstructured != test.wantText {
				t.Fatalf("Unstructured = %q; want %q", got.Unstructured, test.wantText)
			}

			if test.wantRef == "" {
				if got.Structured != nil {
					t.Fatalf("Structured = %+v; want nil", got.Structured)
				}
				return
			}

			creditorReference := got.Structured.CreditorReferenceInformation
			if creditorReference.Reference != test.wantRef || creditorReference.Type.Issuer != test.wantIssuer || creditorReference.Type.CodeOrProprietary.Code != CreditorReferenceTypeCode {
				t.Fatalf("CdtrRefInf = %+v; want Ref %q, Issr %q, Cd SCOR", creditorReference, test.wantRef, test.wantIssuer)
			}
		})
	}
}

func TestDocumentMatchesSchema(t *testing.T) {
	schematest.Validate(t, validDocument())
}

func TestDocumentContainsNewElements(t *testing.T) {
	encoded, err := xml.Marshal(validDocument())
	if err != nil {
		t.Fatalf("marshalling document: %v", err)
	}

	for _, want := range []string{
		"<BtchBookg>false</BtchBookg>",
		"<ChrgBr>SLEV</ChrgBr>",
		"<UltmtDbtr><Nm>Debiteur B.V.</Nm>",
		"<Id><OrgId><Othr><Id>12345678</Id><SchmeNm><Prtry>KVK</Prtry></SchmeNm></Othr></OrgId></Id>",
		"<Strd><CdtrRefInf><Tp><CdOrPrtry><Cd>SCOR</Cd></CdOrPrtry><Issr>CUR</Issr></Tp><Ref>5000056789012345</Ref></CdtrRefInf></Strd>",
		"<Ustrd>Doorbetaling factuur 2026-0043</Ustrd>",
	} {
		if !strings.Contains(string(encoded), want) {
			t.Errorf("document does not contain %s", want)
		}
	}
}

func TestValidateAcceptsValidDocument(t *testing.T) {
	if err := validDocument().Validate(); err != nil {
		t.Fatalf("Validate() = %v; want nil", err)
	}
}

func TestValidateReportsViolations(t *testing.T) {
	tests := []struct {
		name   string
		change func(*Document)
		want   string
	}{
		{
			name:   "control sum does not match",
			change: func(d *Document) { d.CustomerCreditTransferInitiation.GroupHeader.ControlSum = decimal.NewFromInt(1) },
			want:   "GrpHdr.CtrlSum",
		},
		{
			name:   "message id too long",
			change: func(d *Document) { d.CustomerCreditTransferInitiation.GroupHeader.MessageID = strings.Repeat("A", 36) },
			want:   "GrpHdr.MsgId: is 36 characters",
		},
		{
			name:   "invalid IBAN",
			change: func(d *Document) { transaction(d, 0).CreditorAccount.ID.IBAN = "NL91ABNA0417164301" },
			want:   "CdtTrfTxInf[0].CdtrAcct.Id.IBAN",
		},
		{
			name:   "invalid BIC",
			change: func(d *Document) { transaction(d, 1).CreditorAgent.FinancialInstitutionID.BIC = "DEUT" },
			want:   "CdtTrfTxInf[1].CdtrAgt",
		},
		{
			name:   "amount with three decimals",
			change: func(d *Document) { transaction(d, 1).Amount = amount("1.001") },
			want:   "at most 2 decimals",
		},
		{
			name:   "amount in other currency",
			change: func(d *Document) { transaction(d, 1).Amount.InstigatedAmount.Currency = "USD" },
			want:   "currency must be EUR",
		},
		{
			name:   "end-to-end id with double slash",
			change: func(d *Document) { transaction(d, 1).PaymentID.EndToEndID = "A//B" },
			want:   "PmtId.EndToEndId: must not start or end",
		},
		{
			name:   "end-to-end id with other characters",
			change: func(d *Document) { transaction(d, 1).PaymentID.EndToEndID = "factuur_42" },
			want:   "PmtId.EndToEndId: must contain only the SEPA Latin character set",
		},
		{
			name:   "creditor name too long",
			change: func(d *Document) { transaction(d, 1).Creditor.Name = strings.Repeat("A", 71) },
			want:   "CdtTrfTxInf[1].Cdtr.Nm: is 71 characters",
		},
		{
			name: "unstructured and structured",
			change: func(d *Document) {
				transaction(d, 0).RemittanceInfo.Unstructured = "Factuur 42"
			},
			want: "must contain Ustrd or Strd, not both",
		},
		{
			name: "invalid Dutch payment reference",
			change: func(d *Document) {
				transaction(d, 0).RemittanceInfo.Structured.CreditorReferenceInformation.Reference = "4000056789012345"
			},
			want: "Strd.CdtrRefInf.Ref: is not a valid CUR reference",
		},
		{
			name: "ultimate debtor on two levels",
			change: func(d *Document) {
				d.CustomerCreditTransferInitiation.PaymentInfo[0].UltimateDebtor = &pain.Party{Name: "Debiteur B.V."}
			},
			want: "CdtTrfTxInf[0].UltmtDbtr: is also set on the payment info block",
		},
		{
			name: "organisation with two identifications",
			change: func(d *Document) {
				transaction(d, 0).UltimateDebtor.Identification.Organisation.BICOrBEI = "ABNANL2A"
			},
			want: "UltmtDbtr.Id.OrgId: must contain BICOrBEI or Othr",
		},
		{
			name: "three address lines",
			change: func(d *Document) {
				address := transaction(d, 1).Creditor.PostalAddress
				address.AddressLine = append(address.AddressLine, "Nederland")
			},
			want: "Cdtr.PstlAdr.AdrLine: at most 2 lines",
		},
		{
			name:   "other charge bearer",
			change: func(d *Document) { d.CustomerCreditTransferInitiation.PaymentInfo[0].ChargeBearer = "DEBT" },
			want:   "PmtInf[0].ChrgBr: must be SLEV",
		},
		{
			name: "invalid execution date",
			change: func(d *Document) {
				d.CustomerCreditTransferInitiation.PaymentInfo[0].RequestedExecutionDate = "04-11-2026"
			},
			want: "PmtInf[0].ReqdExctnDt",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			document := validDocument()
			test.change(&document)

			err := document.Validate()
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Validate() = %v; want an error that contains %q", err, test.want)
			}
		})
	}
}

func transaction(document *Document, index int) *CreditTransferTransactionInfo {
	return &document.CustomerCreditTransferInitiation.PaymentInfo[0].CreditTransferTransactions[index]
}
