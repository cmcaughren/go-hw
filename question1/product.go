package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/microsoft/go-mssqldb"
)

type Product struct {
	ProductID   int
	ProductName string
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	query := "SELECT * FROM Product"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("Error getting Products: %v\n", err)
		http.Error(w, "Error getting Products", http.StatusInternalServerError)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ProductID,
			&product.ProductName)
		if err != nil {
			fmt.Printf("Scan error: %v\n", err)
			continue
		}
		products = append(products, product)
	}

	//Send JSON reponse
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// Get a single product, by product id
func GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["ProductID"]

	query := "SELECT * FROM Product WHERE ProductID = @p1"
	var product Product
	err := db.QueryRow(query, id).Scan(&product.ProductID,
		&product.ProductName)
	if err == sql.ErrNoRows {
		fmt.Printf("Product Not Found: %v\n", err)
		http.Error(w, "Product Not Found", http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Printf("Scan error: %v\n", err)
		http.Error(w, "Failed to fetch product", http.StatusInternalServerError)
		return
	}

	//Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}
