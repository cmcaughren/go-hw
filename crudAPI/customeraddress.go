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

//TODO Add get customer address by customerID

// Add a new customer address
// allows customer to have multiple addresses!
func AddCustomerAddress(w http.ResponseWriter, r *http.Request) {
	var ca CustomerAddress
	json.NewDecoder(r.Body).Decode(&ca)
	valid, errMsg := validateCustomerAddress(ca)
	if !valid {
		fmt.Print("Bad customer data: &v\n", errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	//"OUTPUT" with QueryRow will get us a return value - use of Exec to insert would give none
	query := `INSERT INTO CustomerAddress (CustomerID, Line1, Line2, City, StateProvince, Country) 
	OUTPUT INSERTED.CustomerID, INSERTED.Line1, INSERTED.Line2, INSERTED.City, INSERTED.Stateprovince, INSERTED.Country
	VALUES (@p1, @p2, @p3, @p4, @p5, @p6)`
	var newCustomerAddress CustomerAddress

	err := db.QueryRow(query,
		ca.CustomerID,
		ca.Line1,
		ca.Line2,
		ca.City,
		ca.StateProvince,
		ca.Country).Scan(
		&newCustomerAddress.CustomerID,
		&newCustomerAddress.Line1,
		&newCustomerAddress.Line2,
		&newCustomerAddress.City,
		&newCustomerAddress.StateProvince,
		&newCustomerAddress.Country)
	if err != nil {
		fmt.Printf("Error adding customer address: %v\n", err)
		http.Error(w, "Error adding customer address", http.StatusInternalServerError)
		return
	}
	//Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newCustomerAddress)
}

func validateCustomerAddress(customerAddress CustomerAddress) (bool, string) {
	// Validate CustomerID
	if customerAddress.CustomerID <= 0 {
		return false, "CustomerID is required and must be positive"
	}

	// Validate Line1
	if len(customerAddress.Line1) == 0 {
		return false, "Line1 is required"
	}
	if len(customerAddress.Line1) > 100 {
		return false, "Line1 must be 100 characters or less"
	}

	// Validate Line2 (optional but has max length)
	if len(customerAddress.Line2) > 100 {
		return false, "Line2 must be 100 characters or less"
	}

	// Validate City
	if len(customerAddress.City) == 0 {
		return false, "City is required"
	}
	if len(customerAddress.City) > 50 {
		return false, "City must be 50 characters or less"
	}

	// Validate StateProvince
	if len(customerAddress.StateProvince) == 0 {
		return false, "StateProvince is required"
	}
	if len(customerAddress.StateProvince) > 50 {
		return false, "StateProvince must be 50 characters or less"
	}

	// Validate Country
	if len(customerAddress.Country) == 0 {
		return false, "Country is required"
	}
	if len(customerAddress.Country) > 50 {
		return false, "Country must be 50 characters or less"
	}

	return true, ""
}

func UpdateCustomerAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["CustomerAddressID"]

	var ca CustomerAddress
	json.NewDecoder(r.Body).Decode(&ca)
	valid, errMsg := validateCustomerAddress(ca)
	if !valid {
		fmt.Printf("Bad value: %v\n", errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	query := `
		UPDATE CustomerAddress
		SET CustomerID = @p1,
		Line1 = @p2,
		Line2 = @p3,
		City = @p4,
		StateProvince = @p5,
		Country = @p6
		OUTPUT INSERTED.CustomerAddressID,
		INSERTED.CustomerID, 
		INSERTED.Line1,
		INSERTED.Line2, 
		INSERTED.City, 
		INSERTED.StateProvince, 
		INSERTED.Country
		WHERE CustomerAddressID = @p7`

	var updatedCustomerAddress CustomerAddress
	err := db.QueryRow(query,
		ca.CustomerID,
		ca.Line1,
		ca.Line2,
		ca.City,
		ca.StateProvince,
		ca.Country,
		id).Scan(&updatedCustomerAddress.CustomerAddressID,
		&updatedCustomerAddress.CustomerID,
		&updatedCustomerAddress.Line1,
		&updatedCustomerAddress.Line2,
		&updatedCustomerAddress.City,
		&updatedCustomerAddress.StateProvince,
		&updatedCustomerAddress.Country)
	if err == sql.ErrNoRows {
		fmt.Printf("Customer Address does not exist: %v\n", err)
		http.Error(w, "Customer Address does not exist", http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Printf("Error updating customer address: %v\n", err)
		http.Error(w, "Error updating customer address", http.StatusInternalServerError)
		return
	}

	//Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCustomerAddress)
}

func DeleteCustomerAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["CustomerAddressID"]

	query := `DELETE FROM CustomerAddress 
		OUTPUT DELETED.CustomerAddressID
		WHERE CustomerAddressID = @p1`

	var deletedID int
	err := db.QueryRow(query, id).Scan(&deletedID)
	if err == sql.ErrNoRows {
		fmt.Printf("Customer Address does not exist: %v\n", err)
		http.Error(w, "Customer Address does not exist", http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Printf("Error deleting customer address: %v\n", err)
		http.Error(w, "Error deleting customer address", http.StatusInternalServerError)
		return
	}

	//Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedID)
}
