package controller

import (
	"hyperscal-go/internal/dto"
	"hyperscal-go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *service.AuthService
}

func NewAuthController(service *service.AuthService) *AuthController {
	return &AuthController{
		service: service,
	}
}

// Register handles POST /api/auth/register
func (c *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", err.Error()))
		return
	}

	response, err := c.service.Register(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Failed to register", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dto.SuccessResponse("Login successful", response))
}

// Login handles POST /api/auth/login
func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", err.Error()))
		return
	}

	response, err := c.service.Login(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse("Failed to login", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Login successful", response))
}
