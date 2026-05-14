package main

import (
	"log"

	"github.com/Abdulbarikassim/restaurant-api/db"
	"github.com/Abdulbarikassim/restaurant-api/routes"
	"github.com/gin-gonic/gin"
)




func main() {

	_, err := db.ConnectDB()

	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	routes.SetUpRoutes(r)

	log.Println("Server running on : 8080")

	r.Run(":8080")
}

