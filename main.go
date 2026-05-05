package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/Abdulbarikassim/restaurant-api/db"
)

func main() {

	// load the .env file

	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found")
	}

	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("DATABASE_URL is not found")
	}

	database , err := db.ConnectDB()

	if 	err != nil {
		log.Fatal("DB connection failed: ", err)
	}

	defer database.Close()


	fmt.Println("Connected to postgreSQL successfully")

	err = database.Ping()

	if err != nil {
		log.Fatal("DB ping failed:", err)
	}

	fmt.Println("App is running")


}

