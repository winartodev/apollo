package entities

import (
	"fmt"
	"github.com/winartodev/apollo/core/errors"
	"github.com/winartodev/apollo/core/helpers"
	"time"
)

const (
	serviceNameRequired         = "service name is required"
	serviceMissingApplicationID = "application id is required"
)

var (
	AllowedServiceFields = map[string]bool{
		"id":         true,
		"slug":       true,
		"name":       true,
		"is_active":  true,
		"created_by": true,
		"created_at": true,
	}
)

type Service struct {
	ID            int64      `json:"id"`
	ApplicationID *int64     `json:"application_id"`
	Slug          string     `json:"slug"`
	Name          string     `json:"name"`
	IsActive      bool       `json:"is_active"`
	Description   string     `json:"description"`
	CreatedBy     int64      `json:"created_by"`
	UpdatedBy     int64      `json:"updated_by"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

func GenerateServiceSlug(data *Service) errors.Errors {
	if data == nil {
		return errors.MissingRequestBodyErr
	}

	dataCp := *data
	if dataCp.Name == "" {
		return errors.BadRequestErr.WithReason(serviceNameRequired)
	}

	if dataCp.ApplicationID == nil {
		return errors.BadRequestErr.WithReason(serviceMissingApplicationID)
	}

	dataCp.Slug = helpers.MakeSlug(fmt.Sprintf("%d %s", *dataCp.ApplicationID, dataCp.Name))

	*data = dataCp

	return nil
}
