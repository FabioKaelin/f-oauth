package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/fabiokaelin/f-oauth/pkg/auth"
	"github.com/fabiokaelin/f-oauth/pkg/middleware"
	"github.com/fabiokaelin/f-oauth/pkg/session"
	user_pkg "github.com/fabiokaelin/f-oauth/pkg/user"

	"github.com/gin-gonic/gin"
)

func AuthRouter(apiGroup *gin.RouterGroup) {
	oauthGroup := apiGroup.Group("/auth")
	{
		oauthGroup.POST("/register", authRegister)
		oauthGroup.POST("/login", authLogin)
		oauthGroup.POST("/logout", middleware.SetUserToContext(), authLogout)
	}
}

// authRegister
func authRegister(ctx *gin.Context) {
	// TODO: Redirect to loginpage when an error occurs with error message
	var payload user_pkg.RegisterUserInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		fmt.Println("1", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "1message": err.Error()})
		return
	}

	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	emailRegex := regexp.MustCompile(emailPattern)

	if !emailRegex.MatchString(payload.Email) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email address pattern"})
		return
	}

	usernamePattern := "^[a-zA-Z0-9_\\-öäüÖÄÜêÊéàèÉÀÈç ]{3,20}$"
	usernameRegesxx := regexp.MustCompile(usernamePattern)

	if !usernameRegesxx.MatchString(payload.Name) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid username pattern"})
		return
	}

	// hashpassword := sha512.Sum512([]byte(payload.Password))
	// hashpasswordValue := hex.EncodeToString(hashpassword[:])
	if err := auth.ValidateNewPassword(payload.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	hashPassword, err := auth.HashAndSalt([]byte(payload.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "error in hashing password"})
		return
	}

	now := time.Now()
	newUser := user_pkg.User{
		Name:      payload.Name,
		Email:     strings.ToLower(payload.Email),
		Password:  hashPassword,
		Role:      "user",
		Provider:  "local",
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	newUser, err = auth.RegisterUser(newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": user_pkg.FilteredResponse(&newUser)}})
}

// authLogin
func authLogin(ctx *gin.Context) {
	// TODO: Redirect to loginpage when an error occurs with error message
	fmt.Println("try to login")
	var payload user_pkg.LoginUserInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		fmt.Println("1", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	fmt.Println("user", payload.Email)

	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	emailRegex := regexp.MustCompile(emailPattern)

	if !emailRegex.MatchString(payload.Email) {
		fmt.Println("invalid email address pattern", payload.Email)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email address pattern"})
		return
	}

	user, err := auth.LoginUser(payload.Email, payload.Password)
	if err != nil {
		fmt.Println("2", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	token, err := session.Create(user.ID.String(), user.Password)
	if err != nil {
		fmt.Println("3", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	session.SetCookie(ctx.Writer, token)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "token": token})
}

// authLogout
func authLogout(ctx *gin.Context) {
	token := ctx.GetString("sessionToken")
	if err := session.Revoke(token); err != nil {
		ctx.JSON(503, gin.H{"status": "fail", "message": "Session revocation unavailable"})
		return
	}
	session.ClearCookie(ctx.Writer)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
