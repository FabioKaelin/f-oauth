package controllers

import (
	"net/http"

	"github.com/fabiokaelin/f-oauth/pkg/middleware"
	"github.com/fabiokaelin/f-oauth/pkg/password"
	user_pkg "github.com/fabiokaelin/f-oauth/pkg/user"
	"github.com/gin-gonic/gin"
)

func Password(apiGroup *gin.RouterGroup) {
	PasswordGroup := apiGroup.Group("/password")
	PasswordGroup.Use(middleware.SetUserToContext())
	{
		PasswordGroup.POST("/reset", resetPasswordPost)
		PasswordGroup.POST("/reset/:secret", resetPasswordUse)
		PasswordGroup.POST("/change", changePassword)
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
	resetPassword, err := password.CreateResetPassword(userid)

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

	err := password.UseResetPassword(secret, body.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})

}

func changePassword(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(user_pkg.User)

	body := struct {
		Password string `json:"password"`
	}{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	err := password.ChangePassword(currentUser.ID.String(), body.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
