package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountById(int) (*Account, error)
}

type PostGresStorage struct {
	db *sql.DB
}

func NewPostGresStorage() (*PostGresStorage, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &PostGresStorage{db: db}, nil
}

func (s *PostGresStorage) CreateAccount(a *Account) error {
	return nil
}

func (s *PostGresStorage) UpdateAccount(a *Account) error {
	return nil
}

func (s *PostGresStorage) DeleteAccount(id int) error {
	return nil
}
func (s *PostGresStorage) GetAccountById(id int) (*Account, error) {
	return nil, nil
}
