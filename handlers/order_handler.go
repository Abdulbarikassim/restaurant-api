package handlers

import (
	"net/http"
	"encoding/json"

	"github.com/gin-gonic/gin"

	m	"github.com/Abdulbarikassim/restaurant-api/models"

	"github.com/Abdulbarikassim/restaurant-api/db"
)

func GetOrders(c *gin.Context) {
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

	rows, err := db.DB.Query(query)

	if err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	defer rows.Close()

	// loop through and scan into go struct.


	var orders []m.Order
	for rows.Next() {

		var order m.Order

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

func CreateNewOrder(c *gin.Context) {

	var order m.Order
	// bind json -> go struct
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// store the total amount.
	var totalAmount float64
	var enrichedItems []m.OrderItem

	for _, item := range order.Items {
		var dbItem m.MenuItem

		err := db.DB.QueryRow(`
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

		enrichedItem := m.OrderItem{
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

	_, err = db.DB.Exec(`
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

	// send SMS using Tracksend API







	// response

	c.JSON(http.StatusCreated, gin.H{
		"message": "order created successfully",
		"total_amount": totalAmount,
		"items": enrichedItems,
	})

}

// getOrderById function.
func GetOrderById(c *gin.Context) {

	id := c.Param("id")

	row := db.DB.QueryRow(`
	SELECT id, customer_phone, status, total_amount, items
	FROM orders
	WHERE id = $1
	`, id)

	var order m.Order

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

// update status

func UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct{
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	validate := map[string]bool {
		"received": true ,
		"preparing": true,
		"completed": true,
		"cancelled": true,
	}

	if !validate[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid status",
		})
		return
	}

	_, err := db.DB.Exec(`
	UPDATE orders SET status = $1 WHERE id = $2
	`, req.Status,id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// response
	c.JSON(http.StatusOK, gin.H{
		"message": "order status updated successfully",
	})
}
