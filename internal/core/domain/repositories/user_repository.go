package repositories

import (
	"context" // To handle request contexts
	"go-playground/internal/core/domain/models"
)

// UserRepository methods definition also known as interface or The Contract
type UserRepository interface {
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
	GetAllUsers(ctx context.Context) ([]*models.User, error)
}
