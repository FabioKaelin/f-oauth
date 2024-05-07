package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fabiokaelin/f-oauth/config"
	"github.com/fabiokaelin/f-oauth/pkg/auth"
	"github.com/fabiokaelin/f-oauth/pkg/db"
	"github.com/fabiokaelin/f-oauth/pkg/middleware"
	"github.com/fabiokaelin/f-oauth/pkg/notification"
	token_pkg "github.com/fabiokaelin/f-oauth/pkg/token"
	user_pkg "github.com/fabiokaelin/f-oauth/pkg/user"

	"github.com/gin-gonic/gin"
)

func AuthRouter(apiGroup *gin.RouterGroup) {
	oauthGroup := apiGroup.Group("/auth")
	{
		oauthGroup.POST("/register", authRegister)
		oauthGroup.POST("/login", authLogin)
		oauthGroup.GET("/logout", middleware.SetUserToContext(), authLogout)
	}
}

// authRegister
func authRegister(ctx *gin.Context) {
	// TODO: Redirect to loginpage when an error occurs with error message
	var payload *user_pkg.RegisterUserInput

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
	hashPassword, err := auth.HashAndSalt([]byte(payload.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "error in hashing password"})
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

	rows, err := db.RunSQL("SELECT `id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at` FROM `users` WHERE `email` = ? LIMIT 1", strings.ToLower(payload.Email))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "error appeared"})
		return
	}
	defer rows.Close()

	isExisting := rows.Next()

	if isExisting {
		// TODO: Redirect to loginpage with error message
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
		return
	}

	row, err := db.RunSQLRow("INSERT INTO `users`(`id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at`) VALUES (UUID(), ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id ;", newUser.Name, newUser.Email, newUser.Password, newUser.Role, newUser.Photo, newUser.Verified, newUser.Provider, newUser.CreatedAt, newUser.UpdatedAt)

	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "2message": "User with that email already exists"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "3message": "Something bad happened"})
		return
	}

	err = row.Scan(&newUser.ID)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "4message": "Something bad happened"})
		return
	}

	notificationConfig := notification.Config{
		Title:   fmt.Sprintf("New User: %s", newUser.Name),
		Message: fmt.Sprintf("Provider: %s\nEmail: %s", newUser.Provider, newUser.Email),
		Type:    "newuser",
	}

	err = notificationConfig.Send()
	if err != nil {
		fmt.Println(err)
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": user_pkg.FilteredResponse(&newUser)}})
}

// authLogin
func authLogin(ctx *gin.Context) {
	// TODO: Redirect to loginpage when an error occurs with error message
	var payload *user_pkg.LoginUserInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		spew.Dump(payload)
		fmt.Println("1", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	emailRegex := regexp.MustCompile(emailPattern)

	if !emailRegex.MatchString(payload.Email) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email address pattern"})
		return
	}

	var user user_pkg.User

	rows, err := db.RunSQL("SELECT `id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at` FROM `users` WHERE `email` = ? LIMIT 1", strings.ToLower(payload.Email))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}
	defer rows.Close()
	isExisting := rows.Next()

	if !isExisting {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Photo, &user.Verified, &user.Provider, &user.CreatedAt, &user.UpdatedAt)

	equal := auth.ComparePasswords(user.Password, []byte(payload.Password))

	if !equal {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	if user.Provider != "local" {
		// TODO: Redirect to loginpage with error message
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You have already signed up with a different method"})
		return
	}

	token, err := token_pkg.GenerateToken(config.TokenExpiresIn, user.ID, config.JWTTokenSecret)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("token", token, config.TokenMaxAge*60, "/", config.TokenURL, false, true)
	ctx.SetCookie("token", token, config.TokenMaxAge*60, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

// authLogout
func authLogout(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", config.TokenURL, false, true)
	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
