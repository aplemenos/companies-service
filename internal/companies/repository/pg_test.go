package repository

import (
	"companies-service/internal/models"
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestCompaniesRepo_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	companiesRepo := NewCompaniesRepository(sqlxDB)

	companyName := "Apple"
	companyDescription := "Computer and consumer electronics company"
	amountOfEmployees := 164000
	registered := true
	companyType := "Corporations"

	t.Run("Create", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"company_name", "company_description",
			"amount_of_employees", "registered", "company_type",
		}).AddRow(companyName, companyDescription,
			amountOfEmployees, registered, companyType)

		company := &models.Company{
			CompanyName:        companyName,
			CompanyDescription: companyDescription,
			AmountOfEmployees:  amountOfEmployees,
			Registered:         registered,
			CompanyType:        companyType,
		}

		mock.ExpectQuery(createCompany).WithArgs(company.CompanyName,
			company.CompanyDescription, company.AmountOfEmployees, company.Registered,
			company.CompanyType).WillReturnRows(rows)

		createdCompany, err := companiesRepo.Create(context.Background(), company)

		require.NoError(t, err)
		require.NotNil(t, createdCompany)
		require.Equal(t, createdCompany, company)
	})

	t.Run("Create ERR", func(t *testing.T) {
		createErr := errors.New("Create company error")

		company := &models.Company{
			CompanyName:        companyName,
			CompanyDescription: companyDescription,
			AmountOfEmployees:  amountOfEmployees,
			Registered:         registered,
			CompanyType:        "Unknown",
		}

		mock.ExpectQuery(createCompany).WithArgs(company.CompanyName,
			company.CompanyDescription, company.AmountOfEmployees, company.Registered,
			company.CompanyType).WillReturnError(createErr)

		createdCompany, err := companiesRepo.Create(context.Background(), company)

		require.Nil(t, createdCompany)
		require.NotNil(t, err)
	})
}

func TestCompaniesRepo_Update(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	companiesRepo := NewCompaniesRepository(sqlxDB)

	companyID := uuid.New()
	companyName := "Apple"
	companyDescription := "Computer and consumer electronics company"
	amountOfEmployees := 164000
	registered := true
	companyType := "Corporations"

	t.Run("Update", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"company_id", "company_name", "company_description",
			"amount_of_employees", "registered", "company_type",
		}).AddRow(companyID, companyName, companyDescription,
			amountOfEmployees, registered, companyType)

		company := &models.Company{
			CompanyID:          companyID,
			CompanyName:        companyName,
			CompanyDescription: companyDescription,
			AmountOfEmployees:  amountOfEmployees,
			Registered:         registered,
			CompanyType:        companyType,
		}

		mock.ExpectQuery(updateCompany).WithArgs(company.CompanyName, company.CompanyDescription,
			company.AmountOfEmployees, company.Registered, company.CompanyType,
			company.CompanyID).WillReturnRows(rows)

		updatedCompany, err := companiesRepo.Update(context.Background(), company)

		require.NoError(t, err)
		require.NotNil(t, updatedCompany)
		require.Equal(t, updatedCompany, company)
	})

	t.Run("Update ERR", func(t *testing.T) {
		updateErr := errors.New("Update company error")

		company := &models.Company{
			CompanyID:          companyID,
			CompanyName:        companyName,
			CompanyDescription: companyDescription,
			AmountOfEmployees:  amountOfEmployees,
			Registered:         registered,
			CompanyType:        "Unknown",
		}

		mock.ExpectQuery(updateCompany).WithArgs(company.CompanyName, company.CompanyDescription,
			company.AmountOfEmployees, company.Registered, company.CompanyType,
			company.CompanyID).WillReturnError(updateErr)

		updatedCompany, err := companiesRepo.Update(context.Background(), company)

		require.NotNil(t, err)
		require.Nil(t, updatedCompany)
	})
}

func TestCompaniesRepo_Delete(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	companiesRepo := NewCompaniesRepository(sqlxDB)

	t.Run("Delete", func(t *testing.T) {
		companyID := uuid.New()
		mock.ExpectExec(deleteCompany).WithArgs(companyID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := companiesRepo.Delete(context.Background(), companyID)

		require.NoError(t, err)
	})

	t.Run("Delete Err", func(t *testing.T) {
		companyID := uuid.New()
		mock.ExpectExec(deleteCompany).WithArgs(companyID).
			WillReturnResult(sqlmock.NewResult(1, 0))

		err := companiesRepo.Delete(context.Background(), companyID)
		require.NotNil(t, err)
	})
}
