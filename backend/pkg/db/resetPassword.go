package db

import "time"

type (
	DatabaseResetPassword struct {
		ID        string    `json:"id,omitempty"`
		Secret    string    `json:"secret,omitempty"`
		UserID    string    `json:"useridfs,omitempty"`
		ExpiresAt time.Time `json:"expires_at,omitempty"`
		Used      bool      `json:"used,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
	}
)

// CreateResetPassword creates a new reset password entry in the database
func CreateResetPassword(resetPassword DatabaseResetPassword) (string, error) {
	row, err := RunSQLRow(
		"INSERT INTO `reset_password`(`id`, `secret`, `useridfs`, `expires_at`, `used`, `created_at`) VALUES (UUID(),?,?,?,?,?) RETURNING id ;",
		resetPassword.Secret, resetPassword.UserID, resetPassword.ExpiresAt, resetPassword.Used, resetPassword.CreatedAt,
	)
	if err != nil {
		return "", err
	}

	if row.Err() != nil {
		return "", row.Err()
	}

	var id string
	err = row.Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

// GetResetPasswordBySecret gets a reset password entry by secret
func GetResetPasswordBySecret(secret string) (DatabaseResetPassword, error) {
	rows, err := RunSQL("SELECT `id`, `secret`, `useridfs`, `expires_at`, `used`, `created_at` FROM `reset_password` WHERE `secret` = ? LIMIT 1", secret)
	if err != nil {
		return DatabaseResetPassword{}, err
	}

	defer rows.Close()

	var resetPassword DatabaseResetPassword
	for rows.Next() {
		rows.Scan(&resetPassword.ID, &resetPassword.Secret, &resetPassword.UserID, &resetPassword.ExpiresAt, &resetPassword.Used, &resetPassword.CreatedAt)
	}

	return resetPassword, nil
}

// MarkResetPasswordUsed marks a reset password entry as used
func MarkResetPasswordUsed(id string) error {
	_, err := RunSQL("UPDATE `reset_password` SET `used` = 1 WHERE `id` = ?", id)
	if err != nil {
		return err
	}
	return nil
}

// CountRecentResetPasswordsByUserID counts reset password entries created after 'since' for a given user
func CountRecentResetPasswordsByUserID(userID string, since time.Time) (int, error) {
	row, err := RunSQLRow("SELECT COUNT(*) FROM `reset_password` WHERE `useridfs` = ? AND `created_at` > ?", userID, since)
	if err != nil {
		return 0, err
	}
	var count int
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// DeleteResetPassword deletes a reset password entry by id
func DeleteResetPassword(id string) error {
	_, err := RunSQL("DELETE FROM `reset_password` WHERE `id` = ?", id)
	if err != nil {
		return err
	}

	return nil
}
