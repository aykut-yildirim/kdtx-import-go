package main

import (
	"github.com/gin-gonic/gin"
	"myapp/internal/api"
	_ "myapp/internal/portals/amazon"
	_ "myapp/internal/portals/etsy"
	_ "myapp/internal/portals/stripe"
)

func main() {
	r := gin.Default()

	r.POST("/file_import", api.FileImport)

	r.POST("/api_import", api.ApiImport)

	r.GET("/healthy", api.Healthy)

	r.Run(":8000")
}
