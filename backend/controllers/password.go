package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fabiokaelin/f-oauth/config"
	"github.com/fabiokaelin/f-oauth/pkg/db"
	"github.com/fabiokaelin/f-oauth/pkg/middleware"
	"github.com/fabiokaelin/f-oauth/pkg/notification"
	"github.com/fabiokaelin/f-oauth/pkg/password"
	user_pkg "github.com/fabiokaelin/f-oauth/pkg/user"
	"github.com/gin-gonic/gin"
)

func Password(apiGroup *gin.RouterGroup) {
	PasswordGroup := apiGroup.Group("/password")
	{
		PasswordGroup.POST("/forgot", forgotPasswordPost)
		PasswordGroup.POST("/reset", middleware.SetUserToContext(), resetPasswordPost)
		PasswordGroup.POST("/reset/:secret", resetPasswordUse)
		PasswordGroup.POST("/change", middleware.SetUserToContext(), changePassword)
	}
}

func forgotPasswordPost(ctx *gin.Context) {
	body := struct {
		Email string `json:"email"`
	}{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	genericResponse := gin.H{"status": "success", "message": "If that email is registered, you will receive a reset link shortly."}

	// Lookup user by email — always return generic response to avoid email enumeration
	user, err := db.GetUserByEmail(body.Email)
	if err != nil || user.ID == "" {
		ctx.JSON(http.StatusOK, genericResponse)
		return
	}

	// Rate limiting: max 3 reset requests per user per hour
	oneHourAgo := time.Now().Add(-time.Hour)
	count, err := db.CountRecentResetPasswordsByUserID(user.ID, oneHourAgo)
	if err == nil && count >= 3 {
		ctx.JSON(http.StatusOK, genericResponse)
		return
	}

	// Create reset token
	resetPasswordToken, err := password.CreateResetPassword(user.ID)
	if err != nil {
		fmt.Println("forgotPasswordPost: failed to create reset token:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to create reset token"})
		return
	}

	// Build reset link and send email asynchronously (errors are logged, not exposed to client)
	resetLink := config.FrontendResetPasswordURL + "?token=" + resetPasswordToken.Secret + "&id=" + resetPasswordToken.ID
	go func() {
		if err := notification.SendPasswordResetEmail(body.Email, resetLink); err != nil {
			fmt.Println("forgotPasswordPost: failed to send reset email:", err)
		}
	}()

	ctx.JSON(http.StatusOK, genericResponse)
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
		status := http.StatusBadRequest
		if err.Error() == "token expired" || err.Error() == "token already used" {
			status = http.StatusGone
		}
		ctx.JSON(status, gin.H{"status": "fail", "message": err.Error()})
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
