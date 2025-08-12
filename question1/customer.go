package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/microsoft/go-mssqldb"
)

type Customer struct {
	CustomerID    int
	FirstName     string
	LastName      string
	DOB           time.Time
	OrderTotal    float64
	OrderTaxTotal float64
}

// w Response Writer used to write back to the client
// r Request contains the request information
func GetCustomers(w http.ResponseWriter, r *http.Request) {
	query := "SELECT * FROM Customer"
	rows, err := db.Query(query)
	if err != nil {
		//send 500 internal error to client
		fmt.Printf("Database error: %v", err)
		http.Error(w, "Failed to fetch Customers", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var customer Customer
		err := rows.Scan(&customer.CustomerID,
			&customer.FirstName,
			&customer.LastName,
			&customer.DOB,
			&customer.OrderTotal,
			&customer.OrderTaxTotal)
		if err != nil {
			fmt.Printf("Scan error: %v", err)
			continue
		}
		customers = append(customers, customer)
	}

	//Send JSON reponse
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

// Get a single customer, by CustomerID
func GetCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["CustomerID"]

	query := "SELECT * FROM Customer WHERE CustomerID = @p1"
	var customer Customer
	err := db.QueryRow(query, id).Scan(&customer.CustomerID,
		&customer.FirstName,
		&customer.LastName,
		&customer.DOB,
		&customer.OrderTotal,
		&customer.OrderTaxTotal)
	if err == sql.ErrNoRows {
		http.Error(w, "Customer Not Found", http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Printf("Scan error: %v", err)
		http.Error(w, "Failed to fetch customer", http.StatusInternalServerError)
		return
	}

	//Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

// Add a new customer
func AddCustomer(w http.ResponseWriter, r *http.Request) {
	var c Customer
	json.NewDecoder(r.Body).Decode(&c)
	valid, errMsg := validateCustomer(c)
	if !valid {
		fmt.Print("Bad customer data: &v\n", errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	//"OUTPUT" with QueryRow will get us a return value - use of Exec to insert would give none
	query := `INSERT INTO Customer (FirstName, LastName, DOB, OrderTotal, OrderTaxTotal) 
	OUTPUT INSERTED.CustomerID, INSERTED.FirstName, INSERTED.LastName, INSERTED.DOB, INSERTED.OrderTotal, INSERTED.OrderTaxTotal
	VALUES (@p1, @p2, @p3, @p4, @p5)`
	var newCustomer Customer

	err := db.QueryRow(query,
		c.FirstName,
		c.LastName,
		c.DOB,
		c.OrderTotal,
		c.OrderTaxTotal).Scan(
		&newCustomer.CustomerID,
		&newCustomer.FirstName,
		&newCustomer.LastName,
		&newCustomer.DOB,
		&newCustomer.OrderTotal,
		&newCustomer.OrderTaxTotal)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		http.Error(w, "Error adding customer", http.StatusInternalServerError)
		return
	}
	//Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCustomer)
}

func validateCustomer(customer Customer) (bool, string) {
	// Validate FirstName
	if len(customer.FirstName) == 0 {
		return false, "FirstName is required"
	}
	if len(customer.FirstName) > 50 {
		return false, "FirstName must be 50 characters or less"
	}

	// Validate LastName
	if len(customer.LastName) == 0 {
		return false, "LastName is required"
	}
	if len(customer.LastName) > 50 {
		return false, "LastName must be 50 characters or less"
	}

	// Validate DOB
	if customer.DOB.IsZero() {
		return false, "DOB is required"
	}
	if customer.DOB.After(time.Now()) {
		return false, "DOB cannot be in the future"
	}

	// Validate OrderTotal
	if customer.OrderTotal < 0 {
		return false, "OrderTotal cannot be negative"
	}

	// Validate OrderTaxTotal
	if customer.OrderTaxTotal < 0 {
		return false, "OrderTaxTotal cannot be negative"
	}
	expectedTax := customer.OrderTotal * 0.1
	if customer.OrderTaxTotal != expectedTax {
		return false, fmt.Sprintf("OrderTaxTotal must be 10%% of OrderTotal (expected: %.2f)", expectedTax)
	}

	return true, ""
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["CustomerID"]

	var c Customer
	json.NewDecoder(r.Body).Decode(&c)
	valid, errMsg := validateCustomer(c)
	if !valid {
		fmt.Printf("Bad value: %v\n", errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	query := `
		UPDATE Customer
		SET FirstName = @p1,
		LastName = @p2,
		DOB = @p3,
		OrderTotal = @p4,
		OrderTaxTotal = @p5
		OUTPUT INSERTED.CustomerId, 
		INSERTED.FirstName,
		INSERTED.LastName, 
		INSERTED.DOB, 
		INSERTED.OrderTotal, 
		INSERTED.OrderTaxTotal
		WHERE CustomerID = @p6`

	var updatedCustomer Customer
	err := db.QueryRow(query,
		c.FirstName,
		c.LastName,
		c.DOB,
		c.OrderTotal,
		c.OrderTaxTotal,
		id).Scan(&updatedCustomer.CustomerID,
		&updatedCustomer.FirstName,
		&updatedCustomer.LastName,
		&updatedCustomer.DOB,
		&updatedCustomer.OrderTotal,
		&updatedCustomer.OrderTaxTotal)
	if err == sql.ErrNoRows {
		fmt.Printf("Customer does not exist: %v\n", err)
		http.Error(w, "Customer does not exist", http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Printf("Error updating customer: %v\n", err)
		http.Error(w, "Error updating customer", http.StatusInternalServerError)
		return
	}

	//Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCustomer)
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["CustomerID"]

	query := `DELETE FROM Customer 
		OUTPUT DELETED.CustomerID
		WHERE CustomerId = @p1`

	var deletedID int
	err := db.QueryRow(query, id).Scan(&deletedID)
	if err == sql.ErrNoRows {
		fmt.Printf("Customer does not exist: %v\n", err)
		http.Error(w, "Customer does not exist", http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Printf("Error deleting customer: %v\n", err)
		http.Error(w, "Error deleting customer", http.StatusInternalServerError)
		return
	}

	//Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedID)
}
