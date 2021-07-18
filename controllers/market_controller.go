package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AnhHoangQuach/go-intern-spores/models"
	"github.com/AnhHoangQuach/go-intern-spores/utils"
	"github.com/gin-gonic/gin"
)

type MarketController struct{}

var MarketModel = new(models.MarketModel)

func HandleDate(input int) string {
	out := strconv.Itoa(input)
	zero := fmt.Sprintf("%c", '0')
	if input < 10 {
		return zero + out
	}
	return out
}

func (m *MarketController) TotalRevenue(c *gin.Context) {
	day := time.Now().Day()
	month := int(time.Now().Month())
	year := time.Now().Year()

	queryType := 1

	yearStart := HandleDate(time.Now().Year())
	monthStart := HandleDate(int(time.Now().Month()))
	dateStart := HandleDate(time.Now().Day())
	tempChar := fmt.Sprintf("%c", '-')

	started := yearStart + tempChar + monthStart + tempChar + dateStart
	to := yearStart + tempChar + monthStart + tempChar + dateStart

	query := c.Request.URL.Query()
	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "day":
			day, _ = strconv.Atoi(queryValue)
			break
		case "month":
			month, _ = strconv.Atoi(queryValue)
			break
		case "year":
			year, _ = strconv.Atoi(queryValue)
			break
		case "started":
			started = queryValue
			queryType = 2
			break
		case "to":
			to = queryValue
			queryType = 2
			break
		}
	}

	sum := MarketModel.CalculateRevenue(day, month, year, queryType, started, to)
	res := utils.BuildResponse(true, "Success", sum)

	c.JSON(http.StatusOK, res)
}

func (m *MarketController) TotalUserRegister(c *gin.Context) {
	sum := MarketModel.CountUserInDay()
	res := utils.BuildResponse(true, "Success", sum)

	c.JSON(http.StatusOK, res)
}

func (m *MarketController) GetItemsNewest(c *gin.Context) {
	result := MarketModel.ListItemsNew()
	res := utils.BuildResponse(true, "Success", result)

	c.JSON(http.StatusOK, res)
}

func (m *MarketController) GetAuctionsNewest(c *gin.Context) {
	result := MarketModel.ListAuctionsNew()
	res := utils.BuildResponse(true, "Success", result)

	c.JSON(http.StatusOK, res)
}

func (m *MarketController) HotItems(c *gin.Context) {
	result := MarketModel.SellestItems()
	res := utils.BuildResponse(true, "Success", result)

	c.JSON(http.StatusOK, res)
}

func (m *MarketController) HotAuctions(c *gin.Context) {
	result := MarketModel.BighestAuctions()
	res := utils.BuildResponse(true, "Success", result)

	c.JSON(http.StatusOK, res)
}
