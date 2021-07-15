package controllers

import (
	"time"

	"github.com/AnhHoangQuach/go-intern-spores/models"
	"github.com/gin-gonic/gin"
)

type CreateAuctionInput struct {
	InitialPrice float64   `json:"initial_price" binding:"required"`
	FinalPrice   float64   `json:"final_price" binding:"required"`
	EndAt        time.Time `json:"end_at" binding:"required"`
}

var auctionModel = new(models.AuctionModel)

type AuctionController struct{}

func (a *AuctionController) CreateAuction(c *gin.Context) {

}
