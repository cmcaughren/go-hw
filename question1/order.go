package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
