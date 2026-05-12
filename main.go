package main

import (
	"database/sql"
	"encoding/json"
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
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	Available   bool    `json:"available"`
}

// order Item

type OrderItem struct {
	MenuItemId int     `json:"menu_item_id"`
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
}

// All the order items

type Order struct {
	ID            int      `json:"id"`
	CustomerPhone string      `json:"customer_phone"`
	Status        string      `json:"status"`
	TotalAmount   float64     `json:"total_amount"`
	Items         []OrderItem `json:"items"`
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
	router.GET("/menu/:id", getMenuItemById)
	router.PUT("/menu/:id", updateMenuItem)
	router.POST("/menu", createMenuItem)
	router.DELETE("/menu/:id", deleteMenuItem)

	// order routes

	router.GET("/orders", getOrders)
	router.GET("/orders/:id", getOrderById)
	router.POST("/orders", createNewOrder)
	router.PUT("/order/:id/status", updateOrderStatus)

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

		items = append(items, item)
	}

	c.JSON(http.StatusOK, items)
}

func createMenuItem(c *gin.Context) {
	var item MenuItem


	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}



	query := `INSERT INTO menu_items (name, description,category, price, available) VALUES ($1,$2,$3,$4,$5) RETURNING id`


	err := database.QueryRow(
		query,
		item.Name,
		item.Description,
		item.Category,
		item.Price,
		item.Available,
	).Scan(&item.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
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
		SET name = $1,
			description = $2,
			category = $3,
			price = $4,
			available = $5
		WHERE id = $6
	`

	result, err := database.Exec(
		query,
		item.Name,
		item.Description,
		item.Category,
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

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "menu item not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "menu item updated",
	})
}
// getMenuById function

func getMenuItemById(c *gin.Context) {

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

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "menu item not found",
			})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, item)
}
// delete menuItem func.

func deleteMenuItem(c *gin.Context) {

	id := c.Param("id")

	query := `DELETE FROM menu_items WHERE id = $1 `

	results , err := database.Exec(query, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"error" : err.Error(),
		})
		return
	}

	rowsEffected, err := results.RowsAffected()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if rowsEffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"eror": "menu item not found",
		})

		return
	}

	c.JSON(http.StatusOK,  gin.H{
		"message" : "menu item deleted successfully" ,
	})


}

// Order Processing logic

// getOrders function gets all the orders

func getOrders(c *gin.Context) {
	// query the database

	query := `
		SELECT
			id,
			customer_phone,
			status,
			total_amount,
			items
		FROM orders
		ORDER BY created_at DESC
	`

	rows, err := database.Query(query)

	if err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	defer rows.Close()

	// loop through and scan into go struct.


	var orders []Order
	for rows.Next() {

		var order Order

		var itemJSON []byte

		err := rows.Scan(
			&order.ID,
			&order.CustomerPhone,
			&order.Status,
			&order.TotalAmount,
			&itemJSON,
		)

		if err != nil {
			continue
		}

		json.Unmarshal(itemJSON, &order.Items)

		orders = append(orders, order )
	}

	c.JSON(http.StatusOK, orders)
}



// create new order function.

func createNewOrder(c *gin.Context) {

	var order Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}


	var totalAmount float64
	var enrichedItems []OrderItem

	for _, item := range order.Items {
		var dbItem MenuItem

		err := database.QueryRow(`
		SELECT name, price
		FROM menu_items
		WHERE id = $1
		`, item.MenuItemId).Scan(
			&dbItem.Name,
			&dbItem.Price,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// build safe order item

		enrichedItem := OrderItem{
			MenuItemId: item.MenuItemId,
			Name: dbItem.Name,
			Quantity: item.Quantity,
			Price: dbItem.Price,
		}

		enrichedItems = append(enrichedItems, enrichedItem)

		totalAmount  += dbItem.Price * float64(item.Quantity)


	}

	itemJSON, err := json.Marshal(enrichedItems)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = database.Exec(`
	INSERT INTO orders
	(customer_phone, total_amount,	Items)
	VALUES ($1,$2,$3)
	`,
	order.CustomerPhone,
	totalAmount,
	itemJSON,


	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// response

	c.JSON(http.StatusOK, gin.H{
		"message": "order created successfully",
		"total_amount": totalAmount,
		"items": enrichedItems,
	})

}

// getOrderById function.
func getOrderById(c *gin.Context) {

	id := c.Param("id")

	row := database.QueryRow(`
	SELECT id, customer_phone, status, total_amount, items
	FROM orders
	WHERE id = $1
	`, id)

	var order Order

	var itemJSON []byte

	err := row.Scan(
		&order.ID,
		&order.CustomerPhone,
		&order.Status,
		&order.TotalAmount,
		&itemJSON,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "order not found",
		})
		return
	}

	json.Unmarshal(itemJSON, &order.Items)

	c.JSON(http.StatusOK, order)
}

