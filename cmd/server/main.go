package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"todo_app/internal/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type CreateProfileRequest struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	PlainPassword string `json:"password"`
}

func main() {
	// Wait for database to be ready
	time.Sleep(3 * time.Second)

	if err := database.Initialize(); err != nil {
		log.Fatalf("failed to initialize database: %s", err.Error())
	}
	defer database.Close()

	// Initialize gin default instance
	router := gin.Default()

	// Create router group
	v1 := router.Group("/api/v1")

	// Register a route
	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	v1.POST("/user/profile", CreateProfile)
	v1.GET("/user/profile/:name", GetProfileByName)

	router.Run() // listen and serve on
}

func CreateProfile(context *gin.Context) {
	var request CreateProfileRequest
	now := time.Now()

	fmt.Printf("Now: %s\n", now.Format("2006-01-02 15:04:05.999999"))

	err := context.ShouldBindJSON(&request)
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.PlainPassword), bcrypt.DefaultCost)
	if err != nil {
		context.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return

	}
	fmt.Printf("Plain Password: %s\n", request.PlainPassword)
	fmt.Printf("Hashed Password: %s\n", hashedPassword)

	_, err = database.DB.Exec(
		`
			INSERT INTO users (name, email, password_hash, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?)
		`,
		request.Name, request.Email, string(hashedPassword), now, now,
	)
	if err != nil {
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

func GetProfileByName(context *gin.Context) {
	requestName := context.Param("name")
	var userId int
	var email string
	var name string

	err := database.DB.QueryRow(
		"SELECT id, name, email FROM users WHERE name = ?",
		requestName,
	).Scan(&userId, &name, &email)

	if err == sql.ErrNoRows {
		context.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "user not found",
			},
		)
		return
	}

	if err != nil {
		context.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	context.JSON(
		http.StatusOK,
		gin.H{
			"id":    userId,
			"name":  name,
			"email": email,
		},
	)
}
