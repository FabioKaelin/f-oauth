package db

type (
	DatabaseResetPassword struct {
		ID     string `json:"id,omitempty"`
		Secret string `json:"secret,omitempty"`
		UserID string `json:"useridfs,omitempty"`
	}
)

// CreateResetPassword creates a new reset password entry in the database
func CreateResetPassword(resetPassword DatabaseResetPassword) (string, error) {
	row, err := RunSQLRow("INSERT INTO `reset_password`(`id`, `secret`, `useridfs`) VALUES (UUID(),?,?) RETURNING id ;", resetPassword.Secret, resetPassword.UserID)
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
	rows, err := RunSQL("SELECT `id`, `secret`, `useridfs` FROM `reset_password` WHERE `secret` = ? LIMIT 1", secret)
	if err != nil {
		return DatabaseResetPassword{}, err
	}

	defer rows.Close()

	var resetPassword DatabaseResetPassword
	for rows.Next() {
		rows.Scan(&resetPassword.ID, &resetPassword.Secret, &resetPassword.UserID)
	}

	return resetPassword, nil
}

// DeleteResetPassword deletes a reset password entry by id
func DeleteResetPassword(id string) error {
	_, err := RunSQL("DELETE FROM `reset_password` WHERE `id` = ?", id)
	if err != nil {
		return err
	}

	return nil
}
