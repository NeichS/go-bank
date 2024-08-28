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
		id primary key serial,
		first_name varchar(50),
		last_name varchar(50),
		bank_number serial,
		balance decimal,
		created_at timestamp,
	)`

	_, err := s.db.Exec(query)
	return err
}


func (s *PostgresConnection) CreateAccount(*Account) error {
	return nil
}

func (s *PostgresConnection) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresConnection) DeleteAccount(id uint) error {
	return nil
}

func (s *PostgresConnection) GetAccountByID(id uint) (*Account, error) {
	return nil, nil
}