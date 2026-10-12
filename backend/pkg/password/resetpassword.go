package password

import (
	"crypto/rand"
	"encoding/hex"
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
	if err := auth.ValidateNewPassword(newPassword); err != nil {
		return err
	}
	hash, err := auth.HashAndSalt([]byte(newPassword))
	if err != nil {
		return err
	}
	return db.CompletePasswordReset(secret, hash)
}
