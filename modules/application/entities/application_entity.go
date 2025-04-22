package entities

import (
	"github.com/winartodev/apollo/core/errors"
	"github.com/winartodev/apollo/core/helpers"
	"github.com/winartodev/apollo/modules/application"
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
	ID          int64               `json:"id"`
	Slug        string              `json:"slug"`
	Name        string              `json:"name"`
	IsActive    bool                `json:"is_active"`
	Description string              `json:"description"`
	CreatedBy   int64               `json:"created_by"`
	UpdatedBy   int64               `json:"updated_by"`
	CreatedAt   *time.Time          `json:"created_at"`
	UpdatedAt   *time.Time          `json:"updated_at"`
	Scopes      []application.Scope `json:"scopes,omitempty"`
	Services    []Service           `json:"services,omitempty"`
}

type ApplicationResponse struct {
	ID          int64      `json:"id"`
	Slug        string     `json:"slug"`
	Name        string     `json:"name"`
	IsActive    bool       `json:"is_active"`
	Description string     `json:"description"`
	CreatedBy   int64      `json:"created_by"`
	UpdatedBy   int64      `json:"updated_by"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Scopes      []string   `json:"scopes"`
	Services    []Service  `json:"services"`
}

func (ar *Application) ToResponse() *ApplicationResponse {
	res := &ApplicationResponse{
		ID:          ar.ID,
		Slug:        ar.Slug,
		Name:        ar.Name,
		IsActive:    ar.IsActive,
		Description: ar.Description,
		CreatedBy:   ar.CreatedBy,
		UpdatedBy:   ar.UpdatedBy,
		CreatedAt:   ar.CreatedAt,
		UpdatedAt:   ar.UpdatedAt,
		Scopes:      nil,
		Services:    nil,
	}

	if ar.Scopes != nil {
		var scopes []string
		for _, scope := range ar.Scopes {
			scopes = append(scopes, scope.ToString())
		}

		res.Scopes = scopes
	}

	if ar.Services != nil {
		var services []Service
		for _, service := range ar.Services {
			services = append(services, service)
		}

		res.Services = services
	}

	return res
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
