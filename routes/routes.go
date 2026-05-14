package routes

import (
	"github.com/Abdulbarikassim/restaurant-api/handlers"

	"github.com/gin-gonic/gin"
)

	// menu routes
func SetUpRoutes (r *gin.Engine) {

		r.GET("/menu", handlers.GetMenuItems)
		r.GET("/menu/:id", handlers.GetMenuItemById)
		r.PUT("/menu/:id",handlers.UpdateMenuItem )
		r.POST("/menu",handlers.CreateMenuItem )
		r.DELETE("/menu/:id",handlers.DeleteMenuItem )

		// order routes

		r.GET("/orders", handlers.GetOrders)
		r.GET("/orders/:id", handlers.GetOrderById)
		r.POST("/orders", handlers.CreateNewOrder)
		r.PUT("/order/:id/status", handlers.UpdateOrderStatus)

	}
