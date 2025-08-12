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

	r.HandleFunc("/customer-addresses", GetCustomerAddresses).Methods("GET")
	r.HandleFunc("/customer-address/{CustomerAddressID}", GetCustomerAddress).Methods("GET")
	r.HandleFunc("/customer-address", AddCustomerAddress).Methods("POST")
	r.HandleFunc("/customer-address/{CustomerAddressID}", UpdateCustomerAddress).Methods("PUT")
	r.HandleFunc("/customer-address/{CustomerAddressID}", DeleteCustomerAddress).Methods("DELETE")

	r.HandleFunc("/products", GetProducts).Methods("GET")
	r.HandleFunc("/product/{ProductID}", GetProduct).Methods("GET")
	r.HandleFunc("/product", AddProduct).Methods("POST")
	r.HandleFunc("/product/{ProductID}", UpdateProduct).Methods("PUT")
	r.HandleFunc("/product/{ProductID}", DeleteProduct).Methods("DELETE")

	r.HandleFunc("/orders", GetOrders).Methods("GET")
	r.HandleFunc("/order/{OrderID}", GetOrder).Methods("GET")
	r.HandleFunc("/order", AddOrder).Methods("POST")
	r.HandleFunc("/order/{OrderID}", UpdateOrder).Methods("PUT")
	r.HandleFunc("/order/{OrderID}", DeleteOrder).Methods("DELETE")

	r.HandleFunc("/line-items", GetLineItems).Methods("GET")
	r.HandleFunc("/line-item/{LineItemID}", GetLineItem).Methods("GET")
	r.HandleFunc("/line-item", AddLineItem).Methods("POST")
	r.HandleFunc("/line-item/{LineItemID}", UpdateLineItem).Methods("PUT")
	r.HandleFunc("/line-item/{LineItemID}", DeleteLineItem).Methods("DELETE")

	r.HandleFunc("/return-items", GetReturnItems).Methods("GET")
	r.HandleFunc("/return-item/{ReturnedItemID}", GetReturnItem).Methods("GET")
	r.HandleFunc("/return-item", AddReturnItem).Methods("POST")
	r.HandleFunc("/return-item/{ReturnedItemID}", UpdateReturnItem).Methods("PUT")
	r.HandleFunc("/return-item/{ReturnedItemID}", DeleteReturnItem).Methods("DELETE")

	//Start Server
	port := ":8080"
	fmt.Printf("Starting server on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))

}
