package main

import (
	"math/rand"
	"time"
)

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Account struct {
	ID         uint      `json:"id"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	BankNumber int64     `json:"bankNumber"`
	Balance    float64   `json:"balance"`
	CreatedAt  time.Time `json:"created_at"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		FirstName:  firstName,
		LastName:   lastName,
		BankNumber: generateBankNumber(),
		CreatedAt:  time.Now().UTC(),
	}
}

func generateBankNumber() int64 {
	//podria agregar logica para que sea unico mas adelante

	return int64(rand.Intn(1000000))
}
