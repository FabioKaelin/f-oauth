package main

// API with gin

import (
	"errors"
	"github.com/fabiokaelin/fcommon/pkg/oauthsession"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	client := oauthsession.NewClient()

	server.GET("", func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("token")
		if err != nil {
			ctx.JSON(200, "no token")
			return
		}
		result, err := client.Validate(ctx.Request.Context(), "https://oauth.fabkli.ch/api/users/me", cookie)
		if err != nil {
			status := http.StatusServiceUnavailable
			if errors.Is(err, oauthsession.ErrUnauthenticated) {
				status = http.StatusUnauthorized
			}
			ctx.AbortWithStatus(status)
			return
		}
		result.ForwardCookies(ctx.Writer)
		responseJson := result.User

		responseBody := ""
		responseBody += "ID: " + responseJson.ID + "\n"
		responseBody += "Name: " + responseJson.Name + "\n"
		responseBody += "Email: " + responseJson.Email + "\n"
		ctx.JSON(200, responseBody)
	})

	router := server.Group("/api")
	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(200, "test of f-oauth")
	})

	server.Run(":8000")

}

type UserResponse struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Role     string `json:"role,omitempty"`
	Provider string `json:"provider,omitempty"`
	Photo    string `json:"photo,omitempty"`
	Verified bool   `json:"verified,omitempty"`
	// CreatedAt time.Time
	// UpdatedAt time.Time
}

// on / reads token and send request to oauth server with token as header and display user info
