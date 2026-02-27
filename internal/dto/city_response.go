package dto

// CityResponse represents the response body for city operations
type CityResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Population  int    `json:"population"`
	CountryID   uint   `json:"country_id"`
	CountryName string `json:"country_name"`
}
