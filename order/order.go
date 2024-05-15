package order

import (
	cryptoRand "crypto/rand"
	"encoding/xml"
	"fmt"
	"github.com/lance-free/pain-xml/document"
	"math/big"
	"time"
)

type Creditor struct {
	Name, Street, PostalCode, Place, Country, IBAN, BIC string
}

type Transaction struct {
	Name, Street, PostalCode, Place, Country, IBAN, BIC, Currency string
	Amount                                                        float64
}

type Order struct {
	ExecutionDate time.Time
	Transactions  []Transaction
	Creditor
}

// idGenerator generates a random string of 15 characters.
func idGenerator() (string, error) {
	const length = 15
	const alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := make([]byte, length)
	maxIndex := big.NewInt(int64(len(alphabet)))

	for i := range bytes {
		n, err := cryptoRand.Int(cryptoRand.Reader, maxIndex)
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		bytes[i] = alphabet[n.Int64()]
	}

	return string(bytes), nil
}

// controlSum calculates the sum of amounts in the given transactions slice and returns a formatted string with two decimal places.
func controlSum(transactions []Transaction) string {
	var sum float64
	for _, t := range transactions {
		sum += t.Amount
	}
	return fmt.Sprintf("%.2f", sum)
}

// ToDocument converts an order to a document.DirectDebit and returns it.
// The function takes an order of type order.Order and maps its properties
// to the corresponding properties of a DirectDebit.
func (order Order) ToDocument() (document.DirectDebit, error) {
	var transactions []document.DirectDebitTransactionInformation
	for _, t := range order.Transactions {
		transactions = append(transactions, document.DirectDebitTransactionInformation{
			PaymentID: document.PaymentID{EndToEndID: t.Name},
			InstigatedAmount: document.InstigatedAmount{
				Currency: t.Currency,
				Text:     fmt.Sprintf("%.2f", t.Amount),
			},
			DirectDebitTransaction: document.DirectDebitTransaction{
				MandateRelatedInformation: document.MandateRelatedInformation{
					MandateID:       order.Creditor.IBAN,
					DateOfSignature: order.ExecutionDate.Format("2006-01-02"),
				},
			},
			DebtorAgent: document.DebtorAgent{
				FinancialInstitutionIdentification: document.FinancialInstitutionIdentification{
					BICFI: order.Creditor.BIC,
				},
			},
			Debtor: document.Debtor{
				Name: t.Name,
				PostalAddress: document.PostalAddress{
					TownName:   t.Place,
					Country:    t.Country,
					StreetName: t.Street,
					PostalCode: t.PostalCode,
				},
			},
			DebtorAccount: document.DebtorAccount{
				IBAN: document.IBAN{IBAN: t.IBAN},
			},
			RemittanceInformation: document.RemittanceInformation{
				Unstructured: t.Name,
			},
		})
	}

	numberOfTransactions := fmt.Sprintf("%d", len(transactions))
	messageID, err := idGenerator()
	if err != nil {
		return document.DirectDebit{}, fmt.Errorf("failed to generate message ID: %w", err)
	}
	paymentInformationId, err := idGenerator()
	if err != nil {
		return document.DirectDebit{}, fmt.Errorf("failed to generate payment information ID: %w", err)
	}
	return document.DirectDebit{
		Xmlns: "urn:iso:std:iso:20022:tech:xsd:pain.008.001.08",
		CustomerDirectDebitInitiation: document.CustomerDirectDebitInitiation{
			GroupHeader: document.GroupHeader{
				MessageID:            messageID,
				CreationDateTime:     order.ExecutionDate.Format("2006-01-02T15:04:05"),
				NumberOfTransactions: numberOfTransactions,
			},
			PaymentInformation: document.PaymentInformation{
				PaymentInformationId: "Incasso SDD" + paymentInformationId,
				PaymentMethod:        "DD",
				NumberOfTransactions: numberOfTransactions,
				ControlSum:           controlSum(order.Transactions),
				PaymentTypeInformation: document.PaymentTypeInformation{
					ServiceLevel: document.ServiceLevel{Code: "SEPA"},
				},
				RequestedCollectionDate: order.ExecutionDate.Format("2006-01-02"),
				Creditor: document.Creditor{
					Name: order.Creditor.Name,
					PostalAddress: document.PostalAddress{
						TownName:   order.Creditor.Place,
						Country:    order.Creditor.Country,
						StreetName: order.Creditor.Street,
						PostalCode: order.Creditor.PostalCode,
					},
				},
				CreditorAccount: document.CreditorAccount{
					IBAN: document.IBAN{IBAN: order.Creditor.IBAN},
				},
				CreditorAgent: document.CreditorAgent{
					FinancialInstitutionIdentification: document.FinancialInstitutionIdentification{
						BICFI: order.Creditor.BIC,
					},
				},
				CreditorSchemeIdentification: document.CreditorSchemeIdentification{
					ID: document.ID{
						PrivateIdentification: document.PrivateIdentification{
							Other: document.Other{
								ID: order.Creditor.IBAN,
							},
						},
					},
				},
				DirectDebitTransactionInformation: transactions,
			},
		},
		XMLName: xml.Name{Local: "Document"},
	}, nil
}
