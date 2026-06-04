package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"myapp/internal/importer"
	"myapp/internal/models"
)

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

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

	fmt.Println(task)

	c.JSON(
		http.StatusOK,
		SuccessResponse{
			Success: true,
			Data:    result,
		},
	)
}

func ApiImport(c *gin.Context) {
	fmt.Println("ApiImport")
}

func Healthy(c *gin.Context) {
	
	fmt.Println("Healthy")
	
	c.JSON(
		http.StatusOK,
		gin.H{
			"status": "ok",
		},
	)
}