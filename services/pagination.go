package services

import (
	"strconv"
	"strings"

	"github.com/AnhHoangQuach/go-intern-spores/models"
	"github.com/gin-gonic/gin"
)

func GeneratePaginationFromRequest(c *gin.Context) models.Pagination {
	limit := 100
	page := 1
	sort := "created_at desc"

	var searchs []models.Search

	query := c.Request.URL.Query()
	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
			break
		case "page":
			page, _ = strconv.Atoi(queryValue)
			break
		case "sort":
			sort = queryValue
			break
		}

		// check if query parameter key contains dot
		if strings.Contains(key, ".") {
			// split query parameter key by dot
			searchKeys := strings.Split(key, ".")

			// create search object
			search := models.Search{Column: searchKeys[0], Action: searchKeys[1], Query: queryValue}

			// add search object to searchs array
			searchs = append(searchs, search)
		}
	}
	return models.Pagination{
		Limit:   limit,
		Page:    page,
		Sort:    sort,
		Searchs: searchs,
	}
}
