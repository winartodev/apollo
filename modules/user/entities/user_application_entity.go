package entities

import "time"

type UserApplication struct {
	UserID        int64      `json:"user_id"`
	ApplicationID int64      `json:"application_id"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

type UserApplicationResponse struct {
	ID       int64  `json:"id"`
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}
