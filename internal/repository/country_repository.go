package repository

import (
	"hyperscal-go/internal/domain"
)

// CountryRepository defines the interface for country data access
// This interface can be implemented by different database adapters (PostgreSQL, Oracle, etc.)
type CountryRepository interface {
	// Create inserts a new country into the database
	Create(country *domain.Country) error

	// FindAll retrieves all countries from the database
	FindAll() ([]domain.Country, error)

	// FindByID retrieves a country by its ID
	FindByID(id uint) (*domain.Country, error)

	// FindByCode retrieves a country by its code
	FindByCode(code string) (*domain.Country, error)

	// Update updates an existing country
	Update(country *domain.Country) error

	// Delete deletes a country by its ID
	Delete(id uint) error
}
