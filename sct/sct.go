// Package sct builds SEPA Credit Transfer messages (pain.001.001.03) from plain payment data. It sets the values that
// SEPA fixes, such as the payment method, the service level, the currency and the charge bearer, and it validates the
// message before it returns the document. Use package ct for documents that this package cannot express.
package sct

import (
	"strings"
	"time"

	"github.com/lance-free/pain-xml/ct"
	"github.com/lance-free/pain-xml/pain"
	"github.com/shopspring/decimal"
)

const (
	serviceLevelSEPA          = "SEPA"
	currencyEUR               = "EUR"
	endToEndIDNotProvided     = "NOTPROVIDED"
	executionDateLayout       = "2006-01-02"
	schemeCompanyRegistration = "COID"
)

// Priority is the instruction priority of the transfers.
type Priority string

const (
	PriorityNormal Priority = "NORM" // PriorityNormal is the normal priority.
	PriorityHigh   Priority = "HIGH" // PriorityHigh is the high priority. Its use must be agreed with the bank.
)

// Message is one credit transfer message: transfers from one debtor account on one execution date.
type Message struct {
	ID              string     // ID identifies the message. It must be unique for the debtor and at most 35 characters.
	PaymentInfoID   string     // PaymentInfoID identifies the payment info block. If empty, ID is used.
	CreatedAt       time.Time  // CreatedAt is the creation time of the message. If zero, the current time is used.
	InitiatingParty string     // InitiatingParty is the name of the party that sends the message. If empty, the name of the debtor is used.
	Debtor          Account    // Debtor is the account the transfers are paid from.
	ExecutionDate   time.Time  // ExecutionDate is the date the bank executes the transfers. If zero, the date of CreatedAt is used.
	BatchBooking    *bool      // BatchBooking requests one booking for all transfers (true) or one booking per transfer (false). If nil, the bank default applies.
	Priority        Priority   // Priority is the instruction priority. If empty, it is not sent.
	UltimateDebtor  *Party     // UltimateDebtor is the party on whose behalf all transfers are paid. Set it here or per transfer, not both.
	Transfers       []Transfer // Transfers are the credit transfers of the message.
}

// Transfer is one credit transfer.
type Transfer struct {
	InstructionID  string          // InstructionID is an optional identification for the debtor's bank.
	EndToEndID     string          // EndToEndID is passed on to the creditor unchanged. If empty, NOTPROVIDED is sent.
	Amount         decimal.Decimal // Amount is the amount in euro, with at most two decimals.
	Creditor       Account         // Creditor is the account the transfer is paid to.
	UltimateDebtor *Party          // UltimateDebtor is the party on whose behalf this transfer is paid.
	Reference      string          // Reference is the creditor reference. A valid Dutch payment reference or RF creditor reference is sent structured.
	Description    string          // Description is the free text for the creditor. It is not sent when Reference is sent structured, because SEPA does not allow both.
}

// Account is a party with its bank account.
type Account struct {
	Name           string          // Name is the name of the account holder, at most 70 characters.
	IBAN           string          // IBAN is the IBAN of the account. Spaces are removed and letters are converted to upper case.
	BIC            string          // BIC is the optional BIC of the bank of the account.
	Address        *Address        // Address is the optional address of the account holder.
	Identification *Identification // Identification optionally identifies the account holder.
}

// Party is a party without a bank account, such as an ultimate debtor.
type Party struct {
	Name           string          // Name is the name of the party, at most 70 characters.
	Address        *Address        // Address is the optional address of the party.
	Identification *Identification // Identification optionally identifies the party.
}

// Address is a postal address. The structured fields and Lines can be combined.
type Address struct {
	Street         string   // Street is the name of the street.
	BuildingNumber string   // BuildingNumber is the number of the building on the street.
	PostCode       string   // PostCode is the postal code.
	Town           string   // Town is the name of the town.
	Country        string   // Country is the ISO 3166 alpha-2 country code (e.g., NL).
	Lines          []string // Lines are at most two free address lines.
}

// Scheme names an identification scheme by an ISO code or by a proprietary name. Set only one of the two.
type Scheme struct {
	Code        string // Code is the ISO code of the scheme (e.g., COID).
	Proprietary string // Proprietary is the proprietary name of the scheme.
}

// Identification identifies an organisation or a person. Create it with OrganisationBIC, OrganisationID,
// CompanyRegistration or PersonID.
type Identification struct {
	bic    string
	other  *pain.GenericIdentification
	person bool
}

// OrganisationBIC identifies an organisation by its BIC or BEI.
func OrganisationBIC(bic string) *Identification {
	return &Identification{bic: strings.ToUpper(bic)}
}

// OrganisationID identifies an organisation by an identification issued under scheme.
func OrganisationID(id string, scheme Scheme, issuer string) *Identification {
	return &Identification{other: genericIdentification(id, scheme, issuer)}
}

// CompanyRegistration identifies an organisation by its company registration number, such as a Dutch Chamber of
// Commerce number with issuer "KVK".
func CompanyRegistration(number, issuer string) *Identification {
	return OrganisationID(number, Scheme{Code: schemeCompanyRegistration}, issuer)
}

