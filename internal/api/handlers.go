package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"myapp/internal/importer"
	"myapp/internal/models"
	"myapp/internal/services"
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
	services.Logger().STATUS("--Api-FileImport")

	var task models.Task

	// services.Logger().STATUS(fmt.Sprintf("%+v", task))

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
	services.Logger().STATUS(fmt.Sprintf("%+v", task))

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
	services.Logger().STATUS("****" + fmt.Sprintf("%+v", task))

	c.JSON(
		http.StatusOK,
		SuccessResponse{
			Success: true,
			Data:    result,
		},
	)
}

func ApiImport(c *gin.Context) {
	services.Logger().STATUS("--Api-ApiImport")
}

func Healthy(c *gin.Context) {
	services.Logger().STATUS("--Api-Healthy")

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": "ok",
		},
	)
}
