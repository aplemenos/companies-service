package http

import (
	"companies-service/config"
	"companies-service/internal/auth"
	"companies-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Map auth routes
func MapAuthRoutes(authGroup *gin.RouterGroup, h auth.Handlers, mw *middleware.MiddlewareManager,
	cfg *config.Config) {
	authGroup.POST("/register", h.Register)
	authGroup.POST("/login", h.Login)
	authGroup.GET("/:user_id", h.GetUserByID)
	authGroup.GET("/me", mw.AuthJWTMiddleware(cfg), h.GetMe)
	authGroup.PUT("/:user_id", mw.AuthJWTMiddleware(cfg), h.Update)
	authGroup.DELETE("/:user_id", mw.AuthJWTMiddleware(cfg), h.Delete)
}
