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
	"companies-service/pkg/httphelper"
	"companies-service/pkg/kafka"
	"companies-service/pkg/metric"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Map Server Handlers
func (s *Server) MapHandlers(r *gin.Engine, kafkaProducer kafka.Producer) error {
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
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Use(cors.New(cors.Config{
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
	r.Use(gin.Logger())
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	r.Use(requestid.New())

	mw := middleware.NewMiddlewareManager(authSrv, s.cfg, []string{"*"}, s.logger)
	r.Use(mw.MetricsMiddleware(metrics))

	r.Use(limits.RequestSizeLimiter(1024 * 1024 * 5)) // 5MB
	if s.cfg.Server.Debug {
		r.Use(mw.DebugMiddleware())
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	v1 := r.Group("/api/v1")

	health := v1.Group("/health")
	authGroup := v1.Group("/auth")
	companiesGroup := v1.Group("/companies")

	authHttp.MapAuthRoutes(authGroup, authHandler, mw, s.cfg)
	companiesHttp.MapCommentsRoutes(companiesGroup, companiesHandler, mw, s.cfg)

	health.GET("", func(c *gin.Context) {
		s.logger.Infof("Health check RequestID: %s", httphelper.GetRequestID(c))
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	return nil
}
