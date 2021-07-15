package controllers

import (
	"net/http"
	"strconv"

	"github.com/AnhHoangQuach/go-intern-spores/models"
	"github.com/AnhHoangQuach/go-intern-spores/services"
	"github.com/AnhHoangQuach/go-intern-spores/utils"
	"github.com/gin-gonic/gin"
)

type TransactionController struct{}

var transactionModel = new(models.TxModel)

func (t *TransactionController) GetTransOfItem(c *gin.Context) {
	var tx models.Transaction
	pagination := services.GeneratePaginationFromRequest(c)

	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("ID is not valid", err.Error(), nil))
		return
	}

	tranLists, totalRows, totalPages, err := transactionModel.TxPagination(&tx, &pagination, uint32(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Failed when fetch pagination", err.Error(), nil))
		return
	}

	result := struct {
		Txs        *[]models.Transaction `json:"transactions"`
		TotalPages int64                 `json:"totalPages"`
		TotalRows  int64                 `json:"totalRows"`
	}{
		Txs:        tranLists,
		TotalPages: totalPages,
		TotalRows:  totalRows,
	}

	res := utils.BuildResponse(true, "Success", result)

	c.JSON(http.StatusOK, res)
}
