package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/microsoft/go-mssqldb"
)

var db *sql.DB

func main() {
	if err := connectDB(); err != nil {
		log.Fatal("Unable to connect to DB:", err)
	}
	defer db.Close()

	fmt.Print("DB successfully connected.\n")

	// Set up HTTP routes
	r := mux.NewRouter()
	r.HandleFunc("/customers", GetCustomers).Methods("GET")
	r.HandleFunc("/customer/{CustomerID}", GetCustomer).Methods("GET")
	r.HandleFunc("/customer", AddCustomer).Methods("POST")
	r.HandleFunc("/customer/{CustomerID}", UpdateCustomer).Methods("PUT")
	r.HandleFunc("/customer/{CustomerID}", DeleteCustomer).Methods("DELETE")

	r.HandleFunc("/customerAddresses", GetCustomerAddresses).Methods("GET")
	r.HandleFunc("/customerAddress/{CustomerAddressID}", GetCustomerAddress).Methods("GET")
	r.HandleFunc("/customerAddress", AddCustomerAddress).Methods("POST")
	r.HandleFunc("/customerAddress/{CustomerAddressID}", UpdateCustomerAddress).Methods("PUT")
	r.HandleFunc("/customerAddress/{CustomerAddressID}", DeleteCustomerAddress).Methods("DELETE")

	r.HandleFunc("/products", GetProducts).Methods("GET")
	r.HandleFunc("/product/{ProductID}", GetProduct).Methods("GET")
	r.HandleFunc("/product", AddProduct).Methods("POST")
	//r.HandleFunc("/product/{ProductID}", UpdateProduct).Methods("PUT")
	//r.HandleFunc("/product/{ProductID}", DeleteProduct).Methods("DELETE")

	//Start Server
	port := ":8080"
	fmt.Printf("Starting server on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))

}
