package entities

import "time"

type GuardianUserRole struct {
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
}

type GuardianRole struct {
	ID          int64      `json:"id"`
	Slug        string     `json:"slug"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}
