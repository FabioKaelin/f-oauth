package initializers

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	FrontEndOrigin string `mapstructure:"FRONTEND_ORIGIN"`

	JWTTokenSecret string        `mapstructure:"JWT_SECRET"`
	TokenExpiresIn time.Duration `mapstructure:"TOKEN_EXPIRED_IN"`
	TokenMaxAge    int           `mapstructure:"TOKEN_MAXAGE"`
	TokenURL       string        `mapstructure:"TOKEN_URL"`

	GoogleClientID         string `mapstructure:"GOOGLE_OAUTH_CLIENT_ID"`
	GoogleClientSecret     string `mapstructure:"GOOGLE_OAUTH_CLIENT_SECRET"`
	GoogleOAuthRedirectUrl string `mapstructure:"GOOGLE_OAUTH_REDIRECT_URL"`

	GitHubClientID         string `mapstructure:"GITHUB_OAUTH_CLIENT_ID"`
	GitHubClientSecret     string `mapstructure:"GITHUB_OAUTH_CLIENT_SECRET"`
	GitHubOAuthRedirectUrl string `mapstructure:"GITHUB_OAUTH_REDIRECT_URL"`

	DatabaseHost     string `mapstructure:"DATABASE_HOST"`
	DatabaseUser     string `mapstructure:"DATABASE_USER"`
	DatabasePassword string `mapstructure:"DATABASE_PASSWORD"`
	DatabasePort     string `mapstructure:"DATABASE_PORT"`

	InternalTippURL    string `mapstructure:"INTERNAL_TIPP_URL"`
	InternalDevTippURL string `mapstructure:"INTERNAL_DEV_TIPP_URL"`

	InternalImageService string `mapstructure:"INTERNAL_IMAGE_SERVICE"`

	NotificationID string `mapstructure:"NOTIFICATION_ID"`
}

var StartConfig Config

func LoadConfig(path string) (config Config, err error) {
	// viper.AddConfigPath(path)
	// viper.SetConfigType("env")
	// viper.SetConfigName("app")

	// viper.AutomaticEnv()

	// err = viper.ReadInConfig()
	// if err != nil {
	// return
	// }

	// err = viper.Unmarshal(&config)
	godotenv.Load("app.env")
	TokenExpiresIn, _ := time.ParseDuration(os.Getenv("TOKEN_EXPIRED_IN"))
	marks, err := strconv.Atoi(os.Getenv("TOKEN_MAXAGE"))
	if err != nil {
		marks = 0
	}

	config = Config{
		FrontEndOrigin:         os.Getenv("FRONTEND_ORIGIN"),
		JWTTokenSecret:         os.Getenv("JWT_SECRET"),
		TokenExpiresIn:         TokenExpiresIn,
		TokenMaxAge:            marks,
		TokenURL:               os.Getenv("TOKEN_URL"),
		GoogleClientID:         os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		GoogleClientSecret:     os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		GoogleOAuthRedirectUrl: os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"),
		GitHubClientID:         os.Getenv("GITHUB_OAUTH_CLIENT_ID"),
		GitHubClientSecret:     os.Getenv("GITHUB_OAUTH_CLIENT_SECRET"),
		GitHubOAuthRedirectUrl: os.Getenv("GITHUB_OAUTH_REDIRECT_URL"),
		DatabaseHost:           os.Getenv("DATABASE_HOST"),
		DatabaseUser:           os.Getenv("DATABASE_USER"),
		DatabasePassword:       os.Getenv("DATABASE_PASSWORD"),
		DatabasePort:           os.Getenv("DATABASE_PORT"),
		InternalTippURL:        os.Getenv("INTERNAL_TIPP_URL"),
		InternalDevTippURL:     os.Getenv("INTERNAL_DEV_TIPP_URL"),
		InternalImageService:   os.Getenv("INTERNAL_IMAGE_SERVICE"),
		NotificationID:         os.Getenv("NOTIFICATION_ID"),
	}

	return
}
