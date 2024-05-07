package controllers

import (
	"github.com/gin-gonic/gin"
)

// InternalRouter defines the routes for the internal API
// it is only for Kubernetes internal communication such as health checks
func InternalRouter(router *gin.Engine) {
	internalGroup := router.Group("/internal")
	{
		internalGroup.GET("/health/live", healthLive)
		internalGroup.GET("/health/ready", healthReady)
	}
}

func healthLive(ctx *gin.Context) {
	ctx.JSON(200, "live")
}

func healthReady(ctx *gin.Context) {
	ctx.JSON(200, "ready")
}
