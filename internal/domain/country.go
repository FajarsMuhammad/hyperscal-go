package domain

import (
	"time"
)

// Country represents the country entity in the database
type Country struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Code      string    `gorm:"type:varchar(10);uniqueIndex;not null" json:"code"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Region    string    `gorm:"type:varchar(50)" json:"region"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for Country model
func (Country) TableName() string {
	return "countries"
}
