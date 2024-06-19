package controllers

import (
	"net/http"

	"github.com/fabiokaelin/f-oauth/pkg/middleware"
	"github.com/fabiokaelin/f-oauth/pkg/resetpassword"
	user_pkg "github.com/fabiokaelin/f-oauth/pkg/user"
	"github.com/gin-gonic/gin"
)

func ResetPassword(apiGroup *gin.RouterGroup) {
	resetPasswordGroup := apiGroup.Group("/reset-password")
	resetPasswordGroup.Use(middleware.SetUserToContext())
	{
		resetPasswordGroup.POST("/", resetPasswordPost)
		resetPasswordGroup.POST("/:secret", resetPasswordUse)
	}
}

func resetPasswordPost(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(user_pkg.User)

	if currentUser.Role != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Unauthorized"})
		return
	}

	userid := ctx.Query("userid")
	if userid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User ID is required"})
		return
	}
	resetPassword, err := resetpassword.CreateResetPassword(userid)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"resetPassword": resetPassword}})
}

func resetPasswordUse(ctx *gin.Context) {
	secret := ctx.Params.ByName("secret")
	body := struct {
		Password string `json:"password"`
	}{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	err := resetpassword.UseResetPassword(secret, body.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})

}
