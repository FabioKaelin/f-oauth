package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fabiokaelin/f-oauth/config"
	"github.com/fabiokaelin/f-oauth/pkg/db"
	token_pkg "github.com/fabiokaelin/f-oauth/pkg/token"
	user_pkg "github.com/fabiokaelin/f-oauth/pkg/user"
	"github.com/gin-gonic/gin"
)

func SetUserToContext() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string
		cookie, err := ctx.Cookie("token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			token = fields[1]
		} else if err == nil {
			token = cookie
		}

		if token == "" {
			fmt.Println("token is empty 1")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, `gin.H{"status": "fail", "message": "You are not logged in"}1`)
			// ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		sub, err := token_pkg.ValidateToken(token, config.JWTTokenSecret)
		if err != nil {
			fmt.Println("error on validate token 2")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		var user user_pkg.User

		rows, err := db.RunSQL("SELECT `id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at` FROM `users` WHERE `id` = ? LIMIT 1;", fmt.Sprint(sub))

		if err != nil {
			fmt.Println("the user belonging to this token no logger exists 3")
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Photo, &user.Verified, &user.Provider, &user.CreatedAt, &user.UpdatedAt)
			break
		}

		// spew.Dump(user)

		ctx.Set("currentUser", user)
		ctx.Next()

	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Validate origin against allowed origins
		allowedOrigin := ""
		if origin == config.FrontEndOrigin || origin == "http://localhost:3000" || origin == "http://localhost:5173" {
			allowedOrigin = origin
		}
		
		// Only set CORS headers if origin is allowed
		if allowedOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Cookie")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
			c.Writer.Header().Set("Vary", "Origin")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func IPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/internal") {
			c.Next()
			return
		}

		ip := c.Request.Header.Get("X-Original-Forwarded-For")
		userAgent := c.Request.UserAgent()
		output := fmt.Sprintf("IP: '%s' - UserAgent: '%s'", ip, userAgent)
		fmt.Println(output)
		c.Next()
	}
}
