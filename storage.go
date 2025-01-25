package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountById(int) (*Account, error)
	GetAccountByNumber(int) (*Account, error)
}

type PostGresStorage struct {
	db *sql.DB
}

func (s *PostGresStorage) Init() error {
	return s.createAccountTable()
}

func (s *PostGresStorage) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts (
	id SERIAL PRIMARY KEY,
	first_name varchar(255),
	last_name varchar(255) ,
	number serial,
	balance serial,
	created_at timestamp  
	)`
	_, err := s.db.Exec(query)
	return err
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
	query := `INSERT INTO accounts (first_name, last_name, number, balance, created_at) VALUES ($1, $2, $3, $4, $5)`
	res, err := s.db.Query(query, a.FirstName, a.LastName, a.BankNumber, a.Balance, a.CreatedAt)
	if err != nil {
		return err
	}
	fmt.Printf("Result: %v\n", res)
	return nil
}

func (s *PostGresStorage) UpdateAccount(a *Account) error {
	return nil
}

func (s *PostGresStorage) DeleteAccount(id int) error {
	_, err := s.db.Query("DELETE FROM accounts WHERE id = $1", id)

	return err
}

func (s *PostGresStorage) GetAccountByNumber(number int) (*Account, error) {
	rows, err := s.db.Query("SELECT * FROM accounts WHERE number = $1", number)
	if err != nil {
		return nil, err

	}

	for rows.Next() {
		return ScanIntoAccount(rows)
	}

	return nil, fmt.Errorf("Account %d not found ", number)

}

func (s *PostGresStorage) GetAccountById(id int) (*Account, error) {
	rows, err := s.db.Query("SELECT * FROM accounts WHERE id = $1", id)
	if err != nil {
		return nil, err

	}

	for rows.Next() {
		return ScanIntoAccount(rows)
	}

	return nil, fmt.Errorf("Account %d not found ", id)

}

func (s *PostGresStorage) GetAccounts() ([]*Account, error) {
	query := `SELECT * FROM accounts`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {

		account, err := ScanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)

	}
	return accounts, nil

}

func ScanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.BankNumber, &account.Balance, &account.CreatedAt)
	if err != nil {
		return nil, err
	}
	return account, nil
}
