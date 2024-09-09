package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func createJWT(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 1500,
		"accountNumber": account.BankNumber,
	}
	
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func withJWTAuth(handleFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		fmt.Println("calling JWT auth middleware")
		
		tokenString := r.Header.Get("x-jwt-token")

		token , err := validateJWT(tokenString)
		if err != nil {
			WriteJSON(w, http.StatusForbidden, ApiError{Error : "invalid token"})
			return
		}
		if !token.Valid {
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "invalid toekn"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		fmt.Println(claims)

		handleFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	
	secret := os.Getenv("JWT_SECRET")
	
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

type ApiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHttpHandleFunc(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr   string
	dbConnection *PostgresConnection
}

func NewAPIServer(listenAddr string, conn *PostgresConnection) *APIServer {
	return &APIServer{
		listenAddr:   listenAddr,
		dbConnection: conn,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHttpHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHttpHandleFunc(s.handleGetAccountByID)))
	router.HandleFunc("/transfer", makeHttpHandleFunc(s.handleTransferAccount))

	log.Println("JSON API server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {

	switch r.Method {
	case "GET":
		return s.handleGetAccount(w)
	case "POST":
		return s.handleCreateAccount(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	default:
		return fmt.Errorf("method not allowed %s", r.Method)
	}
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter) error {
	accounts, err := s.dbConnection.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}
		account, err := s.dbConnection.GetAccountByID(id)
		if err != nil {
			return err
		}
		return WriteJSON(w, http.StatusOK, account)
	} else if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	} else {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)

	if err := s.dbConnection.CreateAccount(account); err != nil {
		return err
	}

	tokenString, err := createJWT(account)
	if err != nil {
		return err
	}

	fmt.Println("JWT token: ", tokenString)
	return WriteJSON(w, http.StatusOK, account)

}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.dbConnection.DeleteAccount(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleTransferAccount(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()
	return WriteJSON(w, http.StatusOK, transferReq)
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}

