package service

import (
	"companies-service/internal/companies/mock"
	"companies-service/internal/models"
	"companies-service/pkg/logger"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/require"
)

func TestCompaniesService_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockCompanyRepo := mock.NewMockRepository(ctrl)
	companiesService := NewCompaniesService(nil, mockCompanyRepo, apiLogger)

	company := &models.Company{
		CompanyID:          uuid.New(),
		CompanyName:        "Apple",
		CompanyDescription: "Computer and consumer electronics company",
		AmountOfEmployees:  164000,
		Registered:         true,
		CompanyType:        "Corporations",
	}

	span, ctx := opentracing.StartSpanFromContext(context.Background(), "companiesService.Create")
	defer span.Finish()

	mockCompanyRepo.EXPECT().Create(ctx, gomock.Eq(company)).Return(company, nil)

	createdCompany, err := companiesService.Create(context.Background(), company)
	require.NoError(t, err)
	require.NotNil(t, createdCompany)
}

func TestCompaniesService_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockCompanyRepo := mock.NewMockRepository(ctrl)
	companiesService := NewCompaniesService(nil, mockCompanyRepo, apiLogger)

	company := &models.Company{
		CompanyID:          uuid.New(),
		CompanyName:        "Apple",
		CompanyDescription: "Computer and consumer electronics company",
		AmountOfEmployees:  164000,
		Registered:         true,
		CompanyType:        "Corporations",
	}

	span, ctxWithTrace := opentracing.StartSpanFromContext(context.Background(),
		"companiesService.Update")
	defer span.Finish()

	mockCompanyRepo.EXPECT().Update(ctxWithTrace, gomock.Eq(company)).Return(company, nil)

	updatedCompany, err := companiesService.Update(context.Background(), company)
	require.NoError(t, err)
	require.NotNil(t, updatedCompany)
}

func TestCompaniesService_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockCompanyRepo := mock.NewMockRepository(ctrl)
	companiesService := NewCompaniesService(nil, mockCompanyRepo, apiLogger)

	company := &models.Company{
		CompanyID: uuid.New(),
	}

	span, ctxWithTrace := opentracing.StartSpanFromContext(context.Background(),
		"companiesService.Delete")
	defer span.Finish()

	mockCompanyRepo.EXPECT().Delete(ctxWithTrace, gomock.Eq(company.CompanyID)).Return(nil)

	err := companiesService.Delete(context.Background(), company.CompanyID)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestCompaniesService_GetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockCompanyRepo := mock.NewMockRepository(ctrl)
	companiesService := NewCompaniesService(nil, mockCompanyRepo, apiLogger)

	company := &models.Company{
		CompanyID:          uuid.New(),
		CompanyName:        "Apple",
		CompanyDescription: "Computer and consumer electronics company",
		AmountOfEmployees:  164000,
		Registered:         true,
		CompanyType:        "Corporations",
	}

	span, ctxWithTrace := opentracing.StartSpanFromContext(context.Background(),
		"companiesService.GetByID")
	defer span.Finish()

	mockCompanyRepo.EXPECT().GetByID(ctxWithTrace,
		gomock.Eq(company.CompanyID)).Return(company, nil)

	companyResponse, err := companiesService.GetByID(context.Background(), company.CompanyID)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, companyResponse)
}
