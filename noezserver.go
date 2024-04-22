package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Customer struct {
	Name string
	Age  int
}

var customers []Customer

func main() {
	http.HandleFunc("/customers", customersHandler)
	http.HandleFunc("/health", healthCheckHandler)

	log.Println("server start listnening on port 8080")
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func customersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getCustomers(w, r)
	case http.MethodPost:
		postCustomers(w, r)
	default:
		http.Error(w, "PIZDA", http.StatusMethodNotAllowed)
	}
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	//json.NewEncoder(w).Encode(customers)
	fmt.Fprintf(w, "get customers: '%v'", customers)
}

func postCustomers(w http.ResponseWriter, r *http.Request) {
	var customer Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "http web-server works")
}
