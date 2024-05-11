package google

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fabiokaelin/f-oauth/config"
	"github.com/fabiokaelin/f-oauth/pkg/db"
	"github.com/fabiokaelin/f-oauth/pkg/image"
	"github.com/fabiokaelin/f-oauth/pkg/notification"
	token_pkg "github.com/fabiokaelin/f-oauth/pkg/token"
	user_pkg "github.com/fabiokaelin/f-oauth/pkg/user"
)

// LoginWithCode return a token the errorcode and an error
// 0 no error
// 1 normal error
// 2 user already signed up with a different method
func LoginWithCode(code string) (string, int, error) {
	tokenRes, err := getGoogleOauthToken(code)

	if err != nil {
		fmt.Println("error, google failed to get token:", err)
		return "", 1, err
	}

	google_user, err := getGoogleUser(tokenRes.Access_token, tokenRes.Id_token)

	if err != nil {
		fmt.Println("error, google failed to get user with token", err)
		return "", 1, err
	}

	now := time.Now()
	email := strings.ToLower(google_user.Email)

	user_data := user_pkg.User{
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

	ifExist, err := db.DoesUserExist(email)

	if err != nil {
		fmt.Println("error, google, failed to get first time user", err)
		return "", 1, err
	}

	fmt.Println("ifExist", ifExist)

	if !ifExist {
		fmt.Println("new user")
		dbuser := db.DatabaseUser{
			Name:      user_data.Name,
			Email:     user_data.Email,
			Password:  user_data.Password,
			Role:      user_data.Role,
			Photo:     user_data.Photo,
			Verified:  user_data.Verified,
			Provider:  user_data.Provider,
			CreatedAt: user_data.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user_data.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		_, err := db.CreateUser(dbuser)
		if err != nil {
			fmt.Println("error, google, failed to create user", err)
			return "", 1, err
		}

		notification.NewUserNotification(user_data)

	}
	user, err := db.GetUserByEmail(email)

	if err != nil {
		fmt.Println("error, google, failed to get user", err)
		return "", 1, err
	}

	if !ifExist {
		image.SaveImage(user.Photo, user.ID)
	}

	spew.Dump(user)

	if user.Provider != "Google" {
		return "", 2, errors.New("you have already signed up with a different method")
	}

	token, err := token_pkg.GenerateToken(config.TokenExpiresIn, user.ID, config.JWTTokenSecret)
	if err != nil {
		fmt.Println("error, google, failed to generate access token", err)
		return "", 1, err
	}

	return token, 0, nil
}
