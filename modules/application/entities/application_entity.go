package entities

import (
	userEntity "github.com/winartodev/apollo/modules/user/entities"
	"time"
)

type ApplicationAccess struct {
	User         userEntity.User `json:"user"`
	Applications []Application   `json:"applications"`
}

type Application struct {
	ID          int64      `json:"id"`
	Slug        string     `json:"slug"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type ApplicationService struct {
	ID            int64      `json:"id"`
	ApplicationID int64      `json:"application_id"`
	Scope         string     `json:"scope"`
	Slug          string     `json:"slug"`
	Name          string     `json:"name"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}
