package service

import (
	"companies-service/config"
	"companies-service/internal/auth/mock"
	"companies-service/internal/models"
	"companies-service/pkg/logger"
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	authService := NewAuthService(cfg, mockAuthRepo, nil, apiLogger)

	user := &models.User{
		Password: "123456",
		Email:    "nick.pap@gmail.com",
	}

	ctx := context.Background()
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "authService.Register")
	defer span.Finish()

	mockAuthRepo.EXPECT().FindByEmail(ctxWithTrace, gomock.Eq(user)).Return(nil, sql.ErrNoRows)
	mockAuthRepo.EXPECT().Register(ctxWithTrace, gomock.Eq(user)).Return(user, nil)

	createdUSer, err := authService.Register(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, createdUSer)
	require.Nil(t, err)
}

func TestAuthService_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)

	mockAuthRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	authUC := NewAuthService(cfg, mockAuthRepo, mockRedisRepo, apiLogger)

	user := &models.User{
		Password: "123456",
		Email:    "nick.pap@gmail.com",
	}
	key := fmt.Sprintf("%s: %s", basePrefix, user.UserID)

	ctx := context.Background()
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "authService.Update")
	defer span.Finish()

	mockAuthRepo.EXPECT().Update(ctxWithTrace, gomock.Eq(user)).Return(user, nil)
	mockRedisRepo.EXPECT().DeleteUserCtx(ctxWithTrace, key).Return(nil)

	updatedUser, err := authUC.Update(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, updatedUser)
	require.Nil(t, err)
}

func TestAuthService_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)

	mockAuthRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	authUC := NewAuthService(cfg, mockAuthRepo, mockRedisRepo, apiLogger)

	user := &models.User{
		Password: "123456",
		Email:    "nick.pap@gmail.com",
	}
	key := fmt.Sprintf("%s: %s", basePrefix, user.UserID)

	ctx := context.Background()
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "authService.Delete")
	defer span.Finish()

	mockAuthRepo.EXPECT().Delete(ctxWithTrace, gomock.Eq(user.UserID)).Return(nil)
	mockRedisRepo.EXPECT().DeleteUserCtx(ctxWithTrace, key).Return(nil)

	err := authUC.Delete(ctx, user.UserID)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestAuthService_GetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	authUC := NewAuthService(cfg, mockAuthRepo, mockRedisRepo, apiLogger)

	user := &models.User{
		Password: "123456",
		Email:    "nick.pap@gmail.com",
	}
	key := fmt.Sprintf("%s: %s", basePrefix, user.UserID)

	ctx := context.Background()
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "authService.GetByID")
	defer span.Finish()

	mockRedisRepo.EXPECT().GetByIDCtx(ctxWithTrace, key).Return(nil, nil)
	mockAuthRepo.EXPECT().GetByID(ctxWithTrace, gomock.Eq(user.UserID)).Return(user, nil)
	mockRedisRepo.EXPECT().SetUserCtx(ctxWithTrace, key, cacheDuration, user).Return(nil)

	u, err := authUC.GetByID(ctx, user.UserID)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, u)
}

func TestAuthService_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	authUC := NewAuthService(cfg, mockAuthRepo, mockRedisRepo, apiLogger)

	ctx := context.Background()
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "authService.Login")
	defer span.Finish()

	user := &models.User{
		Password: "123456",
		Email:    "nick.pap@gmail.com",
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	require.NoError(t, err)

	mockUser := &models.User{
		Email:    "nick.pap@gmail.com",
		Password: string(hashPassword),
	}

	mockAuthRepo.EXPECT().FindByEmail(ctxWithTrace, gomock.Eq(user)).Return(mockUser, nil)

	userWithToken, err := authUC.Login(ctx, user)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, userWithToken)
}
