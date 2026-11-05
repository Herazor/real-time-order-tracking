package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	ConnectDatabase()
	SeedData()

	router := gin.Default()

	router.POST("/api/orders", CreateOrder)
	router.GET("/api/orders", GetOrders)
	router.GET("/api/orders/:id", GetOrderByID)
	router.GET("/api/orders/:id/events", GetOrderEvents)
	router.PATCH("/api/orders/:id/status", UpdateOrderStatus)
	router.POST("/api/orders/:id/cancel", CancelOrder)
	router.PATCH("/api/cancellations/:id/review", ReviewCancellation)

	router.GET("/ws/orders/:id", HandleOrderWebSocket)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Real-Time Order Tracking API is running!",
		})
	})

	router.Run(":8080")
}