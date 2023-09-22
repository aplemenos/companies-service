package http

import (
	"companies-service/config"
	"companies-service/internal/companies"
	"companies-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Map companies routes
func MapCommentsRoutes(companiesGroup *gin.RouterGroup, h companies.Handlers,
	mw *middleware.MiddlewareManager, cfg *config.Config) {
	companiesGroup.POST("", mw.AuthJWTMiddleware(cfg), h.Create)
	companiesGroup.DELETE("/:company_id", mw.AuthJWTMiddleware(cfg), h.Delete)
	companiesGroup.PATCH("/:company_id", mw.AuthJWTMiddleware(cfg), h.Update)
	companiesGroup.GET("/:company_id", h.GetByID)
}
