package controllers

import (
	"fmt"
	"github.com/fabiokaelin/fcommon/pkg/ratelimit"
	"net/http"

	"github.com/fabiokaelin/f-oauth/config"
	"github.com/fabiokaelin/f-oauth/pkg/db"
	"github.com/fabiokaelin/f-oauth/pkg/middleware"
	"github.com/fabiokaelin/f-oauth/pkg/notification"
	"github.com/fabiokaelin/f-oauth/pkg/password"
	"github.com/fabiokaelin/f-oauth/pkg/session"
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

	if !ratelimit.ReserveRecipient(ctx, body.Email) {
		if !ctx.IsAborted() {
			ctx.JSON(http.StatusOK, genericResponse)
		}
		return
	}
	// Lookup userObj by email — always return generic response to avoid email enumeration
	userObj, err := db.GetUserByEmail(body.Email)
	if err != nil || userObj.ID == "" {
		ctx.JSON(http.StatusOK, genericResponse)
		return
	}

	notification.ForgotPasswordNotification(user_pkg.ConvertDBUserToUser(userObj))

	if userObj.Provider != "local" {
		body := fmt.Sprintf("Hi\n\nYour account is registered with %s.\nPlease use that provider to log in. If you did not request a password reset, you can safely ignore this email.\n\nRegards,\n%s", userObj.Provider, config.EmailFromName)
		notification.SendEmail(userObj.Email, "Password Reset Request", body)

		ctx.JSON(http.StatusOK, genericResponse)
		return
	}

	// Create reset token
	resetPasswordToken, err := password.CreateResetPassword(userObj.ID)
	if err != nil {
		fmt.Println("forgotPasswordPost: failed to create reset token:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to create reset token"})
		return
	}

	// Build reset link and send email asynchronously (errors are logged, not exposed to client)
	resetLink := config.FrontEndOrigin + "/reset-password?token=" + resetPasswordToken.Secret + "&id=" + resetPasswordToken.ID
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
		Password        string `json:"password"`
		CurrentPassword string `json:"current_password"`
	}{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if currentUser.Provider != "local" || body.CurrentPassword == "" || len(body.CurrentPassword) > 72 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Current local password required"})
		return
	}
	err := password.ChangePasswordVerified(currentUser.ID.String(), body.CurrentPassword, body.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	session.ClearCookie(ctx.Writer)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
