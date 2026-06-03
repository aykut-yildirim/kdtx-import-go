package api

import  (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()

	r.POST("/file_import", FileImport)
	r.POST("/api_import", ApiImport)

	r.GET("/healthy", Healthy)

	return r
}