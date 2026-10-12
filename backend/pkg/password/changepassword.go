package password

import (
	"errors"
	"github.com/fabiokaelin/f-oauth/pkg/auth"

	"github.com/fabiokaelin/f-oauth/pkg/db"
)

func ChangePassword(userId string, newPassword string) error {
	if err := auth.ValidateNewPassword(newPassword); err != nil {
		return err
	}
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

// ChangePasswordVerified requires the existing local password and rejects a racing reset.
func ChangePasswordVerified(userID, current, replacement string) error {
	if err := auth.ValidateNewPassword(replacement); err != nil {
		return err
	}
	user, err := db.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user.ID == "" || user.Provider != "local" || !auth.ComparePasswords(user.Password, []byte(current)) {
		return errors.New("Current password is incorrect")
	}
	hash, err := auth.HashAndSalt([]byte(replacement))
	if err != nil {
		return err
	}
	return db.UpdatePasswordWithExpected(userID, user.Password, hash)
}
