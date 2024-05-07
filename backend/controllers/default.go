package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Default godoc
//
//	@Summary		Default
//	@Description	Default
//	@Tags			default
//	@Produce		json
//	@Success		200	{string}	string
//	@Router			/ [get]
//	@Router			/api [get]
func Default(c *gin.Context) {
	// notificationConfig := notification.Config{
	// 	Title:   "Welcome to the API for oauth.fabkli.ch",
	// 	Message: "Welcome to the API for oauth.fabkli.ch",
	// 	Type:    "oauthdefault",
	// }
	// err := notificationConfig.Send()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	c.IndentedJSON(http.StatusOK, "Welcome to the API for oauth.fabkli.ch")
}
