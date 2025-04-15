package entities

import (
	"github.com/winartodev/apollo/core/errors"
	"github.com/winartodev/apollo/modules/application"
	"regexp"
	"strings"
	"time"
)

var (
	slugInvalidCharsRegex = regexp.MustCompile(`[^a-z0-9-]+`)
	slugDashesRegex       = regexp.MustCompile(`-+`)
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
	ID          int64      `json:"id"`
	Slug        string     `json:"slug"`
	Name        string     `json:"name"`
	IsActive    bool       `json:"is_active"`
	Description string     `json:"description"`
	CreatedBy   int64      `json:"created_by"`
	UpdatedBy   int64      `json:"updated_by"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

func GenerateServiceSlug(data *Service) errors.Errors {
	if data == nil {
		return errors.MissingRequestBodyErr
	}

	dataCp := *data
	if dataCp.Name == "" {
		return application.ServiceNameIsEmptyErr
	}

	slug := strings.ToLower(strings.TrimSpace(dataCp.Name))
	slug = slugInvalidCharsRegex.ReplaceAllString(slug, "-")
	slug = slugDashesRegex.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "untitled"
	}

	dataCp.Slug = slug

	*data = dataCp

	return nil
}
