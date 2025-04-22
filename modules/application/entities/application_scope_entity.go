package entities

import "time"

type ApplicationScope struct {
	ID            int64      `json:"id"`
	ApplicationID int64      `json:"application_id"`
	ScopeID       int64      `json:"scope_id"`
	IsActive      bool       `json:"is_active"`
	CreatedBy     int64      `json:"created_by"`
	UpdatedBy     int64      `json:"updated_by"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}
