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
	fmt.Println("a")
	initializers.ConnectDB()

	server = gin.Default()
}

func main() {
	fmt.Println("b")
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
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