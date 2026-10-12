package main

import (
	"fmt"
	"net/http"
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

	if err := db.EnsureSessions(); err != nil {
		panic(err)
	}

	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if err := db.DeleteExpiredSessions(); err != nil {
				fmt.Println("Session cleanup unavailable")
			}
		}
	}()
	gin.SetMode(config.GinMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithConfig(logger.LoggerConfig))
	router.Use(middleware.IPMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ProtectWrites())
	router.ForwardedByClientIP = false
	if err := router.SetTrustedProxies(nil); err != nil {
		panic(err)
	}
	router.Use(func(c *gin.Context) { c.Header("X-Content-Type-Options", "nosniff"); c.Next() })
	router.HandleMethodNotAllowed = true
	// router.SetTrustedProxies([]string{"127.0.0.1", "localhost", "oauth.fabkli.ch"})
	router.GET("", controllers.Default)

	controllers.InternalRouter(router)

	apiGroup := router.Group("/api")
	apiGroup.Use(func(c *gin.Context) { c.Header("Cache-Control", "no-store"); c.Next() })

	apiGroup.GET("", controllers.Default)
	apiGroup.GET("/version", func(c *gin.Context) {
		version := config.FVersion
		c.IndentedJSON(200, gin.H{"version": version})
	})

	controllers.AuthRouter(apiGroup)
	controllers.OAuth2Router(apiGroup)
	controllers.UserRouter(apiGroup)
	controllers.Password(apiGroup)

	// router.NoRoute(func(ctx *gin.Context) {
	// ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Route Not Found"})
	// })

	fmt.Println("Server is running")
	httpServer := &http.Server{Addr: "0.0.0.0:8001", Handler: router, ReadHeaderTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 32768}
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
