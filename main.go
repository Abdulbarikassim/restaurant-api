package main

import (
	"database/sql"
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

var database *sql.DB

func main() {

	// load the .env file

	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found")
	}

	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("DATABASE_URL is not found")
	}

	database, err = db.ConnectDB()

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
	//router.GET("/menu/:id", getMenuById)
	router.PUT("/menu/:id", updateMenuItem)
	router.POST("/menu",createMenuItem )


	// order routes

	//router.GET("/orders", getOrders)
	//router.GET("/orders/:id", getOrderById)
	//router.POST("/orders", postNewOrder)
	//router.PUT("/order/:id/status",updateOrderStatus)

	router.Run(":8080")



}


// logic to get menu items

func getMenuItems(c *gin.Context) {
	rows, err := database.Query(`SELECT * FROM menu_items`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	defer rows.Close()

	var items []MenuItem

	for rows.Next() {

		var item MenuItem

		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Category,
			&item.Price,
			&item.Available,
		)

		if err != nil {
			continue
		}

		items =  append(items, item)
	}

	c.JSON(http.StatusOK, items)
}

func createMenuItem(c *gin.Context) {
	var item MenuItem


	if err  := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}


	query := `
		INSERT INTO menu_items
		(name, description, category,price, available)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id
	`

	err := database.QueryRow(
		query,
		item.Name,
		item.Description,
		item.Category,
		item.Price,
		item.Available,
	).Scan(&item.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"error" : err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// logic to update MenuItem


func updateMenuItem(c *gin.Context) {

	id := c.Param("id")

	var item MenuItem

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	query := `
		UPDATE menu_items
		SET
			price = $1,
			available = $2
		WHERE id = $3
	`

	_, err := database.Exec(
		query,
		item.Price,
		item.Available,
		id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "menu item updated",
	})

}


// getMenuById function

func getMenuById(c *gin.Context){

	id := c.Param("id")

	var item MenuItem

	query := `
	SELECT id, name, description,category, price, available FROM menu_items WHERE id = $1
	`

	err := database.QueryRow(query, id).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Category,
		&item.Price,
		&item.Available,
	)

	if err != nil {

		if err == sql.ErrNoRows{
			c.JSON(http.StatusNotFound,gin.H{
				"error": "menu item not found",
			})

			return
		}

		c.JSON(http.StatusInternalServerError,gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK,item)
}







