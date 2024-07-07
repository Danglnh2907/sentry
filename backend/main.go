package main

import (
	//Import standard library
	"fmt"
	"log"
	"net/http"

	//Import standard package
	"sentry/auth"
	"sentry/crud"
)

func main() {
	//Create multiplexer for routing
	mux := http.NewServeMux()

	//Handle auth
	mux.HandleFunc("/create-account", auth.CreateAccount)
	mux.HandleFunc("/login", auth.Login)
	mux.HandleFunc("/get-profile", auth.GetProfile)
	mux.HandleFunc("/delete-profile", auth.DeleteProfile)

	//Handle crud
	mux.HandleFunc("/transaction", crud.PostTransaction)
	mux.HandleFunc("/transactions", crud.PostTransactions)
	mux.HandleFunc("/get-transactions", crud.GetTransactions)

	//Run server at port 8080
	fmt.Println("Sentry running at http://localhost:8080")
	err := http.ListenAndServe("localhost:8080", mux)

	if err != nil {
		var w http.ResponseWriter
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Fatal(err)
	}
}
