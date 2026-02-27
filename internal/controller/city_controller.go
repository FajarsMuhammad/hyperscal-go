package controller

import (
	"hyperscal-go/internal/dto"
	"hyperscal-go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CityController struct {
	service *service.CityService
}

func NewCityController(service *service.CityService) *CityController {
	return &CityController{
		service: service,
	}
}

// Implement CityController methods here (CreateCity, GetAllCities, GetCityByID, UpdateCity, DeleteCity)
func (c *CityController) CreateCity(ctx *gin.Context) {
	var req dto.CreateCityRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", err.Error()))
		return
	}

	city, err := c.service.CreateCity(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to create city", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, dto.SuccessResponse("City created successfully", city))

}

func (c *CityController) GetAllCities(ctx *gin.Context) {

	cities, err := c.service.GetAllCities()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to load cities", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("City fetch successfully", cities))
}

func (c *CityController) SearchCities(ctx *gin.Context) {
	var req dto.SearchCityRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Failed to load cities", err.Error()))
		return
	}

	result, err := c.service.SearchCities(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to search cities", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Citie fetches successfully", result))
}
