package handlers


import (
	"database/sql"

	"net/http"
	m "github.com/Abdulbarikassim/restaurant-api/models"
	"github.com/Abdulbarikassim/restaurant-api/db"
	"github.com/gin-gonic/gin"



)


// logic to get menu items

func GetMenuItems(c *gin.Context) {
	rows, err := db.DB.Query(`SELECT * FROM menu_items`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	defer rows.Close()

	var items []m.MenuItem

	for rows.Next() {

		var item m.MenuItem

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

func CreateMenuItem(c *gin.Context) {
	var item m.MenuItem


	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}



	query := `INSERT INTO menu_items (name, description,category, price, available) VALUES ($1,$2,$3,$4,$5) RETURNING id`


	err := db.DB.QueryRow(
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
func UpdateMenuItem(c *gin.Context) {

	id := c.Param("id")

	var item m.MenuItem

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

	result, err := db.DB.Exec(
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

func GetMenuItemById(c *gin.Context) {

	id := c.Param("id")

	var item m.MenuItem

	query := `
	SELECT id, name, description,category, price, available FROM menu_items WHERE id = $1
	`

	err := db.DB.QueryRow(query, id).Scan(
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

func DeleteMenuItem(c *gin.Context) {

	id := c.Param("id")

	query := `DELETE FROM menu_items WHERE id = $1 `

	results , err := db.DB.Exec(query, id)

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

