package main

import (
	"fmt"
	"time"

	"github.com/fabiokaelin/f-oauth/config"
	"github.com/fabiokaelin/f-oauth/controllers"
	"github.com/fabiokaelin/f-oauth/pkg/db"
	"github.com/fabiokaelin/f-oauth/pkg/logger"
	"github.com/fabiokaelin/f-oauth/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	err := config.Load()

	if err != nil {
		fmt.Println("Error in Load config", err)
		panic(err)
	}

	err = db.UpdateDBConnection()
	if err != nil {
		time.Sleep(10 * time.Second)
		err = db.UpdateDBConnection()
		if err != nil {
			time.Sleep(10 * time.Second)
			err = db.UpdateDBConnection()
			if err != nil {
				panic(err)
			}
		}
	}

	gin.SetMode(config.GinMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithConfig(logger.LoggerConfig))
	router.Use(middleware.IPMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.ForwardedByClientIP = true
	router.HandleMethodNotAllowed = true
	// router.SetTrustedProxies([]string{"127.0.0.1", "localhost", "oauth.fabkli.ch"})
	router.GET("", controllers.Default)

	controllers.InternalRouter(router)

	apiGroup := router.Group("/api")

	apiGroup.GET("", controllers.Default)
	apiGroup.GET("/version", func(c *gin.Context) {
		version := config.FVersion
		c.IndentedJSON(200, gin.H{"version": version})
	})

	controllers.AuthRouter(apiGroup)
	controllers.OAuth2Router(apiGroup)
	controllers.UserRouter(apiGroup)
	controllers.ResetPassword(apiGroup)

	// router.NoRoute(func(ctx *gin.Context) {
	// ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Route Not Found"})
	// })

	fmt.Println("Server is running")
	router.Run("0.0.0.0:8001")
}
