package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/google-github-oath2-golang/initializers"
	"github.com/wpcodevo/google-github-oath2-golang/models"
	"github.com/wpcodevo/google-github-oath2-golang/utils"
)

// SignUp User
func SignUpUser(ctx *gin.Context) {
	var payload *models.RegisterUserInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		fmt.Println("1", err)
		spew.Dump(payload)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "1message": err.Error()})
		return
	}

	now := time.Now()
	newUser := models.User{
		Name:      payload.Name,
		Email:     strings.ToLower(payload.Email),
		Password:  payload.Password,
		Role:      "user",
		Provider:  "local",
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	rows, err := utils.RunSQL("INSERT INTO `users`(`id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at`) VALUES (UUID(),?,?,?,?,?,?,?,?,?) RETURNING id ;", newUser.Name, newUser.Email, newUser.Password, newUser.Role, newUser.Photo, newUser.Verified, newUser.Provider, newUser.CreatedAt, newUser.UpdatedAt)

	for rows.Next() {
		rows.Scan(&newUser.ID)
	}

	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {

		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "2message": "User with that email already exists"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "3message": "Something bad happened"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": models.FilteredResponse(&newUser)}})
}

// SignIn User
func SignInUser(ctx *gin.Context) {
	var payload *models.LoginUserInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		spew.Dump(payload)
		fmt.Println("1", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User

	rows, err := utils.RunSQL("SELECT `id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at` FROM `users` WHERE `email` = ? LIMIT 1", strings.ToLower(payload.Email))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	for rows.Next() {
		rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Photo, &user.Verified, &user.Provider, &user.CreatedAt, &user.UpdatedAt)
		break
	}

	config, _ := initializers.LoadConfig(".")

	token, err := utils.GenerateToken(config.TokenExpiresIn, user.ID, config.JWTTokenSecret)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("token", token, config.TokenMaxAge*60, "/", initializers.StartConfig.TokenURL, false, true) // TODO: lh
	ctx.SetCookie("token", token, config.TokenMaxAge*60, "/", "localhost", false, true)                       // TODO: lh
	// ctx.SetCookie("token", token, config.TokenMaxAge*60, "/", "localhost", false, true) // TODO: lh

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", initializers.StartConfig.TokenURL, false, true)
	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
	// ctx.SetCookie("token", "", -1, "/", "localhost", false, true) // TODO: lh
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func GoogleOAuth(ctx *gin.Context) {
	code := ctx.Query("code")
	var pathUrl string = "/"

	if ctx.Query("state") != "" {
		pathUrl = ctx.Query("state")
	}

	if code == "" {
		fmt.Println(`CODE == ""`)
		ctx.JSON(http.StatusUnauthorized, `"status": "fail", "message": "Authorization code not provided!"`)
		// ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Authorization code not provided!"})
		return
	}

	tokenRes, err := utils.GetGoogleOauthToken(code)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	google_user, err := utils.GetGoogleUser(tokenRes.Access_token, tokenRes.Id_token)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	now := time.Now()
	email := strings.ToLower(google_user.Email)

	user_data := models.User{
		Name:      google_user.Name,
		Email:     email,
		Password:  "",
		Photo:     google_user.Picture,
		Provider:  "Google",
		Role:      "user",
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	rows, err := utils.RunSQL("UPDATE `users` SET `name`= ? ,`password`= ? ,`role`= ? ,`photo`= ? ,`verified`= ? ,`provider`= ? ,`created_at`= ? ,`updated_at`= ?  WHERE `email` = ?;", user_data.Name, user_data.Password, user_data.Role, user_data.Photo, user_data.Verified, user_data.Provider, user_data.CreatedAt, user_data.UpdatedAt, user_data.Email)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message1": err.Error()})
		return
	}

	ifExist := rows.Next()

	if !ifExist {
		// initializers.DB.Create(&user_data)
		_, err := utils.RunSQL("INSERT INTO `users`(`id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at`) VALUES (UUID(),?,?,?,?,?,?,?,?,?) RETURNING id ;", user_data.Name, user_data.Email, user_data.Password, user_data.Role, user_data.Photo, user_data.Verified, user_data.Provider, user_data.CreatedAt, user_data.UpdatedAt)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message2": err.Error()})
			return
		}
	}

	var user models.User
	rows, err = utils.RunSQL("SELECT `id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at` FROM `users` WHERE `email` = ? LIMIT 1", email)
	// initializers.DB.First(&user, "email = ?", email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message3": err.Error()})
		return
	}
	for rows.Next() {
		rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Photo, &user.Verified, &user.Provider, &user.CreatedAt, &user.UpdatedAt)
		break
	}

	config, _ := initializers.LoadConfig(".")

	token, err := utils.GenerateToken(config.TokenExpiresIn, user.ID.String(), config.JWTTokenSecret)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message4": err.Error()})
		return
	}

	// ctx.SetCookie("token", token, config.TokenMaxAge*60, "/", "localhost", false, true) // TODO: lh
	ctx.SetCookie("token", token, config.TokenMaxAge*60, "/", initializers.StartConfig.TokenURL, false, true) // TODO: lh
	ctx.SetCookie("token", token, config.TokenMaxAge*60, "/", "localhost", false, true)                       // TODO: lh
	redirectUrl := ""
	// if pathUrl not begin with http or https
	fmt.Println("pathUrl", pathUrl)
	if !strings.HasPrefix(pathUrl, "http") && !strings.HasPrefix(pathUrl, "https") {
		redirectUrl = fmt.Sprint(config.FrontEndOrigin, pathUrl)
	} else {
		redirectUrl = pathUrl
	}
	fmt.Println("redirectUrl", redirectUrl)

	ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}

