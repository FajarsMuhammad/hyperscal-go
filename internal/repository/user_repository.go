package repository

import "hyperscal-go/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error

	FindByEmail(email string) (*domain.User, error)

	FindByID(id uint) (*domain.User, error)
}
