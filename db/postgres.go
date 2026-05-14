package db

import (
	"database/sql"
	"log"
	"os"
	"fmt"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() (*sql.DB, error) {
	// load the .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// check the DATABASE_URL
	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("DATABASE_URL is not found")
	}

	// connect the database
	connStr :=  os.Getenv("DATABASE_URL")

	database, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	// ping the database
	err = database.Ping()

	if err != nil {
		return nil , err
	}

	fmt.Println("connected PostgreSQL successfully")

	DB = database


	return DB, nil

}
