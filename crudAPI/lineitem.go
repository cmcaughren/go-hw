package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/microsoft/go-mssqldb"
)

type LineItem struct {
	LineItemID        int
	OrderID           int
	ProductID         int
	LineItemUnitPrice float64
	Quantity          int
	LineItemDiscount  float64
}

func GetLineItems(w http.ResponseWriter, r *http.Request) {
	query := "SELECT * FROM LineItem"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("Error getting Line Items: %v\n", err)
		http.Error(w, "Error getting Line Items", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var lineItems []LineItem
	for rows.Next() {
		var lineItem LineItem
		err := rows.Scan(&lineItem.LineItemID,
			&lineItem.OrderID,
			&lineItem.ProductID,
			&lineItem.LineItemUnitPrice,
			&lineItem.Quantity,
			&lineItem.LineItemDiscount)
		if err != nil {
			fmt.Printf("Scan error: %v\n", err)
			continue
		}
		lineItems = append(lineItems, lineItem)
	}

	//Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lineItems)
}

// Get a single line item, by LineItemID
func GetLineItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["LineItemID"]

	query := "SELECT * FROM LineItem WHERE LineItemID = @p1"
	var lineItem LineItem
	err := db.QueryRow(query, id).Scan(&lineItem.LineItemID,
		&lineItem.OrderID,
		&lineItem.ProductID,
		&lineItem.LineItemUnitPrice,
		&lineItem.Quantity,
		&lineItem.LineItemDiscount)
	if err == sql.ErrNoRows {
		fmt.Printf("LineItem Not Found: %v\n", err)
		http.Error(w, "LineItem Not Found", http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Printf("Scan error: %v\n", err)
		http.Error(w, "Failed to fetch line item", http.StatusInternalServerError)
		return
	}

	//Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lineItem)
}
