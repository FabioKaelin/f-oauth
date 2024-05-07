package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fabiokaelin/f-oauth/config"
	"github.com/fabiokaelin/f-oauth/pkg/db"
	"github.com/fabiokaelin/f-oauth/pkg/google"
	"github.com/fabiokaelin/f-oauth/pkg/image"
	token_pkg "github.com/fabiokaelin/f-oauth/pkg/token"
	user_pkg "github.com/fabiokaelin/f-oauth/pkg/user"

	"github.com/fabiokaelin/f-oauth/pkg/github"
	"github.com/gin-gonic/gin"
)

func OAuth2Router(apiGroup *gin.RouterGroup) {
	oauth2Group := apiGroup.Group("/sessions/oauth")
	{
		oauth2Group.GET("/google", oauth2Google)
		oauth2Group.GET("/github", oauth2GitHub)
	}
}

func oauth2Google(ctx *gin.Context) {
	fmt.Println(ctx.Request.URL.String())
	// TODO: Redirect to loginpage when an error occurs with error message
	code := ctx.Query("code")
	// ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s%s", config.FrontEndOrigin, "/login?error=already_signed_up_with_different_method"))
	// return
	// ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You have already signed up with a different method"})
	state := ctx.Query("state") // path/url to redirect after login

	if code == "" {
		ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s%s", config.FrontEndOrigin, "/login"))
		return
	}

	token, errorcode, err := google.LoginWithCode(code)

	if err != nil {
		fmt.Println("error, failed login with google:", err)
		if errorcode == 2 {
			ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s%s", config.FrontEndOrigin, "/login?error=already_signed_up_with_different_method"))
			return
			// ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You have already signed up with a different method"})
		} else {
			ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s%s", config.FrontEndOrigin, "/login?error=error_occured"))
			// ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": "something went wrong"})
			return
		}
	}

	if strings.Contains(config.TokenURL, "localhost") {
		ctx.SetCookie("token", token, config.TokenMaxAge*60, "/", "localhost", true, true)
	} else {
		ctx.SetCookie("token", token, config.TokenMaxAge*60, "/", config.TokenURL, true, true)
	}
	redirectUrl := ""
	// if pathUrl not begin with http or https
	fmt.Println("pathUrl", state)
	if !strings.HasPrefix(state, "http") && !strings.HasPrefix(state, "https") {
		state = "/profile"
		redirectUrl = fmt.Sprintf("%s%s", config.FrontEndOrigin, state)
	} else {
		redirectUrl = state
	}
	fmt.Println("redirectUrl", redirectUrl)

	ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}

func oauth2GitHub(ctx *gin.Context) {
	fmt.Println(ctx.Request.URL.String())

	code := ctx.Query("code")
	var pathUrl string = "/"

	if ctx.Query("state") != "" {
		pathUrl = ctx.Query("state")
	}

	if code == "" {
		fmt.Println(`CODE == ""`)
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Authorization code not provided!"})
		return
	}

	tokenRes, err := github.GetGitHubOauthToken(code)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message5": err.Error()})
		return
	}

	github_user, err := github.GetGitHubUser(tokenRes.Access_token)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message6": err.Error()})
		return
	}

	now := time.Now()
	email := strings.ToLower(github_user.Email)

	user_data := user_pkg.User{
		Name:      github_user.Name,
		Email:     email,
		Password:  "",
		Photo:     github_user.Photo,
		Provider:  "GitHub",
		Role:      "user",
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	fmt.Println("email", email)
	rows, err := db.RunSQL("SELECT `id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at` FROM `users` WHERE `email` = ? LIMIT 1", email)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message1": err.Error()})
		return
	}

	ifExist := rows.Next()
	rows.Close()
	fmt.Println("ifExist", ifExist)

	if !ifExist {
		// config.DB.Create(&user_data)
		fmt.Println("new user")
		rows, err := db.RunSQL("INSERT INTO `users`(`id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at`) VALUES (UUID(),?,?,?,?,?,?,?,?,?) RETURNING id ;", user_data.Name, user_data.Email, user_data.Password, user_data.Role, user_data.Photo, user_data.Verified, user_data.Provider, user_data.CreatedAt, user_data.UpdatedAt)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message2": err.Error()})
			return
		}
		rows.Close()
	}
	rows, err = db.RunSQL("SELECT `id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at` FROM `users` WHERE `email` = ? LIMIT 1", email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message3": err.Error()})
		return
	}
	defer rows.Close()

	var user user_pkg.User
	for rows.Next() {
		rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Photo, &user.Verified, &user.Provider, &user.CreatedAt, &user.UpdatedAt)
		break
	}

	if !ifExist {
		// get image from google
		// save image to public/images
		image.SaveImage(user.Photo, user.ID.String())

	}

	spew.Dump(user)

	if user.Provider != "GitHub" {
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

	fmt.Println("success redirect")

	fmt.Println("pathUrl", pathUrl)

	fmt.Println("config.FrontEndOrigin", config.FrontEndOrigin)

	ctx.Redirect(http.StatusPermanentRedirect, pathUrl)
	// ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(config.FrontEndOrigin, pathUrl))
}
