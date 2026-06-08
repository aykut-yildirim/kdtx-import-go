package main

import (
	"myapp/internal/api"
	"myapp/internal/services"
	"github.com/gin-gonic/gin"

	_ "myapp/internal/portals/amazon"
	_ "myapp/internal/portals/etsy"
	_ "myapp/internal/portals/stripe"
)

func main() {
	services.Logger().STATUS("-- Api / main.go")

	r := gin.Default()

	r.POST("/file_import", api.FileImport)

	r.POST("/api_import", api.ApiImport)

	r.GET("/healthy", api.Healthy)

	r.Run(":8000")
}
