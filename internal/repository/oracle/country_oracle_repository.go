package oracle

import (
	"hyperscal-go/internal/domain"

	"gorm.io/gorm"
)

// CountryOracleRepository is the Oracle implementation of CountryRepository
type CountryOracleRepository struct {
	db *gorm.DB
}

// NewCountryOracleRepository creates a new instance of CountryOracleRepository
func NewCountryOracleRepository(db *gorm.DB) *CountryOracleRepository {
	return &CountryOracleRepository{
		db: db,
	}
}

// Create inserts a new country into Oracle database
func (r *CountryOracleRepository) Create(country *domain.Country) error {
	return r.db.Create(country).Error
}

// FindAll retrieves all countries from Oracle database
func (r *CountryOracleRepository) FindAll() ([]domain.Country, error) {
	var countries []domain.Country
	err := r.db.Order("id ASC").Find(&countries).Error
	return countries, err
}

// FindByID retrieves a country by its ID from Oracle database
func (r *CountryOracleRepository) FindByID(id uint) (*domain.Country, error) {
	var country domain.Country
	err := r.db.First(&country, id).Error
	if err != nil {
		return nil, err
	}
	return &country, nil
}

// FindByCode retrieves a country by its code from Oracle database
func (r *CountryOracleRepository) FindByCode(code string) (*domain.Country, error) {
	var country domain.Country
	err := r.db.Where("code = ?", code).First(&country).Error
	if err != nil {
		return nil, err
	}
	return &country, nil
}

// Update updates an existing country in Oracle database
func (r *CountryOracleRepository) Update(country *domain.Country) error {
	return r.db.Save(country).Error
}

// Delete deletes a country by its ID from Oracle database
func (r *CountryOracleRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Country{}, id).Error
}
