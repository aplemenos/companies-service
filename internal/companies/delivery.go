package companies

import (
	"github.com/gin-gonic/gin"
)

// Companies HTTP Handlers interface
type Handlers interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
}
