package order

import (
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
