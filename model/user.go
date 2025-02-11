package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents users table in the database
type User struct {
	User_id   uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Is_delete bool      `json:"is_delete"`
}

type UserDetailModel struct {
	User_id  uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
}

type UserResponse struct {
	User_id  uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Phone    string    `json:"phone"`
	Token    string    `json:"token"`
}

// Login data model
type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register data model
type RegisterUser struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}
