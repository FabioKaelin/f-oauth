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
	rows, err := db.RunSQL("SELECT `id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at` FROM `users` WHERE `email` = ? LIMIT 1", email)

	if err != nil {
		fmt.Println("error, github, failed to get first time user", err)
		return "", 1, err
	}

	ifExist := rows.Next()
	rows.Close()
	fmt.Println("ifExist", ifExist)

	if !ifExist {
		// config.DB.Create(&user_data)
		fmt.Println("new user")
		rows, err := db.RunSQL("INSERT INTO `users`(`id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at`) VALUES (UUID(),?,?,?,?,?,?,?,?,?) RETURNING id ;", user_data.Name, user_data.Email, user_data.Password, user_data.Role, user_data.Photo, user_data.Verified, user_data.Provider, user_data.CreatedAt, user_data.UpdatedAt)

		if err != nil {
			fmt.Println("error, github, failed to insert new user", err)
			return "", 1, err
		}
		rows.Close()
	}
	rows, err = db.RunSQL("SELECT `id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at` FROM `users` WHERE `email` = ? LIMIT 1", email)
	if err != nil {
		fmt.Println("error, github, failed to get user the second time", err)
		return "", 1, err
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
		return "", 2, errors.New("you have already signed up with a different method")
	}

	token, err := token_pkg.GenerateToken(config.TokenExpiresIn, user.ID.String(), config.JWTTokenSecret)
	if err != nil {
		fmt.Println("error, github, failed to generate access token", err)
		return "", 1, err
	}

	return token, 0, nil
}
