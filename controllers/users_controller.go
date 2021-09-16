package controllers

import (
	"net/http"

	"github.com/AnhHoangQuach/go-intern-spores/models"
	"github.com/AnhHoangQuach/go-intern-spores/services"
	"github.com/AnhHoangQuach/go-intern-spores/utils"
	"github.com/gin-gonic/gin"
)

type RegisterUserInput struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone"    binding:"required"`
	Address  string `json:"address"  binding:"required"`
	Cover    string `json:"cover"`
	Facebook string `json:"facebook"`
}

type VerifyInfo struct {
	Email       string `json:"email"        binding:"required"`
	VerifyToken string `json:"verify_token" binding:"required"`
}

type LoginUserInput struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ResetLinkInput struct {
	Email string `json:"email" binding:"required,email"`
}

type UserProfileInput struct {
	Phone    string `gorm:"size:255;not null;unique" json:"phone"`
	Address  string `gorm:"size:255;not null"        json:"address"`
	Cover    string `                                json:"cover"`
	Facebook string `                                json:"facebook"`
}

type ResetPasswordInput struct {
	Email       string `json:"email"        binding:"required,email"`
	ResetToken  string `json:"reset_token"  binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
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

	err := userModel.SignUp(
		input.Email,
		input.Password,
		input.Phone,
		input.Address,
		input.Cover,
		input.Facebook,
	)

	if err != nil {
		c.JSON(400, gin.H{"message": "Problem creating an account"})
		c.Abort()
		return
	}

	res := utils.BuildResponse(true, "Please check email to verify account", nil)

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
		c.JSON(400, utils.BuildErrorResponse("Active user failed", err.Error(), nil))
		c.Abort()
		return
	}

	res := utils.BuildResponse(true, "Active User Success", nil)

	c.JSON(http.StatusOK, res)
}

func (u *UserController) Login(c *gin.Context) {
	var input LoginUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := userModel.LoginHandler(input.Email, input.Password)

	if err != nil {
		c.JSON(404, utils.BuildErrorResponse("Login failed", err.Error(), nil))
		c.Abort()
		return
	}

	token, err := services.CreateJWT(user.Email)
	if err != nil {
		c.JSON(404, utils.BuildErrorResponse("Token is failed", err.Error(), nil))
		c.Abort()
		return
	}

	res := utils.BuildResponse(true, "Login Success", token)

	c.JSON(http.StatusOK, res)
}

func (u *UserController) Profile(c *gin.Context) {
	getUser, _ := c.Get("User")
	if getUser == nil {
		c.JSON(404, utils.BuildErrorResponse("Please Login", "Authenticate is failed", nil))
		c.Abort()
		return
	}
	user := getUser.(*models.User)
	if user.Email == "" {
		utils.BuildErrorResponse("Please login", "You not logged in", nil)
		return
	}

	result, err := userModel.FindByEmail(user.Email)
	if err != nil {
		c.JSON(404, utils.BuildErrorResponse("Authenticate is failed", err.Error(), nil))
		c.Abort()
		return
	}
	res := utils.BuildResponse(true, "Fetch Profile Success", result)

	c.JSON(http.StatusOK, res)
}

func (u *UserController) ResetByLink(c *gin.Context) {
	var input ResetLinkInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := userModel.ResetLink(input.Email)
	if err != nil {
		c.JSON(404, utils.BuildErrorResponse("Failed", err.Error(), nil))
		c.Abort()
		return
	}

	res := utils.BuildResponse(true, "Please check mail token to reset", nil)

	c.JSON(http.StatusOK, res)
}

func (u *UserController) ResetPasswordUser(c *gin.Context) {
	var input ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := userModel.ResetPassword(input.Email, input.ResetToken, input.NewPassword)
	if err != nil {
		c.JSON(404, utils.BuildErrorResponse("Failed", err.Error(), nil))
		c.Status(http.StatusBadRequest)
		return
	}

	res := utils.BuildResponse(true, "Reset password success", nil)

	c.JSON(http.StatusOK, res)
}

func (u *UserController) ChangeProfile(c *gin.Context) {
	getUser, _ := c.Get("User")
	if getUser == nil {
		c.JSON(404, utils.BuildErrorResponse("Please Login", "Authenticate is failed", nil))
		c.Abort()
		return
	}
	user := getUser.(*models.User)
	if user.Email == "" {
		utils.BuildErrorResponse("Please login", "You not logged in", nil)
		return
	}

	var input UserProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Phone = input.Phone
	user.Address = input.Address
	user.Cover = input.Cover
	user.Facebook = input.Facebook

	err := userModel.Update(user)
	if err != nil {
		c.JSON(404, utils.BuildErrorResponse("Failed", err.Error(), nil))
		c.Status(http.StatusBadRequest)
		return
	}

	res := utils.BuildResponse(true, "Change Profile Success", nil)

	c.JSON(http.StatusOK, res)
}
