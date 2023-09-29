package server

import (
	docs "companies-service/docs"
	authHttp "companies-service/internal/auth/delivery/http"
	authRepository "companies-service/internal/auth/repository"
	authService "companies-service/internal/auth/service"
	companiesHttp "companies-service/internal/companies/delivery/http"
	companiesRepository "companies-service/internal/companies/repository"
	companiesService "companies-service/internal/companies/service"
	"companies-service/internal/middleware"
	"companies-service/pkg/kafka"
	"companies-service/pkg/metric"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Map Server Handlers
func (s *Server) MapHandlers(kafkaProducer kafka.Producer) error {
	metrics, err := metric.CreateMetrics(s.cfg.Metrics.URL, s.cfg.Metrics.ServiceName)
	if err != nil {
		s.logger.Errorf("CreateMetrics Error: %s", err)
	}
	s.logger.Info(
		"Metrics available URL: %s, ServiceName: %s",
		s.cfg.Metrics.URL,
		s.cfg.Metrics.ServiceName,
	)

	// Init repositories
	authRepo := authRepository.NewAuthRepository(s.db)
	companiesRepo := companiesRepository.NewCompaniesRepository(s.db)
	authRedisRepo := authRepository.NewAuthRedisRepo(s.redisClient)

	// Init useCases
	authSrv := authService.NewAuthService(s.cfg, authRepo, authRedisRepo, s.logger)
	companiesSrv := companiesService.NewCompaniesService(s.cfg, companiesRepo, s.logger)

	// Init handlers
	authHandler := authHttp.NewAuthHandlers(s.cfg, authSrv, s.logger)
	companiesHandler := companiesHttp.NewCompaniesHandlers(s.cfg, companiesSrv, kafkaProducer,
		s.logger)

	docs.SwaggerInfo.Title = "Company Service REST API"
	s.gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	s.gin.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST, OPTIONS, GET, PUT, PATCH, DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter
	// even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	s.gin.Use(gin.Logger())
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	s.gin.Use(gin.Recovery())

	s.gin.Use(requestid.New())

	mw := middleware.NewMiddlewareManager(authSrv, s.cfg, []string{"*"}, s.logger)
	s.gin.Use(mw.MetricsMiddleware(metrics))

	s.gin.Use(limits.RequestSizeLimiter(1024 * 1024 * 5)) // 5MB
	if s.cfg.Server.Debug {
		s.gin.Use(mw.DebugMiddleware())
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	v1 := s.gin.Group("/api/v1")

	authGroup := v1.Group("/auth")
	companiesGroup := v1.Group("/companies")

	authHttp.MapAuthRoutes(authGroup, authHandler, mw, s.cfg)
	companiesHttp.MapCommentsRoutes(companiesGroup, companiesHandler, mw, s.cfg)

	return nil
}
