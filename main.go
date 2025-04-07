package main

import (
	"api_boilerplate/db"
	"api_boilerplate/util"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Health")
	})

	dbConn := db.GetDBConnection()
	defer dbConn.Close()

	util.RegisterDomains(r, dbConn)

	r.Run("0.0.0.0:3030")
}
