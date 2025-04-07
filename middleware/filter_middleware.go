package middleware

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func FilterMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		queryParams := ctx.Request.URL.Query()
		queryStr, filters := parseFilters(queryParams)

		ctx.Set("filtersQuery", filters)
		ctx.Set("filtersSQL", queryStr)

		ctx.Next()
	}
}

func parseFilters(filters url.Values) (string, map[string]interface{}) {
	queryFilter := ""
	filtersValues := map[string]interface{}{}

	for key, value := range filters {
		split := strings.Split(value[0], ",")
		if len(split) < 2 {
			continue
		}

		op := split[0]
		val := split[1]

		if len(queryFilter) == 0 {
			queryFilter = "WHERE "
		} else {
			queryFilter += " AND "
		}

		switch op {
		case "eql":
			queryFilter += fmt.Sprintf("%s = :%s", key, key)
			filtersValues[key] = val
		case "lik":
			queryFilter += fmt.Sprintf("%s LIKE :%s", key, key)
			filtersValues[key] = "%" + val + "%"
		default:
		}
	}

	return queryFilter, filtersValues
}
