package repository

import (
	"companies-service/internal/auth"
	"companies-service/internal/models"
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

const (
	createUserQuery = `INSERT INTO users (first_name, last_name, email, password, role, about, phone_number, address,
	               		city, gender, postcode, birthday, created_at, updated_at, login_date)
						VALUES ($1, $2, $3, $4, COALESCE(NULLIF($5, ''), 'user'), $6, $7, $8, $9, $10, $11, $12, now(), now(), now()) 
						RETURNING *`

	updateUserQuery = `UPDATE users 
						SET first_name = COALESCE(NULLIF($1, ''), first_name),
						    last_name = COALESCE(NULLIF($2, ''), last_name),
						    email = COALESCE(NULLIF($3, ''), email),
						    role = COALESCE(NULLIF($4, ''), role),
						    about = COALESCE(NULLIF($5, ''), about),
						    phone_number = COALESCE(NULLIF($7, ''), phone_number),
						    address = COALESCE(NULLIF($8, ''), address),
						    city = COALESCE(NULLIF($9, ''), city),
						    gender = COALESCE(NULLIF($10, ''), gender),
						    postcode = COALESCE(NULLIF($11, 0), postcode),
						    birthday = COALESCE(NULLIF($12, '')::date, birthday),
						    updated_at = now()
						WHERE user_id = $13
						RETURNING *
						`

	deleteUserQuery = `DELETE FROM users WHERE user_id = $1`

	getUserQuery = `SELECT user_id, first_name, last_name, email, role, about, phone_number, 
       				 address, city, gender, postcode, birthday, created_at, updated_at, login_date  
					 FROM users 
					 WHERE user_id = $1`

	findUserByEmail = `SELECT user_id, first_name, last_name, email, role, about, phone_number, 
       			 		address, city, gender, postcode, birthday, created_at, updated_at, login_date, password
				 		FROM users 
				 		WHERE email = $1`
)

// Auth Repository
type authRepo struct {
	db *sqlx.DB
}

// Auth Repository constructor
func NewAuthRepository(db *sqlx.DB) auth.Repository {
	return &authRepo{db: db}
}

// Register creates new user
func (r *authRepo) Register(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authPGRepo.Register")
	defer span.Finish()

	u := &models.User{}
	if err := r.db.QueryRowxContext(ctx, createUserQuery, &user.FirstName, &user.LastName,
		&user.Email, &user.Password, &user.Role, &user.About, &user.PhoneNumber, &user.Address,
		&user.City, &user.Gender, &user.Postcode, &user.Birthday,
	).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "authPGRepo.Register.StructScan")
	}

	return u, nil
}

// Update updates existing user
func (r *authRepo) Update(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authPGRepo.Update")
	defer span.Finish()

	u := &models.User{}
	if err := r.db.GetContext(ctx, u, updateUserQuery, &user.FirstName, &user.LastName,
		&user.Email, &user.Role, &user.About, &user.PhoneNumber, &user.Address, &user.City,
		&user.Gender, &user.Postcode, &user.Birthday, &user.UserID,
	); err != nil {
		return nil, errors.Wrap(err, "authPGRepo.Update.GetContext")
	}

	return u, nil
}

// Delete removes existing user
func (r *authRepo) Delete(ctx context.Context, userID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authPGRepo.Delete")
	defer span.Finish()

	result, err := r.db.ExecContext(ctx, deleteUserQuery, userID)
	if err != nil {
		return errors.WithMessage(err, "authPGRepo Delete ExecContext")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "authPGRepo.Delete.RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "authPGRepo.Delete.rowsAffected")
	}

	return nil
}

// GetByID retreives user by id
func (r *authRepo) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authPGRepo.GetByID")
	defer span.Finish()

	user := &models.User{}
	if err := r.db.QueryRowxContext(ctx, getUserQuery, userID).StructScan(user); err != nil {
		return nil, errors.Wrap(err, "authPGRepo.GetByID.QueryRowxContext")
	}
	return user, nil
}

// FindByEmail searches user by email
func (r *authRepo) FindByEmail(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authPGRepo.FindByEmail")
	defer span.Finish()

	foundUser := &models.User{}
	if err := r.db.QueryRowxContext(ctx, findUserByEmail,
		user.Email).StructScan(foundUser); err != nil {
		return nil, errors.Wrap(err, "authPGRepo.FindByEmail.QueryRowxContext")
	}
	return foundUser, nil
}
