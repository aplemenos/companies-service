package service

import (
	"companies-service/config"
	"companies-service/internal/auth"
	"companies-service/internal/models"
	"companies-service/pkg/authn"
	"companies-service/pkg/httphelper"
	"companies-service/pkg/logger"
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

const (
	basePrefix    = "api-auth:"
	cacheDuration = 3600
)

// Auth Service
type authService struct {
	cfg       *config.Config
	authRepo  auth.Repository
	redisRepo auth.RedisRepository
	logger    logger.Logger
}

// Auth Service constructor
func NewAuthService(
	cfg *config.Config,
	authRepo auth.Repository,
	redisRepo auth.RedisRepository,
	log logger.Logger,
) auth.Service {
	return &authService{cfg: cfg, authRepo: authRepo, redisRepo: redisRepo, logger: log}
}

// Create new user
func (s *authService) Register(
	ctx context.Context, user *models.User,
) (*models.UserWithToken, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authService.Register")
	defer span.Finish()

	existsUser, err := s.authRepo.FindByEmail(ctx, user)
	if existsUser != nil || err == nil {
		return nil,
			httphelper.NewRestErrorWithMessage(http.StatusBadRequest,
				httphelper.ErrEmailAlreadyExists, nil)
	}

	if err = user.PrepareCreate(); err != nil {
		return nil,
			httphelper.NewBadRequestError(errors.Wrap(err, "authService.Register.PrepareCreate"))
	}

	createdUser, err := s.authRepo.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	createdUser.SanitizePassword()

	token, err := authn.GenerateJWTToken(createdUser, s.cfg)
	if err != nil {
		return nil,
			httphelper.NewInternalServerError(errors.Wrap(err,
				"authService.Register.GenerateJWTToken"))
	}

	return &models.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

// Update existing user
func (s *authService) Update(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authService.Update")
	defer span.Finish()

	if err := user.PrepareUpdate(); err != nil {
		return nil,
			httphelper.NewBadRequestError(errors.Wrap(err, "authService.Register.PrepareUpdate"))
	}

	updatedUser, err := s.authRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	updatedUser.SanitizePassword()

	if err = s.redisRepo.DeleteUserCtx(ctx, s.GenerateUserKey(user.UserID.String())); err != nil {
		s.logger.Errorf("AuthService.Update.DeleteUserCtx: %s", err)
	}

	return updatedUser, nil
}

// Delete user
func (s *authService) Delete(ctx context.Context, userID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authService.Delete")
	defer span.Finish()

	if err := s.authRepo.Delete(ctx, userID); err != nil {
		return err
	}

	if err := s.redisRepo.DeleteUserCtx(ctx, s.GenerateUserKey(userID.String())); err != nil {
		s.logger.Errorf("AuthService.Delete.DeleteUserCtx: %s", err)
	}

	return nil
}

// Get user by id
func (u *authService) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authService.GetByID")
	defer span.Finish()

	cachedUser, err := u.redisRepo.GetByIDCtx(ctx, u.GenerateUserKey(userID.String()))
	if err != nil {
		u.logger.Errorf("authService.GetByID.GetByIDCtx: %v", err)
	}
	if cachedUser != nil {
		return cachedUser, nil
	}

	user, err := u.authRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.SetUserCtx(ctx,
		u.GenerateUserKey(userID.String()), cacheDuration, user); err != nil {
		u.logger.Errorf("authService.GetByID.SetUserCtx: %v", err)
	}

	user.SanitizePassword()

	return user, nil
}

// Login user, returns user model with jwt token
func (s *authService) Login(
	ctx context.Context, user *models.User,
) (*models.UserWithToken, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authService.Login")
	defer span.Finish()

	foundUser, err := s.authRepo.FindByEmail(ctx, user)
	if err != nil {
		return nil, err
	}

	if err = foundUser.ComparePasswords(user.Password); err != nil {
		return nil,
			httphelper.NewUnauthorizedError(errors.Wrap(err,
				"authService.GetUsers.ComparePasswords"))
	}

	foundUser.SanitizePassword()

	token, err := authn.GenerateJWTToken(foundUser, s.cfg)
	if err != nil {
		return nil,
			httphelper.NewInternalServerError(errors.Wrap(err,
				"authService.GetUsers.GenerateJWTToken"))
	}

	return &models.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}

func (s *authService) GenerateUserKey(userID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, userID)
}
