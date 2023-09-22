package repository

import (
	"companies-service/internal/companies"
	"companies-service/internal/models"
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

const (
	createCompany = `INSERT INTO companies
						(company_name, company_description,
						 amount_of_employees, registered, company_type) 
						VALUES ($1, $2, $3, $4, $5) RETURNING *`

	updateCompany = `UPDATE companies SET
						company_name = $1, company_description = $2,
						amount_of_employees = $3, registered = $4, 
						company_type = $5
						WHERE company_id = $6 RETURNING *`

	deleteCompany = `DELETE FROM companies WHERE company_id = $1`

	getCompanyByID = `SELECT company_id, company_name, company_description, 
						amount_of_employees, registered, company_type 
						FROM companies WHERE company_id = $1`
)

// Companies Repository
type companiesRepo struct {
	db *sqlx.DB
}

// Companies Repository constructor
func NewCompaniesRepository(db *sqlx.DB) companies.Repository {
	return &companiesRepo{db: db}
}

// Create a company
func (r *companiesRepo) Create(
	ctx context.Context, company *models.Company,
) (*models.Company, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "companiesRepo.Create")
	defer span.Finish()

	c := &models.Company{}
	if err := r.db.QueryRowxContext(
		ctx,
		createCompany,
		&company.CompanyName,
		&company.CompanyDescription,
		&company.AmountOfEmployees,
		&company.Registered,
		&company.CompanyType,
	).StructScan(c); err != nil {
		return nil, errors.Wrap(err, "companiesRepo.Create.StructScan")
	}

	return c, nil
}

// Update a company
func (r *companiesRepo) Update(
	ctx context.Context, company *models.Company,
) (*models.Company, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "companiesRepo.Update")
	defer span.Finish()

	c := &models.Company{}
	if err := r.db.QueryRowxContext(
		ctx,
		updateCompany,
		company.CompanyName,
		company.CompanyDescription,
		company.AmountOfEmployees,
		company.Registered,
		company.CompanyType,
		company.CompanyID).StructScan(c); err != nil {
		return nil, errors.Wrap(err, "companiesRepo.Update.QueryRowxContext")
	}

	return c, nil
}

// Delete a company
func (r *companiesRepo) Delete(ctx context.Context, companyID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "companiesRepo.Delete")
	defer span.Finish()

	result, err := r.db.ExecContext(ctx, deleteCompany, companyID)
	if err != nil {
		return errors.Wrap(err, "companiesRepo.Delete.ExecContext")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "companiesRepo.Delete.RowsAffected")
	}

	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "companiesRepo.Delete.rowsAffected")
	}

	return nil
}

// GetByID company
func (r *companiesRepo) GetByID(
	ctx context.Context, companyID uuid.UUID,
) (*models.Company, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "companiesRepo.GetByID")
	defer span.Finish()

	company := &models.Company{}
	if err := r.db.GetContext(ctx, company, getCompanyByID, companyID); err != nil {
		return nil, errors.Wrap(err, "companiesRepo.GetByID.GetContext")
	}
	return company, nil
}
