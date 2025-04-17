package entities

import (
	"github.com/winartodev/apollo/core/errors"
	"github.com/winartodev/apollo/core/helpers"
	"time"
)

const (
	applicationNameRequired = "application name is required"
)

var (
	AllowedApplicationFields = map[string]bool{
		"id":         true,
		"slug":       true,
		"name":       true,
		"is_active":  true,
		"created_by": true,
		"created_at": true,
	}
)

type Application struct {
	ID        int64      `json:"id"`
	Slug      string     `json:"slug"`
	Name      string     `json:"name"`
	IsActive  bool       `json:"is_active"`
	CreatedBy int64      `json:"created_by"`
	UpdatedBy int64      `json:"updated_by"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func GenerateApplicationSlug(data *Application) errors.Errors {
	if data == nil {
		return errors.MissingRequestBodyErr
	}

	dataCp := *data
	if dataCp.Name == "" {
		return errors.BadRequestErr.WithReason(applicationNameRequired)
	}

	dataCp.Slug = helpers.MakeSlug(dataCp.Name)

	*data = dataCp

	return nil
}
