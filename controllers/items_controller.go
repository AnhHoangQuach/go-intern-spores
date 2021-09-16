package controllers

import (
	"fmt"
	"net/http"

	"github.com/AnhHoangQuach/go-intern-spores/models"
	"github.com/AnhHoangQuach/go-intern-spores/services"
	"github.com/AnhHoangQuach/go-intern-spores/utils"
	"github.com/gin-gonic/gin"
)

type CreateAuctionInput struct {
	InitialPrice float64 `json:"initial_price"`
	EndAt        int     `json:"end_at"`
}

type CreateItemInput struct {
	Name         string             `json:"name" binding:"required`
	Description  string             `json:"description"`
	Price        float64            `json:"price"`
	Currency     string             `json:"currency" binding:"required`
	Type         string             `json:"type" binding:"required`
	Image        string             `json:"image" binding:"required"`
	AuctionInput CreateAuctionInput `json:"create_auction_input"`
}

type UpdateItemInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	Image       string  `json:"image"`
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
var txModel = new(models.TxModel)
var auctionModel = new(models.AuctionModel)

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

	if input.Type == "Auction" {
		input.Price = input.AuctionInput.InitialPrice
	}

	item, err := itemModel.Create(
		input.Name,
		input.Description,
		input.Currency,
		user.Email,
		user.Email,
		input.Type,
		input.Image,
		input.Price,
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.BuildErrorResponse("Problem creating item", err.Error(), nil),
		)
		c.Abort()
		return
	}

	metadata := fmt.Sprintf("localhost:8080/items/%s", item.Id)
	err = itemModel.AddMetadataLink(item.Id, metadata)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.BuildErrorResponse("Problem when add metadata link", err.Error(), nil),
		)
		c.Abort()
		return
	}

	if input.Type == "Auction" {
		auction, err := auctionModel.Create(
			item.Id,
			input.AuctionInput.InitialPrice,
			input.AuctionInput.EndAt,
		)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				utils.BuildErrorResponse("Problem when create auction", err.Error(), nil),
			)
			c.Abort()
			return
		}

		result := struct {
			Item    *models.Item    `json:"item"`
			Auction *models.Auction `json:"auction"`
		}{
			Item:    item,
			Auction: auction,
		}

		res := utils.BuildResponse(true, "Create Item Success", result)
		c.JSON(http.StatusOK, res)
		return
	}

	res := utils.BuildResponse(true, "Create Item Success", item)

	c.JSON(http.StatusOK, res)
}

func (i *ItemController) DeleteItem(c *gin.Context) {
	id := c.Params.ByName("id")

	token := c.Request.Header["Authorization"]
	if len(token) == 0 || token[0] == "" {
		// Abort with error
		utils.BuildErrorResponse("Error", "You are not logged in", nil)
		return
	}

	user, err := services.ParseJWTToken(token[0])

	item, err := itemModel.FindByID(id)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.BuildErrorResponse("Item is not existed", err.Error(), nil),
		)
		return
	}

	if item.Owner != user.Email {
		c.JSON(
			http.StatusBadRequest,
			utils.BuildErrorResponse("Delete Item Failed", "You isn't owner of item", nil),
		)
		return
	}

	err = itemModel.Delete(id)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.BuildErrorResponse("Delete Item Failed", err.Error(), nil),
		)
		return
	}

	res := utils.BuildResponse(true, "Delete Item Success", nil)
	c.JSON(http.StatusOK, res)
}

func (i *ItemController) GetItem(c *gin.Context) {
	id := c.Params.ByName("id")

	item, err := itemModel.FindByID(id)

	auction, _ := auctionModel.FindByItemId(id)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.BuildErrorResponse("Item is not found", err.Error(), nil),
		)
		return
	}

	if item.Status == "Private" {
		getUser, _ := c.Get("User")
		if getUser == nil {
			c.JSON(404, utils.BuildErrorResponse("Please Login", "Authenticate is failed", nil))
			c.Abort()
			return
		}

		user := getUser.(*models.User)

		if user.Email == "" || item.Owner != user.Email {
			c.JSON(404, utils.BuildErrorResponse("Please Login", "Authenticate is failed", nil))
			c.Abort()
			return
		}
	}

	res := utils.BuildResponse(true, "Fetch Item Success", gin.H{
		"item":    item,
		"auction": auction,
	})
	c.JSON(http.StatusOK, res)
}

