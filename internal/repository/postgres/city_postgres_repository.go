package postgres

import (
	"errors"
	"hyperscal-go/internal/domain"
	"hyperscal-go/internal/dto"
	"strings"

	"gorm.io/gorm"
)

type CityPostgresRepository struct {
	db *gorm.DB
}

func NewCityPostgresRepository(db *gorm.DB) *CityPostgresRepository {
	return &CityPostgresRepository{
		db: db,
	}
}

// Implement CityRepository methods here (Create, FindAll, FindByID, Update, Delete)
func (r *CityPostgresRepository) Create(city *domain.City) error {
	return r.db.Create(city).Error
}

func (r *CityPostgresRepository) FindAll() ([]domain.City, error) {
	var cities []domain.City
	err := r.db.Preload("Country").Order("name ASC").Find(&cities).Error
	return cities, err
}

func (r *CityPostgresRepository) FindByID(id uint) (*domain.City, error) {
	var city domain.City
	err := r.db.Preload("Country").First(&city, id).Error
	if err != nil {
		return nil, err
	}
	return &city, nil
}

func (r *CityPostgresRepository) FindByNameAndCountryId(name string, countryId uint) (*domain.City, error) {
	var city domain.City
	err := r.db.Where("LOWER(name) = ? and country_id = ?", strings.ToLower(name), countryId).First(&city).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &city, nil
}

func (r *CityPostgresRepository) Update(city *domain.City) error {
	return r.db.Save(city).Error
}

func (r *CityPostgresRepository) Delete(id uint) error {
	return r.db.Delete(&domain.City{}, id).Error
}

func (r *CityPostgresRepository) Search(req dto.SearchCityRequest) ([]domain.City, int64, error) {
	var cities []domain.City
	var total int64

	//Base query
	query := r.db.Model(&domain.City{}).Preload("Country")

	if req.CountryID > 0 {
		query = query.Where("country_id = ?", req.CountryID)
	}

	if req.Name != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(req.Name)+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Set default values
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	// Apply sorting
	sortBy := req.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortDir := req.SortDir
	if sortDir == "" {
		sortDir = "desc"
	}
	query = query.Order(sortBy + " " + sortDir)

	// apply pagination
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	//Execute query
	if err := query.Find(&cities).Error; err != nil {
		return nil, 0, err
	}

	return cities, total, nil

}
