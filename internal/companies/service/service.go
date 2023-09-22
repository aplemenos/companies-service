package service

import (
	"companies-service/config"
	"companies-service/internal/companies"
	"companies-service/internal/models"
	"companies-service/pkg/logger"
	"context"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
)

// Companies Service
type companiesService struct {
	cfg         *config.Config
	companyRepo companies.Repository
	logger      logger.Logger
}

// Companies Service constructor
func NewCompaniesService(
	cfg *config.Config, companyRepo companies.Repository, logger logger.Logger,
) companies.Service {
	return &companiesService{cfg: cfg, companyRepo: companyRepo, logger: logger}
}

// Create a new company
func (s *companiesService) Create(
	ctx context.Context, company *models.Company,
) (*models.Company, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "companiesService.Create")
	defer span.Finish()
	return s.companyRepo.Create(ctx, company)
}

// Update a company
func (s *companiesService) Update(
	ctx context.Context, company *models.Company,
) (*models.Company, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "companiesService.Update")
	defer span.Finish()

	updatedCompany, err := s.companyRepo.Update(ctx, company)
	if err != nil {
		return nil, err
	}

	return updatedCompany, nil
}

// Delete a company
func (s *companiesService) Delete(ctx context.Context, companyID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "companiesService.Delete")
	defer span.Finish()

	if err := s.companyRepo.Delete(ctx, companyID); err != nil {
		return err
	}

	return nil
}

// GetByID a company
func (s *companiesService) GetByID(
	ctx context.Context, companyID uuid.UUID,
) (*models.Company, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "companiesService.GetByID")
	defer span.Finish()

	return s.companyRepo.GetByID(ctx, companyID)
}
