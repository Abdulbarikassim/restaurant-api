package main

import (
	"fmt"
	"log"
	"os"
	"github.com/gin-gonic/gin"

	"net/http"
	"github.com/Abdulbarikassim/restaurant-api/db"
	"github.com/joho/godotenv"
)

// menu item

type MenuItem struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Category string `json:"category"`
	Price float64 `json:"price"`
	Available bool `json:"available"`
}

// order Item

type OrderItem struct {
	MenuItemId int `json:"menu_item_id"`
	Name string `json:"name"`
	Quantity int `json:"quantity"`
	Price float64 `json:"price"`
}

// All the order items

type Order struct {
	ID 	string `json:"id"`
	CustomerPhone string `json:"customer_phone"`
	Status string `json:"status"`
	TotalAmount float64 `json:"total_amount"`
	Items []OrderItem 	`json:"items"`
}


func main() {

	// load the .env file

	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found")
	}

	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("DATABASE_URL is not found")
	}

	database, err := db.ConnectDB()

	if err != nil {
		log.Fatal("DB connection failed: ", err)
	}

	defer database.Close()

	fmt.Println("Connected to postgreSQL successfully")

	err = database.Ping()

	if err != nil {
		log.Fatal("DB ping failed:", err)
	}

	fmt.Println("App is running")


	router := gin.Default()

	// menu routes

	router.GET("/menu", getMenuItems)
	router.GET("/menu/:id", getMenuById)
	router.PUT("/menu/:id", updateMenu)
	router.POST("/menu", postNewMenu)


	// order routes

	router.GET("/orders", getOrders)
	router.GET("/orders/:id", getOrderById)
	router.POST("/orders", postNewOrder)
	router.PUT("/order/:id/status",updateOrderStatus)







}
