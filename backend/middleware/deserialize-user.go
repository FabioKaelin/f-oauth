package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/google-github-oath2-golang/initializers"
	"github.com/wpcodevo/google-github-oath2-golang/models"
	"github.com/wpcodevo/google-github-oath2-golang/utils"
)

func DeserializeUser() gin.HandlerFunc {
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

		config, _ := initializers.LoadConfig(".")
		sub, err := utils.ValidateToken(token, config.JWTTokenSecret)
		if err != nil {
			fmt.Println("error on validate token 2")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		var user models.User

		rows, err := utils.RunSQL("SELECT `id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at` FROM `users` WHERE `id` = ? LIMIT 1;", fmt.Sprint(sub))

		if err != nil {
			fmt.Println("the user belonging to this token no logger exists 3")
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
			return
		}

		for rows.Next() {
			rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Photo, &user.Verified, &user.Provider, &user.CreatedAt, &user.UpdatedAt)
			break
		}

		spew.Dump(user)

		ctx.Set("currentUser", user)
		ctx.Next()

	}
}
