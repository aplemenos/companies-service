package auth

import (
	"companies-service/internal/models"
	"context"

	"github.com/google/uuid"
)

// Auth repository interface
type Repository interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, user *models.User) (*models.User, error)
}
