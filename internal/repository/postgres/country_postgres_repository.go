package postgres

import (
	"hyperscal-go/internal/domain"

	"gorm.io/gorm"
)

// CountryPostgresRepository is the PostgreSQL implementation of CountryRepository
type CountryPostgresRepository struct {
	db *gorm.DB
}

// NewCountryPostgresRepository creates a new instance of CountryPostgresRepository
func NewCountryPostgresRepository(db *gorm.DB) *CountryPostgresRepository {
	return &CountryPostgresRepository{
		db: db,
	}
}

// Create inserts a new country into PostgreSQL database
func (r *CountryPostgresRepository) Create(country *domain.Country) error {
	return r.db.Create(country).Error
}

// FindAll retrieves all countries from PostgreSQL database
func (r *CountryPostgresRepository) FindAll() ([]domain.Country, error) {
	var countries []domain.Country
	err := r.db.Order("id ASC").Find(&countries).Error
	return countries, err
}

// FindByID retrieves a country by its ID from PostgreSQL database
func (r *CountryPostgresRepository) FindByID(id uint) (*domain.Country, error) {
	var country domain.Country
	err := r.db.First(&country, id).Error
	if err != nil {
		return nil, err
	}
	return &country, nil
}

// FindByCode retrieves a country by its code from PostgreSQL database
func (r *CountryPostgresRepository) FindByCode(code string) (*domain.Country, error) {
	var country domain.Country
	err := r.db.Where("code = ?", code).First(&country).Error
	if err != nil {
		return nil, err
	}
	return &country, nil
}

// Update updates an existing country in PostgreSQL database
func (r *CountryPostgresRepository) Update(country *domain.Country) error {
	return r.db.Save(country).Error
}

// Delete deletes a country by its ID from PostgreSQL database
func (r *CountryPostgresRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Country{}, id).Error
}
