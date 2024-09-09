package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts([]*Account, error)
	GetAccountByID(int) (*Account, error)
}

type PostgresConnection struct {
	db *sql.DB
}

func NewPostgresConnection() (*PostgresConnection, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error handling .env file")
	}

	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	//por defecto se usa 5432 y localhost
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbName, dbPass)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresConnection{db: db}, nil
}

func (s *PostgresConnection) Init() error {
	return s.createAccountTable()
}

func (s *PostgresConnection) createAccountTable() error {
	query := `create table if not exists account(
		id serial primary key ,
		first_name varchar(50),
		last_name varchar(50),
		bank_number serial,
		balance decimal,
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresConnection) CreateAccount(acc *Account) error {
	query := `
	insert into account
	(first_name, last_name, bank_number, balance, created_at)
	values ($1, $2, $3, $4, $5)`
	resp, err := s.db.Query(query, acc.FirstName, acc.LastName, acc.BankNumber, acc.Balance, acc.CreatedAt)

	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp)

	return nil
}

func (s *PostgresConnection) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresConnection) DeleteAccount(id int) error {
	_, err := s.db.Query("DELETE FROM account WHERE id = $1", id)
	return err
}

func (s *PostgresConnection) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query("select * from account where id = $1", id)
	if err != nil {
		return nil, err
	}
	
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)
}
func (s *PostgresConnection) GetAccounts() ([]*Account, error) {
	query := `select * from account`

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		account := new(Account)
		account, err := scanIntoAccount(rows) 
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}


func scanIntoAccount(rows *sql.Rows) (*Account, error){
	
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.BankNumber,
		&account.Balance,
		&account.CreatedAt,
	)

	return account, err
}