func GitHubOAuth(ctx *gin.Context) {
	code := ctx.Query("code")
	var pathUrl string = "/"

	if ctx.Query("state") != "" {
		pathUrl = ctx.Query("state")
	}

	if code == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Authorization code not provided!"})
		return
	}

	tokenRes, err := utils.GetGitHubOauthToken(code)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message5": err.Error()})
		return
	}

	github_user, err := utils.GetGitHubUser(tokenRes.Access_token)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message6": err.Error()})
		return
	}

	now := time.Now()
	email := strings.ToLower(github_user.Email)

	user_data := models.User{
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

	rows, err := utils.RunSQL("UPDATE `users` SET `name`= ? ,`password`= ? ,`role`= ? ,`photo`= ? ,`verified`= ? ,`provider`= ? ,`created_at`= ? ,`updated_at`= ?  WHERE `email` = ?;", user_data.Name, user_data.Password, user_data.Role, user_data.Photo, user_data.Verified, user_data.Provider, user_data.CreatedAt, user_data.UpdatedAt, user_data.Email)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if !rows.Next() {
		// initializers.DB.Create(&user_data)
		_, err := utils.RunSQL("INSERT INTO `users`(`id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at`) VALUES (UUID(),?,?,?,?,?,?,?,?,?) RETURNING id ;", user_data.Name, user_data.Email, user_data.Password, user_data.Role, user_data.Photo, user_data.Verified, user_data.Provider, user_data.CreatedAt, user_data.UpdatedAt)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
			return
		}
	}

	var user models.User
	rows, err = utils.RunSQL("SELECT `id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at` FROM `users` WHERE `email` = ? LIMIT 1", email)
	// initializers.DB.First(&user, "email = ?", email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	for rows.Next() {
		rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Photo, &user.Verified, &user.Provider, &user.CreatedAt, &user.UpdatedAt)
		break
	}
	config, _ := initializers.LoadConfig(".")

	token, err := utils.GenerateToken(config.TokenExpiresIn, user.ID, config.JWTTokenSecret)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("token", token, config.TokenMaxAge*60, "/", "localhost", false, true) // TODO: lh

	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(config.FrontEndOrigin, pathUrl))
}
