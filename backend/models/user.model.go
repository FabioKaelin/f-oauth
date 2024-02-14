package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Password string

	Role  string
	Photo string

	Verified  bool
	Provider  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RegisterUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" bindinig:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserInput struct {
	Email    string `json:"email" bindinig:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Role     string `json:"role,omitempty"`
	Provider string `json:"provider,omitempty"`
	Photo    string `json:"photo,omitempty"`
	Verified bool   `json:"verified,omitempty"`
	// CreatedAt time.Time
	// UpdatedAt time.Time
}

func FilteredResponse(user *User) UserResponse {
	return UserResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		Name:     user.Name,
		Role:     user.Role,
		Verified: user.Verified,
		Photo:    user.Photo,
		Provider: user.Provider,
		// CreatedAt: user.CreatedAt,
		// UpdatedAt: user.UpdatedAt,
	}
}
