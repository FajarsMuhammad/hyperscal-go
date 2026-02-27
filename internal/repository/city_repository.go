package repository

import (
	"hyperscal-go/internal/domain"
	"hyperscal-go/internal/dto"
)

type CityRepository interface {
	// Create inserts a new city into the database
	Create(city *domain.City) error

	// FindAll retrieves all cities from the database
	FindAll() ([]domain.City, error)

	// FindByID retrieves a city by its ID
	FindByID(id uint) (*domain.City, error)

	//FindByName
	FindByNameAndCountryId(name string, id uint) (*domain.City, error)

	// Update updates an existing city
	Update(city *domain.City) error

	// Delete deletes a city by its ID
	Delete(id uint) error

	Search(req dto.SearchCityRequest) ([]domain.City, int64, error)
}
