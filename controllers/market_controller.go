package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AnhHoangQuach/go-intern-spores/models"
	"github.com/AnhHoangQuach/go-intern-spores/utils"
	"github.com/gin-gonic/gin"
)

type MarketController struct{}

var MarketModel = new(models.MarketModel)

func (m *MarketController) TotalRevenue(c *gin.Context) {
	cal_type := "date"
	time_query := time.Now().Day()
	var sum float64 = 0

	query := c.Request.URL.Query()
	for key, value := range query {
		queryValue := value[len(value)-1]
		if key == "cal_type" {
			cal_type = queryValue
		}

		if key == "date" || key == "month" || key == "year" {
			time_query, _ = strconv.Atoi(queryValue)
		}
	}

	txLists := MarketModel.CalculateRevenue(cal_type, time_query)

	for _, tx := range txLists {
		sum += (float64(tx.Price) + float64(tx.Fee))
	}
	res := utils.BuildResponse(true, "Success", sum)

	c.JSON(http.StatusOK, res)
}
