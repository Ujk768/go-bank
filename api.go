package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	listenAddr string
	store Storage
}

// Constructor for ApiServer
func NewAPIServer(listenAddr string, store Storage) *ApiServer {
	fmt.Println("Inside NewAPIServer")
	return &ApiServer{
		listenAddr: listenAddr,
		store: store,
	}
}

// Define the type of the function that will be used to handle the API requests
type ApiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
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
	router.HandleFunc("/account", makeHttpHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHttpHandleFunc(s.handleAccount))
	log.Println("APi Server started...")
	http.ListenAndServe(s.listenAddr, router)
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
	// handle get account
	id := mux.Vars(r)["id"]
	fmt.Println("ID: ", id)
	account := NewAccount(1, "John", "Doe", 123456, 1000)
	return writeJSON(w, http.StatusOK, account)
}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	// handle get account
	return nil
}
func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	// handle delete account
	return nil
}

func (s *ApiServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	// handle transfer
	return nil
}