// PersonID identifies a person by an identification issued under scheme.
func PersonID(id string, scheme Scheme, issuer string) *Identification {
	return &Identification{other: genericIdentification(id, scheme, issuer), person: true}
}

// Document returns the pain.001.001.03 document of the message, which can be marshalled with encoding/xml or passed
// to the ABN AMRO SDK. It returns the violations of the SEPA rules if the message is not valid.
func (message Message) Document() (ct.Document, error) {
	createdAt := message.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	executionDate := message.ExecutionDate
	if executionDate.IsZero() {
		executionDate = createdAt
	}

	document := ct.NewDocument(ct.CustomerCreditTransferInitiation{
		GroupHeader: pain.GroupHeader{
			MessageID:        message.ID,
			CreationDateTime: pain.CreationDateTime(createdAt),
			InitiatingParty:  pain.InitiatingParty{Name: orDefault(message.InitiatingParty, message.Debtor.Name)},
		},
		PaymentInfo: []ct.PaymentInfo{
			{
				PaymentInfoID:              orDefault(message.PaymentInfoID, message.ID),
				PaymentMethod:              ct.PaymentMethodTransfer,
				BatchBooking:               message.BatchBooking,
				PaymentTypeInfo:            pain.PaymentTypeInfo{InstructionPriority: string(message.Priority), ServiceLevel: pain.ServiceLevel{Code: serviceLevelSEPA}},
				RequestedExecutionDate:     executionDate.Format(executionDateLayout),
				Debtor:                     message.Debtor.party(),
				DebtorAccount:              message.Debtor.account(),
				DebtorAgent:                ct.FinancialInstitution{FinancialInstitutionID: ct.FinancialInstitutionID{BIC: strings.ToUpper(message.Debtor.BIC)}},
				UltimateDebtor:             message.UltimateDebtor.pain(),
				ChargeBearer:               ct.ChargeBearerFollowingServiceLevel,
				CreditTransferTransactions: transactions(message.Transfers),
			},
		},
	})

	if err := document.Validate(); err != nil {
		return ct.Document{}, err
	}

	return document, nil
}

func transactions(transfers []Transfer) []ct.CreditTransferTransactionInfo {
	result := make([]ct.CreditTransferTransactionInfo, 0, len(transfers))
	for _, transfer := range transfers {
		result = append(result, transfer.transaction())
	}

	return result
}

func (transfer Transfer) transaction() ct.CreditTransferTransactionInfo {
	return ct.CreditTransferTransactionInfo{
		PaymentID: pain.PaymentID{
			InstructionID: transfer.InstructionID,
			EndToEndID:    orDefault(transfer.EndToEndID, endToEndIDNotProvided),
		},
		Amount:          ct.Amount{InstigatedAmount: pain.InstigatedAmount{Currency: currencyEUR, Text: transfer.Amount}},
		UltimateDebtor:  transfer.UltimateDebtor.pain(),
		CreditorAgent:   transfer.Creditor.agent(),
		Creditor:        transfer.Creditor.party(),
		CreditorAccount: transfer.Creditor.account(),
		RemittanceInfo:  ct.NewRemittanceInformation(transfer.Reference, transfer.Description),
	}
}

func (account Account) party() pain.Party {
	return pain.Party{Name: account.Name, PostalAddress: account.Address.pain(), Identification: account.Identification.pain()}
}

func (account Account) account() pain.Account {
	iban := strings.ToUpper(strings.ReplaceAll(account.IBAN, " ", ""))
	return pain.Account{ID: pain.AccountID{IBAN: iban}}
}

func (account Account) agent() *ct.FinancialInstitution {
	if account.BIC == "" {
		return nil
	}

	return &ct.FinancialInstitution{FinancialInstitutionID: ct.FinancialInstitutionID{BIC: strings.ToUpper(account.BIC)}}
}

func (party *Party) pain() *pain.Party {
	if party == nil {
		return nil
	}

	return &pain.Party{Name: party.Name, PostalAddress: party.Address.pain(), Identification: party.Identification.pain()}
}

func (address *Address) pain() *pain.PostalAddress {
	if address == nil {
		return nil
	}

	return &pain.PostalAddress{
		StreetName:     address.Street,
		BuildingNumber: address.BuildingNumber,
		PostCode:       address.PostCode,
		TownName:       address.Town,
		Country:        address.Country,
		AddressLine:    address.Lines,
	}
}

func (identification *Identification) pain() *pain.PartyIdentification {
	if identification == nil {
		return nil
	}

	if identification.person {
		return &pain.PartyIdentification{Private: &pain.PersonIdentification{Other: identification.other}}
	}

	return &pain.PartyIdentification{
		Organisation: &pain.OrganisationIdentification{BICOrBEI: identification.bic, Other: identification.other},
	}
}

func genericIdentification(id string, scheme Scheme, issuer string) *pain.GenericIdentification {
	identification := &pain.GenericIdentification{ID: id, Issuer: issuer}
	if scheme != (Scheme{}) {
		identification.SchemeName = &pain.SchemeName{Code: scheme.Code, Proprietary: scheme.Proprietary}
	}

	return identification
}

func orDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}
