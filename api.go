package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	jwt "github.com/golang-jwt/jwt/v5"
)

type ApiServer struct {
	listenAddr string
	store      Storage
}

// Constructor for ApiServer
func NewAPIServer(listenAddr string, store Storage) *ApiServer {
	fmt.Println("Inside NewAPIServer")
	return &ApiServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

// Define the type of the function that will be used to handle the API requests
type ApiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

// Create JSON response to set the status code and content type
func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

// Make this function to check error and use this func as http handler
func makeHttpHandleFunc(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc(("/login"), makeHttpHandleFunc(s.handleLogin))
	router.HandleFunc("/account", makeHttpHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHttpHandleFunc(s.handleGetAccountById), s.store))
	router.HandleFunc("/transfer", makeHttpHandleFunc(s.handleTransfer))
	http.ListenAndServe(s.listenAddr, router)
}

func (s *ApiServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("Method Not Allowed %s", r.Method)
	}
	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return err
	}
	
	writeJSON(w, http.StatusOK, request)
	return nil
}

// Define the functions to handle the API requests
func (s *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	// handle  account
	if r.Method == http.MethodGet {
		return s.handleGetAccount(w, r)
	}
	if r.Method == http.MethodPost {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == http.MethodDelete {
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("Method Not Allowed %s", r.Method)
}

func (s *ApiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, accounts)
}

func (s *ApiServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		id, err := getId(r)
		if err != nil {
			return err
		}
		account, err := s.store.GetAccountById(id)
		if err != nil {
			return err
		}
		return writeJSON(w, http.StatusOK, account)
	}
	if r.Method == http.MethodDelete {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("Method Not Allowed %s", r.Method)
}

// handle create account
func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	account, err := NewAccount(req.FirstName, req.LastName, req.Password)
	if err != nil {
		return err
	}

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, account)

}

// handle delete account
func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)
	if err != nil {
		return err
	}
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *ApiServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	// handle transfer
	trnansferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(trnansferReq); err != nil {
		return err
	}
	defer r.Body.Close()

	return writeJSON(w, http.StatusOK, trnansferReq)
}

func permissionDenied(w http.ResponseWriter) {
	writeJSON(w, http.StatusUnauthorized, ApiError{Error: "Access Denied"})
}

func withJWTAuth(handleFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Calling JWT middleware")
		tokenString := r.Header.Get("x-jwt-token")
		token, err := ValidateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}

		userId, err := getId(r)
		if err != nil {
			permissionDenied(w)
			return
		}
		account, err := s.GetAccountById(userId)
		if err != nil {
			writeJSON(w, http.StatusForbidden, ApiError{Error: "Invalid Token"})
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		fmt.Println("Claims: ", claims["accountNumber"])
		fmt.Println("Account: ", account.BankNumber)

		if account.BankNumber != claims["accountNumber"].(float64) {
			permissionDenied(w)
			return
		}
		handleFunc(w, r)
	}
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected SIgning Mehtod %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

}

func getId(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("Invalid ID %s", idStr)
	}
	return id, nil
}

func createJWT(account *Account) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	// Create the Claims
	claims := &jwt.MapClaims{
		"expiresAt":     150000,
		"accountNumber": account.BankNumber,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))

}