func (i *ItemController) GetPrivateItems(c *gin.Context) {
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

	itemLists, totalRows, totalPages, err := itemModel.Pagination(
		&item,
		&pagination,
		"Private",
		user.Email,
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.BuildErrorResponse("Failed when fetch pagination", err.Error(), nil),
		)
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

func (i *ItemController) UpdateItem(c *gin.Context) {
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

	item, err := itemModel.FindByID(c.Params.ByName("id"))

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.BuildErrorResponse("Item is not existed", err.Error(), nil),
		)
		return
	}

	if item.Owner != user.Email {
		c.JSON(
			http.StatusBadRequest,
			utils.BuildErrorResponse("Update Item Failed", "You isn't owner of item", nil),
		)
		return
	}

	var input UpdateItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != "" && input.Name != item.Name {
		item.Name = input.Name
	}
	if input.Description != "" && input.Description != item.Description {
		item.Description = input.Description
	}
	if input.Price != 0 && input.Price != item.Price {
		item.Price = input.Price
	}
	if input.Currency != "" && input.Currency != item.Currency {
		item.Currency = input.Currency
	}
	if input.Image != "" && input.Image != item.Image {
		item.Image = input.Image
	}

	if input.Name == "" && input.Description == "" && input.Price == 0 && input.Currency == "" &&
		input.Image == "" {
		c.JSON(
			http.StatusBadRequest,
			utils.BuildErrorResponse("Failed", "Please provide info to update item", nil),
		)
		return
	}

	err = itemModel.Update(item)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.BuildErrorResponse("Delete Item Failed", err.Error(), nil),
		)
		return
	}

	res := utils.BuildResponse(true, "Update Item Success", nil)
	c.JSON(http.StatusOK, res)
}

func (i *ItemController) BuyItem(c *gin.Context) {
	token := c.Request.Header["Authorization"]
	fmt.Println(token)
	if len(token) == 0 || token[0] == "" {
		// Abort with error
		utils.BuildErrorResponse("Error", "You are not logged in", nil)
		return
	}

	user, err := services.ParseJWTToken(token[0])

	item, err := itemModel.FindByID(c.Params.ByName("id"))

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.BuildErrorResponse("Item is not existed", err.Error(), nil),
		)
		return
	}

	if item.Owner == user.Email {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Failed", "This is your item", nil))
		return
	}

	if item.Status == "Private" {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Failed", "Item not on sale", nil))
		return
	}

	hash := utils.NewSHA1Hash()

	tx, err := txModel.Create(
		hash,
		item.Id,
		user.Email,
		item.Owner,
		item.Price,
		float64(item.Price)*0.1,
	)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.BuildErrorResponse("Transaction Failed", err.Error(), nil),
		)
		return
	}

	item.Owner = user.Email
	item.Status = "Private"
	err = itemModel.Update(item)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Buy Error", err.Error(), nil))
		return
	}

	res := utils.BuildResponse(true, "Buy Success", tx)
	c.JSON(http.StatusOK, res)
}

func (i *ItemController) GetPublicItems(c *gin.Context) {
	var item models.Item
	pagination := services.GeneratePaginationFromRequest(c)

	itemLists, totalRows, totalPages, err := itemModel.Pagination(&item, &pagination, "Public")

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.BuildErrorResponse("Failed when fetch pagination", err.Error(), nil),
		)
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

func (i *ItemController) PutOnMarket(c *gin.Context) {
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

	item, err := itemModel.FindByID(c.Params.ByName("id"))

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.BuildErrorResponse("Item is not existed", err.Error(), nil),
		)
		return
	}

	if item.Owner != user.Email {
		c.JSON(
			http.StatusBadRequest,
			utils.BuildErrorResponse("Failed", "This is not your item", nil),
		)
		return
	}

	item.Status = "Public"

	err = itemModel.Update(item)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Something error", err.Error(), nil))
		return
	}

	res := utils.BuildResponse(true, "Success", nil)

	c.JSON(http.StatusOK, res)

}
