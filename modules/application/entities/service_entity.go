package entities

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var (
	slugInvalidCharsRegex = regexp.MustCompile(`[^a-z0-9-]+`)
	slugDashesRegex       = regexp.MustCompile(`-+`)
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

func GenerateServiceSlug(data *Service) error {
	if data == nil {
		return errors.New("data is nil")
	}

	dataCp := *data
	if dataCp.Name == "" {
		return errors.New("data is required")
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
