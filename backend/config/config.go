package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	GinMode string

	FrontEndOrigin string

	JWTTokenSecret string
	TokenExpiresIn time.Duration
	TokenMaxAge    int
	TokenURL       string

	GoogleClientID         string
	GoogleClientSecret     string
	GoogleOAuthRedirectUrl string

	GitHubClientID         string
	GitHubClientSecret     string
	GitHubOAuthRedirectUrl string

	DatabaseHost     string
	DatabaseUser     string
	DatabasePassword string
	DatabasePort     string

	InternalTippURL      string
	InternalDevTippURL   string
	InternalImageService string

	NotificationID string
)

func getString(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("key '%s' not found", key)
	}
	return value, nil
}

func getInt(key string) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return 0, fmt.Errorf("key '%s' not found", key)
	}
	return strconv.Atoi(value)
}

// func getBool(key string) (bool, error) {
// 	value := os.Getenv(key)
// 	if value == "" {
// 		return false, fmt.Errorf("key '%s' not found", key)
// 	}
// 	return strconv.ParseBool(value)
// }

func getDuration(key string) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		return 0, fmt.Errorf("key '%s' not found", key)
	}
	return time.ParseDuration(value)
}

func Load() error {
	godotenv.Load("app.env")

	ginmode, err := getString("GIN_MODE")
	if err != nil {
		return err
	}
	if ginmode != "debug" && ginmode != "release" {
		ginmode = "debug"
	}
	GinMode = ginmode

	FrontEndOrigin, err = getString("FRONTEND_ORIGIN")
	if err != nil {
		return err
	}

	JWTTokenSecret, err = getString("JWT_SECRET")
	if err != nil {
		return err
	}

	TokenExpiresIn, err = getDuration("TOKEN_EXPIRED_IN")
	if err != nil {
		return err
	}

	TokenMaxAge, err = getInt("TOKEN_MAXAGE")
	if err != nil {
		return err
	}

	TokenURL, err = getString("TOKEN_URL")
	if err != nil {
		return err
	}

	GoogleClientID, err = getString("GOOGLE_OAUTH_CLIENT_ID")
	if err != nil {
		return err
	}

	GoogleClientSecret, err = getString("GOOGLE_OAUTH_CLIENT_SECRET")
	if err != nil {
		return err
	}

	GoogleOAuthRedirectUrl, err = getString("GOOGLE_OAUTH_REDIRECT_URL")
	if err != nil {
		return err
	}

	GitHubClientID, err = getString("GITHUB_OAUTH_CLIENT_ID")
	if err != nil {
		return err
	}

	GitHubClientSecret, err = getString("GITHUB_OAUTH_CLIENT_SECRET")
	if err != nil {
		return err
	}

	GitHubOAuthRedirectUrl, err = getString("GITHUB_OAUTH_REDIRECT_URL")
	if err != nil {
		return err
	}

	DatabaseHost, err = getString("DATABASE_HOST")
	if err != nil {
		return err
	}

	DatabaseUser, err = getString("DATABASE_USER")
	if err != nil {
		return err
	}

	DatabasePassword, err = getString("DATABASE_PASSWORD")
	if err != nil {
		return err
	}

	DatabasePort, err = getString("DATABASE_PORT")
	if err != nil {
		return err
	}

	InternalTippURL, err = getString("INTERNAL_TIPP_URL")
	if err != nil {
		return err
	}

	InternalDevTippURL, err = getString("INTERNAL_DEV_TIPP_URL")
	if err != nil {
		return err
	}

	InternalImageService, err = getString("INTERNAL_IMAGE_SERVICE")
	if err != nil {
		return err
	}

	NotificationID, err = getString("NOTIFICATION_ID")
	if err != nil {
		return err
	}

	return nil
}
