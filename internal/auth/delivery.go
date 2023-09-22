package auth

import "github.com/gin-gonic/gin"

// Auth HTTP Handlers interface
type Handlers interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetUserByID(c *gin.Context)
	GetMe(c *gin.Context)
}
