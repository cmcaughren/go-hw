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

type Order struct {
	OrderID            int
	OrderNumber        string
	CustomerID         int
	OrderCreateDate    time.Time
	OrderFulfilledDate *time.Time //pointer so we can set this to null if needed
	OrderTotal         float64
	OrderTaxTotal      float64
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	query := "SELECT * FROM [Order]"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("Error getting Orders: %v\n", err)
		http.Error(w, "Error getting Orders", http.StatusInternalServerError)
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.OrderID,
			&order.OrderNumber,
			&order.CustomerID,
			&order.OrderCreateDate,
			&order.OrderFulfilledDate,
			&order.OrderTotal,
			&order.OrderTaxTotal)
		if err != nil {
			fmt.Printf("Scan error: %v\n", err)
			continue
		}
		orders = append(orders, order)
	}

	//Send JSON reponse
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// Get a single order, by OrderID
func GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["OrderID"]

	query := "SELECT * FROM [Order] WHERE OrderID = @p1"
	var order Order
	err := db.QueryRow(query, id).Scan(&order.OrderID,
		&order.OrderNumber,
		&order.CustomerID,
		&order.OrderCreateDate,
		&order.OrderFulfilledDate,
		&order.OrderTotal,
		&order.OrderTaxTotal)
	if err == sql.ErrNoRows {
		fmt.Printf("Order Not Found: %v\n", err)
		http.Error(w, "Order Not Found", http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Printf("Scan error: %v\n", err)
		http.Error(w, "Failed to fetch order", http.StatusInternalServerError)
		return
	}

	//Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// Add a new order
func AddOrder(w http.ResponseWriter, r *http.Request) {
	var o Order
	json.NewDecoder(r.Body).Decode(&o)
	valid, errMsg := validateOrder(o)
	if !valid {
		fmt.Printf("Bad order data: %v\n", errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	//"OUTPUT" with QueryRow will get us a return value - use of Exec to insert would give none
	query := `INSERT INTO [Order] (OrderNumber, CustomerID, OrderCreateDate, OrderFulfilledDate, OrderTotal, OrderTaxTotal) 
	OUTPUT INSERTED.OrderID, INSERTED.OrderNumber, INSERTED.CustomerID, INSERTED.OrderCreateDate, INSERTED.OrderFulfilledDate, INSERTED.OrderTotal, INSERTED.OrderTaxTotal
	VALUES (@p1, @p2, @p3, @p4, @p5, @p6)`
	var newOrder Order

	err := db.QueryRow(query,
		o.OrderNumber,
		o.CustomerID,
		o.OrderCreateDate,
		o.OrderFulfilledDate,
		o.OrderTotal,
		o.OrderTaxTotal).Scan(
		&newOrder.OrderID,
		&newOrder.OrderNumber,
		&newOrder.CustomerID,
		&newOrder.OrderCreateDate,
		&newOrder.OrderFulfilledDate,
		&newOrder.OrderTotal,
		&newOrder.OrderTaxTotal)
	if err != nil {
		fmt.Printf("Error adding order: %v\n", err)
		http.Error(w, "Error adding order", http.StatusInternalServerError)
		return
	}
	//Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newOrder)
}

func validateOrder(order Order) (bool, string) {
	// Validate OrderNumber
	if len(order.OrderNumber) == 0 {
		return false, "OrderNumber is required"
	}
	if len(order.OrderNumber) > 50 {
		return false, "OrderNumber must be 50 characters or less"
	}

	// Validate CustomerID
	if order.CustomerID <= 0 {
		return false, "CustomerID is required and must be positive"
	}

	// Validate OrderCreateDate
	if order.OrderCreateDate.IsZero() {
		return false, "OrderCreateDate is required"
	}
	if order.OrderCreateDate.After(time.Now()) {
		return false, "OrderCreateDate cannot be in the future"
	}

	// Validate OrderFulfilledDate (if provided, must be after create date)
	if order.OrderFulfilledDate != nil && order.OrderFulfilledDate.Before(order.OrderCreateDate) {
		return false, "OrderFulfilledDate must be after OrderCreateDate"
	}

	// Validate OrderTotal
	if order.OrderTotal < 0 {
		return false, "OrderTotal cannot be negative"
	}

	// Validate OrderTaxTotal
	if order.OrderTaxTotal < 0 {
		return false, "OrderTaxTotal cannot be negative"
	}

	return true, ""
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["OrderID"]

	var o Order
	json.NewDecoder(r.Body).Decode(&o)
	valid, errMsg := validateOrder(o)
	if !valid {
		fmt.Printf("Bad value: %v\n", errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	query := `
		UPDATE [Order]
		SET OrderNumber = @p1,
		CustomerID = @p2,
		OrderCreateDate = @p3,
		OrderFulfilledDate = @p4,
		OrderTotal = @p5,
		OrderTaxTotal = @p6
		OUTPUT INSERTED.OrderID,
		INSERTED.OrderNumber, 
		INSERTED.CustomerID,
		INSERTED.OrderCreateDate, 
		INSERTED.OrderFulfilledDate, 
		INSERTED.OrderTotal, 
		INSERTED.OrderTaxTotal
		WHERE OrderID = @p7`

	var updatedOrder Order
	err := db.QueryRow(query,
		o.OrderNumber,
		o.CustomerID,
		o.OrderCreateDate,
		o.OrderFulfilledDate,
		o.OrderTotal,
		o.OrderTaxTotal,
		id).Scan(&updatedOrder.OrderID,
		&updatedOrder.OrderNumber,
		&updatedOrder.CustomerID,
		&updatedOrder.OrderCreateDate,
		&updatedOrder.OrderFulfilledDate,
		&updatedOrder.OrderTotal,
		&updatedOrder.OrderTaxTotal)
	if err == sql.ErrNoRows {
		fmt.Printf("Order does not exist: %v\n", err)
		http.Error(w, "Order does not exist", http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Printf("Error updating order: %v\n", err)
		http.Error(w, "Error updating order", http.StatusInternalServerError)
		return
	}

	//Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedOrder)
}
