package db

type (
	DatabaseUser struct {
		ID       string `json:"id,omitempty"`
		Name     string `json:"name,omitempty"`
		Email    string `json:"email,omitempty"`
		Password string `json:"password,omitempty"`

		Role  string `json:"role,omitempty"`
		Photo string `json:"photo,omitempty"`

		Verified  bool   `json:"verified,omitempty"`
		Provider  string `json:"provider,omitempty"`
		CreatedAt string `json:"created_at,omitempty"`
		UpdatedAt string `json:"updated_at,omitempty"`
	}
)

func CreateUser(user DatabaseUser) (string, error) {
	row, err := RunSQLRow("INSERT INTO `users`(`id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at`) VALUES (UUID(),?,?,?,?,?,?,?,?,?) RETURNING id ;", user.Name, user.Email, user.Password, user.Role, user.Photo, user.Verified, user.Provider, user.CreatedAt, user.UpdatedAt)
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

func GetUserByEmail(email string) (DatabaseUser, error) {
	rows, err := RunSQL("SELECT `id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at` FROM `users` WHERE `email` = ? LIMIT 1", email)
	if err != nil {
		return DatabaseUser{}, err
	}

	defer rows.Close()

	var user DatabaseUser
	for rows.Next() {
		rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Photo, &user.Verified, &user.Provider, &user.CreatedAt, &user.UpdatedAt)
	}

	return user, nil
}

func DoesUserExist(email string) (bool, error) {
	rows, err := RunSQL("SELECT `id` FROM `users` WHERE `email` = ? LIMIT 1", email)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	ifExist := rows.Next()
	return ifExist, nil
}
