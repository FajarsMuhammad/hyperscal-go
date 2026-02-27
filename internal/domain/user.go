package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password string `gorm:"type:varchar(255)not null" json:"password"`
	Name     string `gorm:"type:varchar(100)not null" json:"name"`
}
