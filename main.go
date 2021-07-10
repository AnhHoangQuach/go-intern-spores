package main

import (
	"github.com/AnhHoangQuach/go-intern-spores/controllers"
	"github.com/AnhHoangQuach/go-intern-spores/models"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect DB
	r := gin.Default()

	// Connect to database
	models.ConnectDB()

	// Routes

	// Define the user controller
	user := new(controllers.UserController)

	userApi := r.Group("/user-api")
	{
		userApi.POST("/active", user.VerifyUser)
		userApi.POST("/signup", user.SignUp)
	}

	// Run the server
	r.Run(":8080")
}
