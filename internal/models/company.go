package models

import (
	"github.com/google/uuid"
)

// Company base model
type Company struct {
	CompanyID          uuid.UUID `json:"company_id" db:"company_id" validate:"omitempty,uuid"`
	CompanyName        string    `json:"company_name" db:"company_name" validate:"required,lte=15"`
	CompanyDescription string    `json:"company_description" db:"company_description" validate:"omitempty,lte=3000"`
	AmountOfEmployees  int       `json:"amount_of_employees" db:"amount_of_employees" validate:"required"`
	Registered         bool      `json:"registered" db:"registered" validate:"required"`
	CompanyType        string    `json:"company_type" db:"company_type" validate:"required"`
}
