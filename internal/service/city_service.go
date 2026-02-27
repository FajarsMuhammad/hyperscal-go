package service

import (
	"errors"
	"hyperscal-go/internal/domain"
	"hyperscal-go/internal/dto"
	"hyperscal-go/internal/repository"
)

type CityService struct {
	repo repository.CityRepository
}

func NewCityService(repo repository.CityRepository) *CityService {
	return &CityService{
		repo: repo,
	}
}

// Implement CityService methods here (CreateCity, GetAllCities, GetCityByID, UpdateCity, DeleteCity)
func (s *CityService) CreateCity(req dto.CreateCityRequest) (*dto.CityResponse, error) {
	existing, err := s.repo.FindByNameAndCountryId(req.Name, req.CountryID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, errors.New("City already exists")
	}

	city := &domain.City{
		Name:       req.Name,
		Population: req.Population,
		CountryID:  req.CountryID,
	}

	if err := s.repo.Create(city); err != nil {
		return nil, err
	}

	// Fetch city with preloaded Country
	city, err = s.repo.FindByID(city.ID)
	if err != nil {
		return nil, err
	}

	return s.toResponse(city), nil
}

func (s *CityService) GetAllCities() ([]dto.CityResponse, error) {
	cities, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.CityResponse, len(cities))
	for i, city := range cities {
		responses[i] = *s.toResponse(&city)
	}

	return responses, nil
}

func (s *CityService) SearchCities(req dto.SearchCityRequest) (*dto.PaginationResponse, error) {
	cities, total, err := s.repo.Search(req)
	if err != nil {
		return nil, err
	}

	// 2. Convert domain ke response DTO
	cityResponses := make([]dto.CityResponse, len(cities))
	for i, city := range cities {
		cityResponses[i] = *s.toResponse(&city)
	}

	// 3. Set default pagination values
	page := req.Page
	if page < 1 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &dto.PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
		Data:       cityResponses,
	}, nil

}

// toResponse converts domain.City to dto.CityResponse
func (s *CityService) toResponse(city *domain.City) *dto.CityResponse {
	var countryID uint
	var countryName string

	if city.Country.ID != 0 {
		countryID = city.Country.ID
		countryName = city.Country.Name
	} else {
		countryID = city.CountryID
	}

	return &dto.CityResponse{
		ID:          city.ID,
		Name:        city.Name,
		Population:  city.Population,
		CountryID:   countryID,
		CountryName: countryName,
	}
}
