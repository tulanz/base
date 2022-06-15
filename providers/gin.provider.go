package providers

import (
	"github.com/gin-gonic/gin"
)

func NewGinProvider() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin := gin.New()
	return gin
}
