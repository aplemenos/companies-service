package server

import (
	"companies-service/pkg/httphelper"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) runHealthCheck() {
	s.gin.GET("/live", func(c *gin.Context) {
		s.logger.Infof("Health check RequestID: %s", httphelper.GetRequestID(c))
		if err := s.configureHealthCheckEndpoints(c); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})
}

func (s *Server) configureHealthCheckEndpoints(ctx context.Context) error {
	// Liveness check of (Postgres) DB
	if err := s.db.Ping(); err != nil {
		s.logger.Warnf("(DB Liveness Check) err: {%v}", err)
		return err
	}

	// Liveness check of Redis
	if err := s.redisClient.Ping(ctx).Err(); err != nil {
		s.logger.Warnf("(Redis Liveness Check) err: {%v}", err)
		return err
	}

	return nil
}
