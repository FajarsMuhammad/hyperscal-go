package dto

// CreateCityRequest represents the request body for creating a city
type CreateCityRequest struct {
	Name       string `json:"name" binding:"required,min=2,max=100"`
	Population int    `json:"population" binding:"required,gt=0"`
	CountryID  uint   `json:"country_id" binding:"required"`
}

type SearchCityRequest struct {
	PaginationRequest
	CountryID uint   `form:"country_id" binding:"omitempty,min=1"`
	Name      string `form:"name"`
}
