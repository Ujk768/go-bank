package main

type Account struct {
	ID         int     `json:"id"`
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	BankNumber float64 `json:"bankNumber"`
	Balance    float64 `json:"balance"`
}

func NewAccount(id int, firstName, lastName string, bankNumber, balance float64) *Account {
	return &Account{
		ID:         id,
		FirstName:  firstName,
		LastName:   lastName,
		BankNumber: bankNumber,
		Balance:    balance,
	}
}
