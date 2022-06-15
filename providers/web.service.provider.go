package providers

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func StartGinWebService(lifecycle fx.Lifecycle, gin *gin.Engine) {

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return gin.Run(":4008")
		},
	})
}
