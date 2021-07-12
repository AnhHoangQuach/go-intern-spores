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
