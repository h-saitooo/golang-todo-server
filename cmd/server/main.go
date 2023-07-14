package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateProfileRequest struct {
	Email string `json:"email"`
}

func main() {
	// Initialize gin default instance
	router := gin.Default()

	// Register a route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/user/profile/:name", updateProfile)

	router.Run() // listen and serve on
}

func updateProfile(context *gin.Context) {
	var request UpdateProfileRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	context.JSON(
		http.StatusOK,
		gin.H{
			"name":  context.Param("name"),
			"email": request.Email,
		},
	)
}
