package main

import (
	"encoding/json"
	"fmt"
	"net/http"

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
