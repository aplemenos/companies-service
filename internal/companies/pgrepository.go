package companies

import (
	"companies-service/internal/models"
	"context"

	"github.com/google/uuid"
)

// Companies repository interface
type Repository interface {
	Create(ctx context.Context, company *models.Company) (*models.Company, error)
	Update(ctx context.Context, comment *models.Company) (*models.Company, error)
	Delete(ctx context.Context, companyID uuid.UUID) error
	GetByID(ctx context.Context, companyID uuid.UUID) (*models.Company, error)
}
