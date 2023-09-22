package middleware

import (
	"companies-service/config"
	"companies-service/internal/auth"
	"companies-service/pkg/logger"
)

// Middleware manager
type MiddlewareManager struct {
	authService auth.Service
	cfg         *config.Config
	origins     []string
	logger      logger.Logger
}

// Middleware manager constructor
func NewMiddlewareManager(
	authService auth.Service, cfg *config.Config, origins []string, logger logger.Logger,
) *MiddlewareManager {
	return &MiddlewareManager{authService: authService, cfg: cfg, origins: origins, logger: logger}
}
