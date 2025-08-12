package main

import (
	"encoding/json"
	"fmt"
	"net/http"

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
