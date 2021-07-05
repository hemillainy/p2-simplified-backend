package main

import (
	"github.com/gorilla/mux"
	"github.com/hemillainy/backend/transfer"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/transaction", transfer.TransactionHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
