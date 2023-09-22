package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

// Debug dump request middleware
func (mw *MiddlewareManager) DebugMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if mw.cfg.Server.Debug {
			dump, err := httputil.DumpRequest(c.Request, true)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				c.Abort()
				return
			}
			mw.logger.Info(fmt.Sprintf(
				"\nRequest dump begin :--------------\n\n%s\n\nRequest dump end :--------------",
				dump))
		}

		c.Next()
	}
}
