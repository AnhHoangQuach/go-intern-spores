package middlewares

import (
	"github.com/AnhHoangQuach/go-intern-spores/models"
	"github.com/AnhHoangQuach/go-intern-spores/services"
	"github.com/AnhHoangQuach/go-intern-spores/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		requiredToken := c.Request.Header["Authorization"]
		if len(requiredToken) == 0 || requiredToken[0] == "" {
			// Abort with error
			utils.BuildErrorResponse("Authenticate is failed", "You are not logged in", nil)
			return
		}

		user, err := services.ParseJWTToken(requiredToken[0])
		if err != nil {
			utils.BuildErrorResponse("Authenticate is failed", err.Error(), nil)
			return
		}
		result, err := new(models.UserModel).FindByEmail(user.Email)

		if err != nil {
			utils.BuildErrorResponse("User account not found", err.Error(), nil)
			return
		}

		c.Set("User", result)

		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
