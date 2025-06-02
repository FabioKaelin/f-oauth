package password

import (
	"github.com/fabiokaelin/f-oauth/pkg/auth"

	"github.com/fabiokaelin/f-oauth/pkg/db"
)

func ChangePassword(userId string, newPassword string) error {
	hashedpassword, err := auth.HashAndSalt([]byte(newPassword))
	if err != nil {
		return err
	}

	hashedPassword := string(hashedpassword)

	// Update the user's password in the database
	err = db.UpdatePassword(userId, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}
