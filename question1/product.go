package main

import (
	"encoding/json"
	"fmt"
	"net/http"

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
