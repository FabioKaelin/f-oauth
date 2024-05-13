package github

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

func LoginWithCode(code string) (string, int, error) {
	tokenRes, err := getGitHubOauthToken(code)

	if err != nil {
		fmt.Println("error, github failed to get token:", err)
		return "", 1, err
	}

	github_user, err := getGitHubUser(tokenRes.Access_token)

	if err != nil {
		fmt.Println("error, github failed to get user with token", err)
		return "", 1, err
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

	ifExist, err := db.DoesUserExist(email)

	if err != nil {
		fmt.Println("error, github, failed to get first time user", err)
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
			fmt.Println("error, github, failed to create user", err)
			return "", 1, err
		}

		notification.NewUserNotification(user_data)
	}

	user, err := db.GetUserByEmail(email)

	if err != nil {
		fmt.Println("error, github, failed to get user", err)
		return "", 1, err
	}

	if !ifExist {
		image.SaveImage(user.Photo, user.ID)
	}

	spew.Dump(user)

	if user.Provider != "GitHub" {
		return "", 2, errors.New("you have already signed up with a different method")
	}

	token, err := token_pkg.GenerateToken(config.TokenExpiresIn, user.ID, config.JWTTokenSecret)
	if err != nil {
		fmt.Println("error, github, failed to generate access token", err)
		return "", 1, err
	}

	return token, 0, nil
}
