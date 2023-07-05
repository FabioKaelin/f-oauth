package main

// API with gin

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	server.GET("", func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("token")
		if err != nil {
			ctx.JSON(200, "no token")
			return
		}
		client := http.Client{}
		req, err := http.NewRequest("GET", "https://oauth.fabkli.ch/api/users/me", nil)
		if err != nil {
			ctx.JSON(200, "error"+err.Error())
			return
		}

		req.Header = http.Header{
			// "Host":          {"www.host.com"},
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer " + cookie},
		}

		res, err := client.Do(req)
		if err != nil {
			ctx.JSON(200, "error"+err.Error())
			return
		}
		responseJson := UserResponse{}
		json.NewDecoder(res.Body).Decode(&responseJson)

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
