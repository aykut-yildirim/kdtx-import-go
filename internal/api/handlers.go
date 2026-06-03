package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"myapp/internal/importer"
	"myapp/internal/models"
)

func FileImport(c *gin.Context) {
	fmt.Println("FileImport")
	
	var task models.Task
	
	fmt.Println(task)
	if err := c.ShouldBindJSON(&task); err != nil {
		
		c.JSON(
			http.StatusBadRequest,
			ErrorResponse{
				Success: false,
				Error:   err.Error(),
			},
		)
		
		return
	}
	fmt.Println(task)

	result, err := importer.Service{}.Execute(task)

	if err != nil {

		c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{
				Success: false,
				Error:   err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		SuccessResponse{
			Success: true,
			Data:    result,
		},
	)
}

func ApiImport(c *gin.Context) {
	FileImport(c)
}

func Healthy(c *gin.Context) {

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": "ok",
		},
	)
}