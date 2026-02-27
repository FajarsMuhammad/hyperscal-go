package dto

// CreateCountryRequest represents the request body for creating a country
type CreateCountryRequest struct {
	Code   string `json:"code" binding:"required,min=2,max=10"`
	Name   string `json:"name" binding:"required,min=3,max=100"`
	Region string `json:"region" binding:"max=50"`
}

// UpdateCountryRequest represents the request body for updating a country
type UpdateCountryRequest struct {
	Code   string `json:"code" binding:"omitempty,min=2,max=10"`
	Name   string `json:"name" binding:"omitempty,min=3,max=100"`
	Region string `json:"region" binding:"omitempty,max=50"`
}
