package repository

import (
	"companies-service/internal/models"
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestAuthRepo_Register(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("Register", func(t *testing.T) {
		gender := "male"
		role := "admin"

		rows := sqlmock.NewRows([]string{
			"first_name", "last_name", "password", "email", "role", "gender",
		}).AddRow("Nick", "Papadopoulos", "123456", "nick.pap@gmail.com", "admin", &gender)

		user := &models.User{
			FirstName: "Nick",
			LastName:  "Papadopoulos",
			Email:     "nick.pap@gmail.com",
			Password:  "123456",
			Role:      &role,
			Gender:    &gender,
		}

		mock.ExpectQuery(createUserQuery).WithArgs(&user.FirstName, &user.LastName, &user.Email,
			&user.Password, &user.Role, &user.About, &user.PhoneNumber, &user.Address, &user.City,
			&user.Gender, &user.Postcode, &user.Birthday).WillReturnRows(rows)

		createdUser, err := authRepo.Register(context.Background(), user)

		require.NoError(t, err)
		require.NotNil(t, createdUser)
		require.Equal(t, createdUser, user)
	})
}

func TestAuthRepo_GetByID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("GetByID", func(t *testing.T) {
		uuid := uuid.New()

		rows := sqlmock.NewRows([]string{
			"user_id", "first_name", "last_name", "email",
		}).AddRow(uuid, "Nick", "Papadopoulos", "nick.pap@gmail.com")

		testUser := &models.User{
			UserID:    uuid,
			FirstName: "Nick",
			LastName:  "Papadopoulos",
			Email:     "nick.pap@gmail.com",
		}

		mock.ExpectQuery(getUserQuery).
			WithArgs(uuid).
			WillReturnRows(rows)

		user, err := authRepo.GetByID(context.Background(), uuid)
		require.NoError(t, err)
		require.Equal(t, user.FirstName, testUser.FirstName)
		fmt.Printf("test user: %s \n", testUser.FirstName)
		fmt.Printf("user: %s \n", user.FirstName)
	})
}

func TestAuthRepo_Delete(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("Delete", func(t *testing.T) {
		uid := uuid.New()
		mock.ExpectExec(deleteUserQuery).WithArgs(uid).WillReturnResult(sqlmock.NewResult(1, 1))

		err := authRepo.Delete(context.Background(), uid)
		require.Nil(t, err)
	})

	t.Run("Delete No rows", func(t *testing.T) {
		uid := uuid.New()
		mock.ExpectExec(deleteUserQuery).WithArgs(uid).WillReturnResult(sqlmock.NewResult(1, 0))

		err := authRepo.Delete(context.Background(), uid)
		require.NotNil(t, err)
	})
}

func TestAuthRepo_Update(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("Update", func(t *testing.T) {
		gender := "male"
		role := "admin"

		rows := sqlmock.NewRows([]string{
			"first_name", "last_name", "password", "email", "role", "gender",
		}).AddRow("Nick", "Papadopoulos", "123456", "nick.pap@gmail.com", "admin", &gender)

		user := &models.User{
			FirstName: "Nick",
			LastName:  "Papadopoulos",
			Email:     "nick.pap@gmail.com",
			Password:  "123456",
			Role:      &role,
			Gender:    &gender,
		}

		mock.ExpectQuery(updateUserQuery).WithArgs(&user.FirstName, &user.LastName, &user.Email,
			&user.Role, &user.About, &user.PhoneNumber, &user.Address, &user.City, &user.Gender,
			&user.Postcode, &user.Birthday, &user.UserID).WillReturnRows(rows)

		updatedUser, err := authRepo.Update(context.Background(), user)

		require.NoError(t, err)
		require.NotNil(t, updatedUser)
		require.Equal(t, user, updatedUser)
	})
}

func TestAuthRepo_FindByEmail(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("FindByEmail", func(t *testing.T) {
		uuid := uuid.New()

		rows := sqlmock.NewRows([]string{
			"user_id", "first_name", "last_name", "email",
		}).AddRow(uuid, "Nick", "Papadopoulos", "nick.pap@gmail.com")

		testUser := &models.User{
			UserID:    uuid,
			FirstName: "Nick",
			LastName:  "Papadopoulos",
			Email:     "nick.pap@mail.com",
		}

		mock.ExpectQuery(findUserByEmail).WithArgs(testUser.Email).WillReturnRows(rows)

		foundUser, err := authRepo.FindByEmail(context.Background(), testUser)

		require.NoError(t, err)
		require.NotNil(t, foundUser)
		require.Equal(t, foundUser.FirstName, testUser.FirstName)
	})
}
