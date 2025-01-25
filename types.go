package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	BankNumber        float64   `json:"bankNumber"`
	EncryptedPassword string    `json:"-"`
	Balance           float64   `json:"balance"`
	CreatedAt         time.Time `json:"createdAt"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

func NewAccount(firstName, lastName, password string) (*Account, error) {
	encPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		BankNumber:        float64(rand.Intn(1000)),
		EncryptedPassword: string(encPw),
		CreatedAt:         time.Now().UTC(),
	}, nil
}

type TransferRequest struct {
	amount int `json:"amount"`
}

type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}
