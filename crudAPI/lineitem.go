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

// Add a new line item
func AddLineItem(w http.ResponseWriter, r *http.Request) {
	var li LineItem
	json.NewDecoder(r.Body).Decode(&li)
	valid, errMsg := validateLineItem(li)
	if !valid {
		fmt.Printf("Bad line item data: %v\n", errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	//"OUTPUT" with QueryRow will get us a return value - use of Exec to insert would give none
	query := `INSERT INTO LineItem (OrderID, ProductID, LineItemUnitPrice, Quantity, LineItemDiscount) 
	OUTPUT INSERTED.LineItemID, INSERTED.OrderID, INSERTED.ProductID, INSERTED.LineItemUnitPrice, INSERTED.Quantity, INSERTED.LineItemDiscount
	VALUES (@p1, @p2, @p3, @p4, @p5)`
	var newLineItem LineItem

	err := db.QueryRow(query,
		li.OrderID,
		li.ProductID,
		li.LineItemUnitPrice,
		li.Quantity,
		li.LineItemDiscount).Scan(
		&newLineItem.LineItemID,
		&newLineItem.OrderID,
		&newLineItem.ProductID,
		&newLineItem.LineItemUnitPrice,
		&newLineItem.Quantity,
		&newLineItem.LineItemDiscount)
	if err != nil {
		fmt.Printf("Error adding line item: %v\n", err)
		http.Error(w, "Error adding line item", http.StatusInternalServerError)
		return
	}
	//Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newLineItem)
}

func validateLineItem(lineItem LineItem) (bool, string) {
	// Validate OrderID
	if lineItem.OrderID <= 0 {
		return false, "OrderID is required and must be positive"
	}

	// Validate ProductID
	if lineItem.ProductID <= 0 {
		return false, "ProductID is required and must be positive"
	}

	// Validate LineItemUnitPrice
	if lineItem.LineItemUnitPrice < 0 {
		return false, "LineItemUnitPrice cannot be negative"
	}

	// Validate Quantity
	if lineItem.Quantity <= 0 {
		return false, "Quantity must be positive"
	}

	// Validate LineItemDiscount
	if lineItem.LineItemDiscount < 0 {
		return false, "LineItemDiscount cannot be negative"
	}
	if lineItem.LineItemDiscount > lineItem.LineItemUnitPrice * float64(lineItem.Quantity) {
		return false, "LineItemDiscount cannot exceed total price"
	}

	return true, ""
}

func UpdateLineItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["LineItemID"]

	var li LineItem
	json.NewDecoder(r.Body).Decode(&li)
	valid, errMsg := validateLineItem(li)
	if !valid {
		fmt.Printf("Bad value: %v\n", errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	query := `
		UPDATE LineItem
		SET OrderID = @p1,
		ProductID = @p2,
		LineItemUnitPrice = @p3,
		Quantity = @p4,
		LineItemDiscount = @p5
		OUTPUT INSERTED.LineItemID,
		INSERTED.OrderID, 
		INSERTED.ProductID,
		INSERTED.LineItemUnitPrice, 
		INSERTED.Quantity, 
		INSERTED.LineItemDiscount
		WHERE LineItemID = @p6`

	var updatedLineItem LineItem
	err := db.QueryRow(query,
		li.OrderID,
		li.ProductID,
		li.LineItemUnitPrice,
		li.Quantity,
		li.LineItemDiscount,
		id).Scan(&updatedLineItem.LineItemID,
		&updatedLineItem.OrderID,
		&updatedLineItem.ProductID,
		&updatedLineItem.LineItemUnitPrice,
		&updatedLineItem.Quantity,
		&updatedLineItem.LineItemDiscount)
	if err == sql.ErrNoRows {
		fmt.Printf("Line Item does not exist: %v\n", err)
		http.Error(w, "Line Item does not exist", http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Printf("Error updating line item: %v\n", err)
		http.Error(w, "Error updating line item", http.StatusInternalServerError)
		return
	}

	//Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedLineItem)
}
