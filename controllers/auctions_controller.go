package controllers

import (
	"net/http"
	"strconv"

	"github.com/AnhHoangQuach/go-intern-spores/models"
	"github.com/AnhHoangQuach/go-intern-spores/utils"
	"github.com/gin-gonic/gin"
)

var aModel = new(models.AuctionModel)
var iModel = new(models.ItemModel)

type AuctionController struct{}

type UpdateAuctionInput struct {
	InitialPrice float64 `json:"initial_price"`
	FinalPrice   float64 `json:"final_price"`
	EndAt        int     `json:"end_at"`
}

type BidAuctionInput struct {
	Amount float64 `json:"amount"`
}

func (a *AuctionController) UpdateAuction(c *gin.Context) {
	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("ID is not valid", err.Error(), nil))
		return
	}

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

	auction, err := aModel.FindByID(uint32(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Auction is not existed", err.Error(), nil))
		return
	}

	item, err := iModel.FindByID(auction.ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Item is not existed", err.Error(), nil))
		return
	}

	if item.Owner != user.Email {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Failed", "You isn't owner of item", nil))
		return
	}

	var input UpdateAuctionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.InitialPrice != 0 && input.InitialPrice != auction.InitialPrice {
		auction.InitialPrice = input.InitialPrice
	}
	if input.FinalPrice != 0 && input.FinalPrice != auction.FinalPrice {
		auction.FinalPrice = input.FinalPrice
	}
	if input.EndAt != 0 {
		auction.EndAt = auction.CreatedAt.AddDate(0, 0, input.EndAt)
	}

	if input.InitialPrice == 0 && input.FinalPrice == 0 && input.EndAt == 0 {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Failed", "Please provide info to update item", nil))
		return
	}

	err = aModel.Update(auction)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Update Auction Failed", err.Error(), nil))
		return
	}

	res := utils.BuildResponse(true, "Update Auction Success", auction)
	c.JSON(http.StatusOK, res)
}

func (a *AuctionController) DeleteAuction(c *gin.Context) {
	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
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

	auction, err := aModel.FindByID(uint32(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Auction is not existed", err.Error(), nil))
		return
	}

	item, err := iModel.FindByID(auction.ItemID)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Item is not existed", err.Error(), nil))
		return
	}

	if item.Owner != user.Email {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Delete Item Failed", "You isn't owner of item", nil))
		return
	}

	err = aModel.Delete(uint32(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Delete Auction Failed", err.Error(), nil))
		return
	}

	res := utils.BuildResponse(true, "Delete Auction Success", nil)
	c.JSON(http.StatusOK, res)
}

func (a *AuctionController) BidAuction(c *gin.Context) {
	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
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

	auction, err := aModel.FindByID(uint32(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Auction is not existed", err.Error(), nil))
		return
	}

	item, err := iModel.FindByID(auction.ItemID)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Item is not existed", err.Error(), nil))
		return
	}

	if item.Owner == user.Email {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Bid Auction Failed", "This is your item", nil))
		return
	}

	var input BidAuctionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	auctionAfterBid, err := auctionModel.Bid(uint32(id), input.Amount)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Bid Auction Failed", err.Error(), nil))
		return
	}

	hash := utils.NewSHA1Hash()

	tx, err := txModel.Create(hash, item.ID, user.Email, item.Owner, input.Amount, float64(input.Amount)*0.1)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Transaction Failed", err.Error(), nil))
		return
	}

	result := struct {
		Tx      *models.Transaction `json:"tx"`
		Auction *models.Auction     `json:"auction"`
	}{
		Tx:      tx,
		Auction: auctionAfterBid,
	}

	res := utils.BuildResponse(true, "Bid Auction Success", result)
	c.JSON(http.StatusOK, res)
}
