package oracle

import (
	"gorm.io/gorm"
)

type UserOracleRepository struct {
	db *gorm.DB
}

func NewUserOracleRepository(db *gorm.DB) *UserOracleRepository {
	return &UserOracleRepository{
		db: db,
	}
}

// Create, FindByEmail, FindByID methods sama seperti PostgreSQL
