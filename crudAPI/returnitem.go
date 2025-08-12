package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/microsoft/go-mssqldb"
)

type ReturnItem struct {
	ReturnedItemID int
	OrderID        int
	LineItemID     int
	Quantity       int
	AmountReturned float64
}

func GetReturnItems(w http.ResponseWriter, r *http.Request) {
	query := "SELECT * FROM ReturnItem"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("Error getting Return Items: %v\n", err)
		http.Error(w, "Error getting Return Items", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var returnItems []ReturnItem
	for rows.Next() {
		var returnItem ReturnItem
		err := rows.Scan(&returnItem.ReturnedItemID,
			&returnItem.OrderID,
			&returnItem.LineItemID,
			&returnItem.Quantity,
			&returnItem.AmountReturned)
		if err != nil {
			fmt.Printf("Scan error: %v\n", err)
			continue
		}
		returnItems = append(returnItems, returnItem)
	}

	//Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returnItems)
}

// Get a single return item, by ReturnedItemID
func GetReturnItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["ReturnedItemID"]

	query := "SELECT * FROM ReturnItem WHERE ReturnedItemID = @p1"
	var returnItem ReturnItem
	err := db.QueryRow(query, id).Scan(&returnItem.ReturnedItemID,
		&returnItem.OrderID,
		&returnItem.LineItemID,
		&returnItem.Quantity,
		&returnItem.AmountReturned)
	if err == sql.ErrNoRows {
		fmt.Printf("ReturnItem Not Found: %v\n", err)
		http.Error(w, "ReturnItem Not Found", http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Printf("Scan error: %v\n", err)
		http.Error(w, "Failed to fetch return item", http.StatusInternalServerError)
		return
	}

	//Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returnItem)
}

// Add a new return item
func AddReturnItem(w http.ResponseWriter, r *http.Request) {
	// TODO: Verify OrderID exists
	// TODO: Verify LineItemID exists
	// TODO: Verify quantity doesn't exceed original
	// TODO: Ensure AmountReturned is correct
	// TODO: Update Order/Customer totals
	var ri ReturnItem
	json.NewDecoder(r.Body).Decode(&ri)
	valid, errMsg := validateReturnItem(ri)
	if !valid {
		fmt.Printf("Bad return item data: %v\n", errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	//"OUTPUT" with QueryRow will get us a return value - use of Exec to insert would give none
	query := `INSERT INTO ReturnItem (OrderID, LineItemID, Quantity, AmountReturned) 
	OUTPUT INSERTED.ReturnedItemID, INSERTED.OrderID, INSERTED.LineItemID, INSERTED.Quantity, INSERTED.AmountReturned
	VALUES (@p1, @p2, @p3, @p4)`
	var newReturnItem ReturnItem

	err := db.QueryRow(query,
		ri.OrderID,
		ri.LineItemID,
		ri.Quantity,
		ri.AmountReturned).Scan(
		&newReturnItem.ReturnedItemID,
		&newReturnItem.OrderID,
		&newReturnItem.LineItemID,
		&newReturnItem.Quantity,
		&newReturnItem.AmountReturned)
	if err != nil {
		fmt.Printf("Error adding return item: %v\n", err)
		http.Error(w, "Error adding return item", http.StatusInternalServerError)
		return
	}
	//Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newReturnItem)
}

func validateReturnItem(returnItem ReturnItem) (bool, string) {
	// Validate OrderID
	if returnItem.OrderID <= 0 {
		return false, "OrderID is required and must be positive"
	}

	// Validate LineItemID
	if returnItem.LineItemID <= 0 {
		return false, "LineItemID is required and must be positive"
	}

	// Validate Quantity
	if returnItem.Quantity <= 0 {
		return false, "Quantity must be positive"
	}

	// Validate AmountReturned
	if returnItem.AmountReturned < 0 {
		return false, "AmountReturned cannot be negative"
	}
	// TODO: Verify AmountReturned doesn't exceed original line item value

	return true, ""
}

func UpdateReturnItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["ReturnedItemID"]

	// TODO: Verify OrderID exists
	// TODO: Verify LineItemID exists
	// TODO: Update Order/Customer totals
	var ri ReturnItem
	json.NewDecoder(r.Body).Decode(&ri)
	valid, errMsg := validateReturnItem(ri)
	if !valid {
		fmt.Printf("Bad value: %v\n", errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	query := `
		UPDATE ReturnItem
		SET OrderID = @p1,
		LineItemID = @p2,
		Quantity = @p3,
		AmountReturned = @p4
		OUTPUT INSERTED.ReturnedItemID,
		INSERTED.OrderID, 
		INSERTED.LineItemID,
		INSERTED.Quantity, 
		INSERTED.AmountReturned
		WHERE ReturnedItemID = @p5`

	var updatedReturnItem ReturnItem
	err := db.QueryRow(query,
		ri.OrderID,
		ri.LineItemID,
		ri.Quantity,
		ri.AmountReturned,
		id).Scan(&updatedReturnItem.ReturnedItemID,
		&updatedReturnItem.OrderID,
		&updatedReturnItem.LineItemID,
		&updatedReturnItem.Quantity,
		&updatedReturnItem.AmountReturned)
	if err == sql.ErrNoRows {
		fmt.Printf("Return Item does not exist: %v\n", err)
		http.Error(w, "Return Item does not exist", http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Printf("Error updating return item: %v\n", err)
		http.Error(w, "Error updating return item", http.StatusInternalServerError)
		return
	}

	//Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedReturnItem)
}

func DeleteReturnItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["ReturnedItemID"]

	// TODO: Update Order/Customer totals
	query := `DELETE FROM ReturnItem
		OUTPUT DELETED.ReturnedItemID
		WHERE ReturnedItemID = @p1`

	var deletedID int
	err := db.QueryRow(query, id).Scan(&deletedID)
	if err == sql.ErrNoRows {
		fmt.Printf("Return Item does not exist: %v\n", err)
		http.Error(w, "Return Item does not exist", http.StatusNotFound)
		return
	} else if err != nil && strings.Contains(err.Error(), "REFERENCE constraint") {
		fmt.Printf("Foreign key constraint error: %v\n", err)
		http.Error(w, "Cannot delete return item: it is referenced by existing records", http.StatusConflict)
		return
	} else if err != nil {
		fmt.Printf("Error deleting return item: %v\n", err)
		http.Error(w, "Error deleting return item", http.StatusInternalServerError)
		return
	}

	//Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedID)
}
