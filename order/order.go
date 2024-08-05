package order

import (
	"github.com/shopspring/decimal"
	"time"
)

// Party represents information about a party in a transaction.
type Party struct {
	// Name is a field that represents the name of the person involved in the transaction.
	Name string

	// Street represents the street of the address associated with a financial transaction.
	Street string

	// PostalCode is a field that represents the postal code of a person's address.
	PostalCode string

	// Place represents the name of a place associated with a transaction or debtor.
	Place string

	// Country represents the country of a financial transaction.
	// The length of the country code is limited to two characters.
	Country string

	// IBAN represents an International Bank Account Number used for financial transactions.
	IBAN string

	// BIC represents the Bank Identifier Code for a financial transaction.
	// It is a unique identification code for a specific bank.
	BIC string
}

// Transaction represents a financial transaction with its details.
type Transaction struct {
	// Creditor represents information about the party who is owed the money in a transaction.
	Creditor Party

	// Currency represents the currency of a financial transaction.
	Currency string

	// Amount is a field that represents the amount of money involved in the transaction.
	Amount decimal.Decimal
}

// Order represents an order with its execution date, transactions, and debtor information.
// An order is a financial document that includes information about when it should be executed,
// the list of transactions associated with it, and details about the debtor involved in the order.
type Order struct {
	// ExecutionDate represents the date when an order should be executed.
	ExecutionDate time.Time

	// Transactions represents a slice of Transaction objects.
	Transactions []Transaction

	// Debtor represents information about the debtor involved in a financial transaction.
	Debtor Party
}
