package domain

import (
	"time"
)

type City struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"type:varchar(100)" json:"name"`
	Population int       `gorm:"not null" json:"population"`
	CountryID  uint      `gorm:"not null" json:"country_id"`
	Country    Country   `gorm:"foreignKey:CountryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"country"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for City model
func (City) TableName() string {
	return "cities"
}
