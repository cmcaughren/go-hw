package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/microsoft/go-mssqldb"
)

type CustomerAddress struct {
	CustomerAddressID int
	CustomerID        int
	Line1             string
	Line2             string
	City              string
	StateProvince     string
	Country           string
}

func GetCustomerAddresses(w http.ResponseWriter, r *http.Request) {
	query := "SELECT * FROM CustomerAddress"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("Error getting Customer Addresses: %v\n", err)
		http.Error(w, "Error getting Customer Addresses", http.StatusInternalServerError)
	}
	defer rows.Close()

	var customerAddresses []CustomerAddress
	for rows.Next() {
		var customerAddress CustomerAddress
		err := rows.Scan(&customerAddress.CustomerAddressID,
			&customerAddress.CustomerID,
			&customerAddress.Line1,
			&customerAddress.Line2,
			&customerAddress.City,
			&customerAddress.StateProvince,
			&customerAddress.Country)
		if err != nil {
			fmt.Printf("Scan error: %v\n", err)
			continue
		}
		customerAddresses = append(customerAddresses, customerAddress)
	}

	//Send JSON reponse
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customerAddresses)
}

// Get a single customer address, by CustomerAddressID
func GetCustomerAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["CustomerAddressID"]

	query := "SELECT * FROM CustomerAddress WHERE CustomerAddressID = @p1"
	var customerAddress CustomerAddress
	err := db.QueryRow(query, id).Scan(&customerAddress.CustomerAddressID,
		&customerAddress.CustomerID,
		&customerAddress.Line1,
		&customerAddress.Line2,
		&customerAddress.City,
		&customerAddress.StateProvince,
		&customerAddress.Country)
	if err == sql.ErrNoRows {
		fmt.Printf("CustomerAddress Not Found: %v\n", err)
		http.Error(w, "CustomerAddress Not Found", http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Printf("Scan error: %v\n", err)
		http.Error(w, "Failed to fetch customer", http.StatusInternalServerError)
		return
	}

	//Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customerAddress)
}
