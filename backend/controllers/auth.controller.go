package controllers

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fabiokaelin/f-oauth/initializers"
	"github.com/fabiokaelin/f-oauth/models"
	"github.com/fabiokaelin/f-oauth/utils"
	"github.com/gin-gonic/gin"
)

// SignUp User
func SignUpUser(ctx *gin.Context) {
	var payload *models.RegisterUserInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		fmt.Println("1", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "1message": err.Error()})
		return
	}

	hashpassword := sha512.Sum512([]byte(payload.Password))
	hashpasswordValue := hex.EncodeToString(hashpassword[:])

	now := time.Now()
	newUser := models.User{
		Name:      payload.Name,
		Email:     strings.ToLower(payload.Email),
		Password:  hashpasswordValue,
		Role:      "user",
		Provider:  "local",
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	rows, err := utils.RunSQL("SELECT `id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at` FROM `users` WHERE `email` = ? LIMIT 1", strings.ToLower(payload.Email))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "error appeared"})
		return
	}

	isExisting := rows.Next()

	if isExisting {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
		return
	}

	row, err := utils.RunSQLRow("INSERT INTO `users`(`id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at`) VALUES (UUID(), ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id ;", newUser.Name, newUser.Email, newUser.Password, newUser.Role, newUser.Photo, newUser.Verified, newUser.Provider, newUser.CreatedAt, newUser.UpdatedAt)

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
	isExisting := rows.Next()

	if !isExisting {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Photo, &user.Verified, &user.Provider, &user.CreatedAt, &user.UpdatedAt)

	hashpassword := sha512.Sum512([]byte(payload.Password))
	hashpasswordValue := hex.EncodeToString(hashpassword[:])

	if user.Password != hashpasswordValue {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	if user.Provider != "local" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You have already signed up with a different method"})
		return
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

	fmt.Println("email", email)
	rows, err := utils.RunSQL("SELECT `id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at` FROM `users` WHERE `email` = ? LIMIT 1", email)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message1": err.Error()})
		return
	}

	ifExist := rows.Next()
	fmt.Println("ifExist", ifExist)

	if !ifExist {
		// initializers.DB.Create(&user_data)
		fmt.Println("new user")
		_, err := utils.RunSQL("INSERT INTO `users`(`id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at`) VALUES (UUID(),?,?,?,?,?,?,?,?,?) RETURNING id ;", user_data.Name, user_data.Email, user_data.Password, user_data.Role, user_data.Photo, user_data.Verified, user_data.Provider, user_data.CreatedAt, user_data.UpdatedAt)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message2": err.Error()})
			return
		}
	}
	rows, err = utils.RunSQL("SELECT `id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at` FROM `users` WHERE `email` = ? LIMIT 1", email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message3": err.Error()})
		return
	}

	var user models.User
	for rows.Next() {
		rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Photo, &user.Verified, &user.Provider, &user.CreatedAt, &user.UpdatedAt)
		break
	}

	if !ifExist {
		// get image from google
		// save image to public/images
		saveImage(user.Photo, user.ID.String())

	}

	spew.Dump(user)

	if user.Provider != "Google" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You have already signed up with a different method"})
		return
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
		fmt.Println(`CODE == ""`)
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

	fmt.Println("email", email)
	rows, err := utils.RunSQL("SELECT `id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at` FROM `users` WHERE `email` = ? LIMIT 1", email)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message1": err.Error()})
		return
	}

	ifExist := rows.Next()
	fmt.Println("ifExist", ifExist)

	if !ifExist {
		// initializers.DB.Create(&user_data)
		fmt.Println("new user")
		_, err := utils.RunSQL("INSERT INTO `users`(`id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at`) VALUES (UUID(),?,?,?,?,?,?,?,?,?) RETURNING id ;", user_data.Name, user_data.Email, user_data.Password, user_data.Role, user_data.Photo, user_data.Verified, user_data.Provider, user_data.CreatedAt, user_data.UpdatedAt)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message2": err.Error()})
			return
		}
	}
	rows, err = utils.RunSQL("SELECT `id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at` FROM `users` WHERE `email` = ? LIMIT 1", email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message3": err.Error()})
		return
	}

	var user models.User
	for rows.Next() {
		rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Photo, &user.Verified, &user.Provider, &user.CreatedAt, &user.UpdatedAt)
		break
	}

	if !ifExist {
		// get image from google
		// save image to public/images
		saveImage(user.Photo, user.ID.String())

	}

	spew.Dump(user)

	if user.Provider != "GitHub" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You have already signed up with a different method"})
		return
	}
	config, _ := initializers.LoadConfig(".")

	token, err := utils.GenerateToken(config.TokenExpiresIn, user.ID, config.JWTTokenSecret)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("token", token, config.TokenMaxAge*60, "/", "localhost", false, true) // TODO: lh

	fmt.Println("success redirect")

	fmt.Println("pathUrl", pathUrl)

	fmt.Println("config.FrontEndOrigin", config.FrontEndOrigin)

	ctx.Redirect(http.StatusPermanentRedirect, config.FrontEndOrigin)
	// ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(config.FrontEndOrigin, pathUrl))
}

func saveImage(url string, userid string) error {
	// Create new file name
	newFileName := "profileimage-" + userid + ".png"
	fmt.Println("newFileName", newFileName)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode the image
	imageFile, _, err := image.Decode(resp.Body)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, imageFile); err != nil {
		return err
	}

	// push to INTERNAL_IMAGE_SERVICE as formFile with name "image"
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Create a new form file
	fw, err := w.CreateFormFile("image", newFileName)
	if err != nil {
		return err
	}

	// Write the image data to the form file
	if _, err = io.Copy(fw, buf); err != nil {
		return err
	}

	// Close the multipart writer
	if err = w.Close(); err != nil {
		return err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", initializers.StartConfig.InternalImageService+"/api/users/"+userid, &b)
	if err != nil {
		return err
	}

	// Set the content type, this is very important
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Do the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	return nil
}
