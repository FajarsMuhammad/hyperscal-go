package controller

import (
	"hyperscal-go/internal/dto"
	"hyperscal-go/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CountryController handles HTTP requests for country operations
type CountryController struct {
	service *service.CountryService
}

// NewCountryController creates a new instance of CountryController
func NewCountryController(service *service.CountryService) *CountryController {
	return &CountryController{
		service: service,
	}
}

// CreateCountry handles POST /api/countries
func (c *CountryController) CreateCountry(ctx *gin.Context) {
	var req dto.CreateCountryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", err.Error()))
		return
	}

	country, err := c.service.CreateCountry(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Failed to create country", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, dto.SuccessResponse("Country created successfully", country))
}

// GetAllCountries handles GET /api/countries
func (c *CountryController) GetAllCountries(ctx *gin.Context) {
	countries, err := c.service.GetAllCountries()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to retrieve countries", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Countries retrieved successfully", countries))
}

// GetCountryByID handles GET /api/countries/:id
func (c *CountryController) GetCountryByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid country ID", err.Error()))
		return
	}

	country, err := c.service.GetCountryByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Country not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Country retrieved successfully", country))
}

// UpdateCountry handles PUT /api/countries/:id
func (c *CountryController) UpdateCountry(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid country ID", err.Error()))
		return
	}

	var req dto.UpdateCountryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", err.Error()))
		return
	}

	country, err := c.service.UpdateCountry(uint(id), req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Failed to update country", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Country updated successfully", country))
}

// DeleteCountry handles DELETE /api/countries/:id
func (c *CountryController) DeleteCountry(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid country ID", err.Error()))
		return
	}

	if err := c.service.DeleteCountry(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Failed to delete country", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Country deleted successfully", nil))
}
