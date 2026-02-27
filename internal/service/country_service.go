package service

import (
	"errors"
	"hyperscal-go/internal/domain"
	"hyperscal-go/internal/dto"
	"hyperscal-go/internal/repository"

	"gorm.io/gorm"
)

// CountryService handles business logic for country operations
type CountryService struct {
	repo repository.CountryRepository
}

// NewCountryService creates a new instance of CountryService
func NewCountryService(repo repository.CountryRepository) *CountryService {
	return &CountryService{
		repo: repo,
	}
}

// CreateCountry creates a new country
func (s *CountryService) CreateCountry(req dto.CreateCountryRequest) (*dto.CountryResponse, error) {
	// Check if country code already exists
	existing, err := s.repo.FindByCode(req.Code)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("country with this code already exists")
	}

	// Create new country
	country := &domain.Country{
		Code:   req.Code,
		Name:   req.Name,
		Region: req.Region,
	}

	if err := s.repo.Create(country); err != nil {
		return nil, err
	}

	return s.toResponse(country), nil
}

// GetAllCountries retrieves all countries
func (s *CountryService) GetAllCountries() ([]dto.CountryResponse, error) {
	countries, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.CountryResponse, len(countries))
	for i, country := range countries {
		responses[i] = *s.toResponse(&country)
	}

	return responses, nil
}

// GetCountryByID retrieves a country by its ID
func (s *CountryService) GetCountryByID(id uint) (*dto.CountryResponse, error) {
	country, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("country not found")
		}
		return nil, err
	}

	return s.toResponse(country), nil
}

// UpdateCountry updates an existing country
func (s *CountryService) UpdateCountry(id uint, req dto.UpdateCountryRequest) (*dto.CountryResponse, error) {
	// Find existing country
	country, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("country not found")
		}
		return nil, err
	}

	// Check if code is being changed and if new code already exists
	if req.Code != "" && req.Code != country.Code {
		existing, err := s.repo.FindByCode(req.Code)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("country with this code already exists")
		}
		country.Code = req.Code
	}

	// Update fields if provided
	if req.Name != "" {
		country.Name = req.Name
	}
	if req.Region != "" {
		country.Region = req.Region
	}

	if err := s.repo.Update(country); err != nil {
		return nil, err
	}

	return s.toResponse(country), nil
}

// DeleteCountry deletes a country by its ID
func (s *CountryService) DeleteCountry(id uint) error {
	// Check if country exists
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("country not found")
		}
		return err
	}

	return s.repo.Delete(id)
}

// toResponse converts domain.Country to dto.CountryResponse
func (s *CountryService) toResponse(country *domain.Country) *dto.CountryResponse {
	return &dto.CountryResponse{
		ID:        country.ID,
		Code:      country.Code,
		Name:      country.Name,
		Region:    country.Region,
		CreatedAt: country.CreatedAt,
		UpdatedAt: country.UpdatedAt,
	}
}
