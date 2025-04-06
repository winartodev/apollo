package entities

import "time"

type User struct {
	ID              int64      `json:"id"`
	UUID            string     `json:"uuid"`
	Email           string     `json:"email"`
	PhoneNumber     string     `json:"phone_number"`
	Username        string     `json:"username"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	ProfilePicture  string     `json:"profile_picture"`
	Password        *string    `json:"password,omitempty"`
	RefreshToken    *string    `json:"refresh_token,omitempty"`
	IsEmailVerified bool       `json:"is_email_verified"`
	IsPhoneVerified bool       `json:"is_phone_verified"`
	LastLogin       *time.Time `json:"last_login,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
}

type UserUniqueField struct {
	Email       string
	PhoneNumber string
	Username    string
}

type UserUniqueFieldExists struct {
	IsEmailExists    bool
	IsPhoneExists    bool
	IsUsernameExists bool
}

type UserRole struct {
	RoleID int64  `json:"role_id"`
	Slug   string `json:"slug"`
	Name   string `json:"name"`
}
