package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AnhHoangQuach/go-intern-spores/models"
	"github.com/AnhHoangQuach/go-intern-spores/services"
	"github.com/AnhHoangQuach/go-intern-spores/utils"
	"github.com/gin-gonic/gin"
)

type CreateItemInput struct {
	Name        string `json:"name" binding:"required`
	Description string `json:"description"`
	Price       int64  `json:"price" binding:"required`
	Currency    string `json:"currency" binding:"required`
	Owner       string `json:"owner" binding:"required`
	Creator     string `json:"creator" binding:"required`
	Metadata    string `json:"metadata" binding:"required`
}

type Pagination struct {
	Limit      int         `json:"limit"`
	Page       int         `json:"page"`
	Sort       string      `json:"sort"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Rows       interface{} `json:"rows"`
}

var itemModel = new(models.ItemModel)

type ItemController struct{}

func (i *ItemController) CreateItem(c *gin.Context) {
	getUser, _ := c.Get("User")
	if getUser == nil {
		c.JSON(404, utils.BuildErrorResponse("Please Login", "Authenticate is failed", nil))
		c.Abort()
		return
	}
	user := getUser.(*models.User)
	if user.Email == "" {
		c.JSON(404, utils.BuildErrorResponse("Please Login", "Authenticate is failed", nil))
		c.Abort()
		return
	}
	var input CreateItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := itemModel.Create(input.Name, input.Description, input.Currency, user.Email, user.Email, input.Price)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Problem creating item", err.Error(), nil))
		c.Abort()
		return
	}

	metadata := fmt.Sprintf("localhost:8080/items/%d", item.ID)
	err = itemModel.AddMetadataLink(item.ID, metadata)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Problem when add metadata link", err.Error(), nil))
		c.Abort()
		return
	}

	res := utils.BuildResponse(true, "Create Item Success", item)

	c.JSON(http.StatusOK, res)
}

func (i *ItemController) DeleteItem(c *gin.Context) {
	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("ID is not valid", err.Error(), nil))
		return
	}

	token := c.Request.Header["Authorization"]
	if len(token) == 0 || token[0] == "" {
		// Abort with error
		utils.BuildErrorResponse("Error", "You are not logged in", nil)
		return
	}

	user, err := services.ParseJWTToken(token[0])

	item, err := itemModel.FindByID(uint32(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Item is not existed", err.Error(), nil))
		return
	}

	if item.Owner != user.Email {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Delete Item Failed", "You isn't owner of item", nil))
		return
	}

	err = itemModel.Delete(uint32(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Delete Item Failed", err.Error(), nil))
		return
	}

	res := utils.BuildResponse(true, "Delete Item Success", nil)
	c.JSON(http.StatusOK, res)
}

func (i *ItemController) GetItem(c *gin.Context) {
	getUser, _ := c.Get("User")
	if getUser == nil {
		c.JSON(404, utils.BuildErrorResponse("Please Login", "Authenticate is failed", nil))
		c.Abort()
		return
	}

	user := getUser.(*models.User)

	if user.Email == "" {
		c.JSON(404, utils.BuildErrorResponse("Please Login", "Authenticate is failed", nil))
		c.Abort()
		return
	}

	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("ID is not valid", err.Error(), nil))
		return
	}

	item, err := itemModel.FindByID(uint32(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Item is not found", err.Error(), nil))
		return
	}

	res := utils.BuildResponse(true, "Fetch Item Success", item)
	c.JSON(http.StatusOK, res)
}

func (i *ItemController) GetAllItems(c *gin.Context) {
	getUser, _ := c.Get("User")
	if getUser == nil {
		c.JSON(404, utils.BuildErrorResponse("Please Login", "Authenticate is failed", nil))
		c.Abort()
		return
	}

	user := getUser.(*models.User)

	if user.Email == "" {
		c.JSON(404, utils.BuildErrorResponse("Please Login", "Authenticate is failed", nil))
		c.Abort()
		return
	}

	var item models.Item
	pagination := services.GeneratePaginationFromRequest(c)

	itemLists, totalRows, totalPages, err := itemModel.GetItemsPagination(&item, &pagination, user.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Failed when fetch pagination", err.Error(), nil))
		return
	}

	result := struct {
		Items      *[]models.Item `json:"items"`
		TotalPages int64          `json:"totalPages"`
		TotalRows  int64          `json:"totalRows"`
	}{
		Items:      itemLists,
		TotalPages: totalPages,
		TotalRows:  totalRows,
	}

	res := utils.BuildResponse(true, "Success", result)

	c.JSON(http.StatusOK, res)
}
