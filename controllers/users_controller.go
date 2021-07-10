package controllers

import (
	"net/http"

	"github.com/AnhHoangQuach/go-intern-spores/models"
	"github.com/AnhHoangQuach/go-intern-spores/utils"
	"github.com/gin-gonic/gin"
)

type RegisterUserInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Address  string `json:"address" binding:"required"`
}

type VerifyInfo struct {
	Email       string `json:"email" binding:"required"`
	VerifyToken string `json:"verify_token" binding:"required"`
}

// Import the userModel from the models
var userModel = new(models.UserModel)

type UserController struct{}

func (u *UserController) SignUp(c *gin.Context) {
	var input RegisterUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := userModel.SignUp(input.Email, input.Password, input.Phone, input.Address)

	if err != nil {
		c.JSON(400, gin.H{"message": "Problem creating an account"})
		c.Abort()
		return
	}

	res := utils.BuildResponse(true, "Please check email to verify account", input)

	c.JSON(http.StatusOK, res)
}

func (u *UserController) VerifyUser(c *gin.Context) {
	var input VerifyInfo
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := userModel.ActiveUser(input.Email, input.VerifyToken)

	if err != nil {
		c.JSON(400, gin.H{"message": "Problem when verify account"})
		c.Abort()
		return
	}

	res := utils.BuildResponse(true, "Active User Success", input)

	c.JSON(http.StatusOK, res)
}

// func FindAllUsers(c *gin.Context) {
// 	var users []models.User
// 	models.DB.Find(&users)

// 	c.JSON(http.StatusOK, gin.H{"data": users})
// }

// func Delete(c *gin.Context) {
// 	var user models.User
// 	if err := models.DB.Where("email = ?", c.Param("email")).First(&user).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
// 		return
// 	}

// 	models.DB.Delete(&user)

// 	c.JSON(http.StatusOK, gin.H{"data": true})
// }
