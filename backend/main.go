package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/google-github-oath2-golang/controllers"
	"github.com/wpcodevo/google-github-oath2-golang/initializers"
	"github.com/wpcodevo/google-github-oath2-golang/middleware"
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

	// server.Use(CORSMiddleware())

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080,http://localhost:3000,http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, access-control-allow-origin, Cookie, caches, Pragma, Expires")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

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

	router.StaticFS("/images", http.Dir("public"))
	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Route Not Found"})
	})

	log.Fatal(server.Run(":" + "8000"))
}
