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
	c.IndentedJSON(http.StatusOK, "Welcome to the API for oauth.fabkli.ch")
}
