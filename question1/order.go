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
	OrderFulfilledDate time.Time
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
