package middleware

import (
	"companies-service/config"
	"companies-service/pkg/authn"
	"companies-service/pkg/httphelper"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// JWT way of auth using cookie or Authorization header
func (mw *MiddlewareManager) AuthJWTMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := authn.ExtractBearerToken(c.Request)

		mw.logger.Infof("auth middleware bearerHeader %s", tokenString)

		if tokenString != "" {
			if err := mw.validateJWTToken(tokenString, cfg); err != nil {
				mw.logger.Error("middleware validateJWTToken",
					zap.String("headerJWT", err.Error()))
				c.String(http.StatusUnauthorized, httphelper.ErrUnauthorized.Error())
				c.Abort()
				return
			}

			c.Next()
			return
		}

		cookie, err := c.Cookie("jwt-token")
		if err != nil {
			mw.logger.Errorf("c.Cookie", err.Error())
			c.String(http.StatusUnauthorized, httphelper.ErrUnauthorized.Error())
			c.Abort()
			return
		}

		if err = mw.validateJWTToken(cookie, cfg); err != nil {
			mw.logger.Errorf("validateJWTToken", err.Error())
			c.String(http.StatusUnauthorized, httphelper.ErrUnauthorized.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}

func (mw *MiddlewareManager) validateJWTToken(tokenString string, cfg *config.Config) error {
	if tokenString == "" {
		return httphelper.ErrInvalidJWTToken
	}

	_, err := authn.ValidateJWT(tokenString, cfg)
	if err != nil {
		return err
	}

	return nil
}
