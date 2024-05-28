package controller

import (
	"go-app/src/service/dto"
	"go-app/src/service/spec"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v2

type RetailsController struct {
	svc spec.IRetailsService
}

func NewRetailsController(svc spec.IRetailsService) *RetailsController {
	return &RetailsController{svc: svc}
}

func (c *RetailsController) CustomerOrderStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dto.ApiResponse{ErrorCode: 0, Data: "hehe"})
}
