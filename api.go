package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"

)

type APIServer struct {
	listenAddr string
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func WriteJSON(w http.ResponseWriter, status int, resp any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	//create a new encoder to convert Go data structures (like maps, structs, slices, etc.) to JSON format and then encode response
	return json.NewEncoder(w).Encode(resp)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}


func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/", makeHTTPHandleFunc(s.handleAccount))
	// router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccount))
	log.Println("JSON API server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router) // starts http server on on address specified and listens for incoming requests
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}
