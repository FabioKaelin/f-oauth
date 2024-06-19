package resetpassword

import (
	"math/rand"

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
	secret := randStringBytes(16)

	resetPasswordDB := db.DatabaseResetPassword{
		Secret: secret,
		UserID: userId,
	}

	id, err := db.CreateResetPassword(resetPasswordDB)
	if err != nil {
		return ResetPassword{}, err
	}

	resetPassword2 := ResetPassword{
		ID:     id,
		Secret: secret,

		// User: user.User{
		// 	ID: userId,
		// },
	}

	return resetPassword2, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// UseResetPassword resets the password of a user
func UseResetPassword(secret string, newPassword string) error {
	resetPassword2, err := db.GetResetPasswordBySecret(secret)
	if err != nil {
		return err
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

	err = db.DeleteResetPassword(resetPassword2.ID)
	if err != nil {
		return err
	}

	return nil
}
