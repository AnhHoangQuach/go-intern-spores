package main

import (
	"github.com/AnhHoangQuach/go-intern-spores/controllers"
	"github.com/AnhHoangQuach/go-intern-spores/middlewares"
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

	userApi := r.Group("/auth")
	{
		userApi.POST("/verifyEmail", user.VerifyUser)
		userApi.POST("/register", user.SignUp)
		userApi.POST("/login", user.Login)
		userApi.GET("/profile", middlewares.Authenticate(), user.Profile)
	}

	item := new(controllers.ItemController)

	itemApi := r.Group("/items")

	{
		itemApi.POST("/", middlewares.Authenticate(), item.CreateItem)
		itemApi.DELETE("/:id", middlewares.Authenticate(), item.DeleteItem)
		itemApi.GET("/:id", middlewares.Authenticate(), item.GetItem)
		itemApi.GET("/", middlewares.Authenticate(), item.GetAllItems)
		itemApi.PATCH("/:id", middlewares.Authenticate(), item.UpdateItem)
		itemApi.PUT("/:id", middlewares.Authenticate(), item.UpdateItem)
		itemApi.POST("/buy/:id", middlewares.Authenticate(), item.BuyItem)

	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not found"})
	})

	// Run the server
	r.Run(":8080")
}
