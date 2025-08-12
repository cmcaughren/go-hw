package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func connectDB() error {
	// Load .env file
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	server := os.Getenv("SERVER_NAME")
	user := os.Getenv("ADMIN_USER")
	password := os.Getenv("ADMIN_PASSWORD")
	database := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("sqlserver://%s:%s@%s:1433?database=%s&encrypt=true",
		user, password, server, database)

	var err error
	db, err = sql.Open("sqlserver", connStr)
	return err

}
