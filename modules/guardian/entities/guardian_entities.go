package entities

import (
	applicationEntity "github.com/winartodev/apollo/modules/application/entities"
	userEntity "github.com/winartodev/apollo/modules/user/entities"
)

type GuardianUserAccessPermission struct {
	User        *GuardianUser        `json:"user"`
	Application *GuardianApplication `json:"application,omitempty"`
}

func (g *GuardianUserAccessPermission) Build(user *userEntity.User, userRole *userEntity.UserRole, application *applicationEntity.Application, appService *applicationEntity.ApplicationService) GuardianUserAccessPermission {
	return GuardianUserAccessPermission{
		User: &GuardianUser{
			ID:          user.ID,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			GuardianRole: &GuardianRole{
				ID:   userRole.RoleID,
				Slug: userRole.Slug,
				Name: userRole.Name,
			},
		},
		Application: &GuardianApplication{
			ID:       application.ID,
			Slug:     application.Slug,
			Name:     application.Name,
			IsActive: application.IsActive,
			Service: &GuardianApplicationService{
				ID:    appService.ID,
				Scope: appService.Scope,
				Slug:  appService.Slug,
				Name:  appService.Name,
			},
		},
	}
}

type GuardianUser struct {
	ID           int64         `json:"id"`
	Email        string        `json:"email"`
	PhoneNumber  string        `json:"phone_number"`
	GuardianRole *GuardianRole `json:"role,omitempty"`
}

type GuardianApplication struct {
	ID       int64                       `json:"id"`
	Slug     string                      `json:"slug"`
	Name     string                      `json:"name"`
	IsActive bool                        `json:"is_active"`
	Service  *GuardianApplicationService `json:"service,omitempty"`
}

type GuardianApplicationService struct {
	ID          int64    `json:"id"`
	Scope       string   `json:"scope"`
	Slug        string   `json:"slug"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}
