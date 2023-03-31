package main

import "github.com/gin-gonic/gin"

func main() {
	// Initialize gin default instance
	r := gin.Default()

	// Register a route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run() // listen and serve on
}
