package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/google-github-oath2-golang/controllers"
	"github.com/wpcodevo/google-github-oath2-golang/initializers"
	"github.com/wpcodevo/google-github-oath2-golang/middleware"
	"github.com/wpcodevo/google-github-oath2-golang/utils"
)

var server *gin.Engine

func init() {
	startconfig, err := initializers.LoadConfig(".")
	initializers.StartConfig = startconfig
	if err != nil {
		fmt.Println("Error", err)
	}

	// initializers.ConnectDB()

	server = gin.Default()
	utils.UpdateDBConnection()

	// server.Use(CORSMiddleware())

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		origin := c.Request.Header.Get("Origin")
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080,http://localhost:3000,http://localhost:8000,http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, access-control-allow-origin, Cookie, caches, Pragma, Expires")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		// Vary: Origin
		c.Writer.Header().Set("Vary", "Origin")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	// corsConfig := cors.DefaultConfig()
	// corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173"}
	// corsConfig.AllowCredentials = true

	err := utils.UpdateDBConnection()
	if err != nil {
		newErr := errors.Join(errors.New("error durring updating db connection"), err)
		fmt.Println(newErr)
	}

	// server.Use(cors.New(corsConfig))
	server.Use(CORSMiddleware())

	router := server.Group("/api")
	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "This is the backend of the authentication server for F-Products")
	})
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Implement Google OAuth2 in Golang"})
	})

	auth_router := router.Group("/auth")
	auth_router.POST("/register", controllers.SignUpUser) // Migrated to sql
	auth_router.POST("/login", controllers.SignInUser)
	auth_router.GET("/logout", middleware.DeserializeUser(), controllers.LogoutUser)

	router.GET("/sessions/oauth/google", controllers.GoogleOAuth)
	router.GET("/sessions/oauth/github", controllers.GitHubOAuth)
	router.GET("/users/me", middleware.DeserializeUser(), controllers.GetMe)
	router.PUT("/users/me", middleware.DeserializeUser(), controllers.UpdateMe)

	router.StaticFS("/images", http.Dir("public"))
	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Route Not Found"})
	})

	log.Fatal(server.Run(":" + "8001"))
}
