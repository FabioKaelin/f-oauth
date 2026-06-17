package password

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/fabiokaelin/f-oauth/config"
	"github.com/fabiokaelin/f-oauth/pkg/auth"
	"github.com/fabiokaelin/f-oauth/pkg/db"
	"github.com/fabiokaelin/f-oauth/pkg/user"
)

type (
	ResetPassword struct {
		ID     string    `json:"id,omitempty"`
		Secret string    `json:"secret,omitempty"`
		User   user.User `json:"user,omitempty"`
	}
)

// CreateResetPassword creates a new reset password entry in the database
func CreateResetPassword(userId string) (ResetPassword, error) {
	secret, err := generateSecureToken(32)
	if err != nil {
		return ResetPassword{}, err
	}

	now := time.Now()
	expiresAt := now.Add(time.Duration(config.ResetPasswordTokenExpiry) * time.Second)

	resetPasswordDB := db.DatabaseResetPassword{
		Secret:    secret,
		UserID:    userId,
		ExpiresAt: expiresAt,
		Used:      false,
		CreatedAt: now,
	}

	id, err := db.CreateResetPassword(resetPasswordDB)
	if err != nil {
		return ResetPassword{}, err
	}

	return ResetPassword{
		ID:     id,
		Secret: secret,
	}, nil
}

// generateSecureToken generates a cryptographically secure random hex token of 2*length characters
func generateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// UseResetPassword resets the password of a user using a valid, unexpired, unused token
func UseResetPassword(secret string, newPassword string) error {
	resetPassword2, err := db.GetResetPasswordBySecret(secret)
	if err != nil {
		return err
	}

	if resetPassword2.ID == "" {
		return errors.New("invalid token")
	}

	if resetPassword2.Used {
		return errors.New("token already used")
	}

	if time.Now().After(resetPassword2.ExpiresAt) {
		return errors.New("token expired")
	}

	user, err := db.GetUserByID(resetPassword2.UserID)
	if err != nil {
		return err
	}

	hashedpassword, err := auth.HashAndSalt([]byte(newPassword))
	if err != nil {
		return err
	}

	user.Password = string(hashedpassword)

	err = db.UpdatePassword(user.ID, user.Password)
	if err != nil {
		return err
	}

	err = db.MarkResetPasswordUsed(resetPassword2.ID)
	if err != nil {
		return err
	}

	return nil
}
