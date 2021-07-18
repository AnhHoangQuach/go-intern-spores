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

	tx := new(controllers.TransactionController)

	itemApi := r.Group("/items")

	{
		itemApi.POST("/", middlewares.Authenticate(), item.CreateItem)
		itemApi.DELETE("/:id", middlewares.Authenticate(), item.DeleteItem)
		itemApi.GET("/:id", middlewares.Authenticate(), item.GetItem)
		itemApi.GET("/", middlewares.Authenticate(), item.GetAllItems)
		itemApi.PATCH("/:id", middlewares.Authenticate(), item.UpdateItem)
		itemApi.PUT("/:id", middlewares.Authenticate(), item.UpdateItem)
		itemApi.POST("/:id/buy", middlewares.Authenticate(), item.BuyItem)
		itemApi.GET("/:id/transactions", tx.GetTransOfItem)

	}

	auction := new(controllers.AuctionController)

	auctionApi := r.Group("/auction")

	{
		auctionApi.PATCH("/:id", middlewares.Authenticate(), auction.UpdateAuction)
		auctionApi.PUT("/:id", middlewares.Authenticate(), auction.UpdateAuction)
		auctionApi.DELETE("/:id", middlewares.Authenticate(), auction.DeleteAuction)
		auctionApi.POST("/:id/bid", middlewares.Authenticate(), auction.BidAuction)
	}

	market := new(controllers.MarketController)

	marketApi := r.Group("/market")

	{
		marketApi.GET("/revenue", middlewares.Authenticate(), market.TotalRevenue)
		marketApi.GET("/users", middlewares.Authenticate(), market.TotalUserRegister)
		marketApi.GET("/items/newest", market.GetItemsNewest)
		marketApi.GET("/auctions/newest", market.GetAuctionsNewest)
		marketApi.GET("/auctions/hot", market.HotAuctions)
		marketApi.GET("/items/hot", market.HotItems)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not found"})
	})

	// Run the server
	r.Run(":8080")
}
