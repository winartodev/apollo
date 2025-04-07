package entities

import (
	userEntity "github.com/winartodev/apollo/modules/user/entities"
)

type GuardianUserAccessPermission struct {
	User        *GuardianUser        `json:"user"`
	Application *GuardianApplication `json:"application,omitempty"`
}

func (g *GuardianUserAccessPermission) Build(user *userEntity.User, userRole *userEntity.UserRoleResponse,
	userApplication *userEntity.UserApplicationResponse, applicationServices *GuardianApplicationService) GuardianUserAccessPermission {

	services := make([]GuardianApplicationService, 0)
	services = append(services, *applicationServices)

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
			ID:       userApplication.ID,
			Slug:     userApplication.Slug,
			Name:     userApplication.Name,
			IsActive: userApplication.IsActive,
			Services: services,
		},
	}
}

type GuardianUser struct {
	ID           int64         `json:"id"`
	Email        string        `json:"email"`
	PhoneNumber  string        `json:"phone_number"`
	GuardianRole *GuardianRole `json:"role,omitempty"`
}

type GuardianRole struct {
	ID   int64  `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type GuardianApplication struct {
	ID       int64                        `json:"id"`
	Slug     string                       `json:"slug"`
	Name     string                       `json:"name"`
	IsActive bool                         `json:"is_active"`
	Services []GuardianApplicationService `json:"service,omitempty"`
}

type GuardianApplicationService struct {
	ID          int64                `json:"id"`
	Scope       string               `json:"scope"`
	Slug        string               `json:"slug"`
	Name        string               `json:"name"`
	Permissions []GuardianPermission `json:"permissions,omitempty"`
}

type GuardianPermission struct {
	ID   int64  `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name,omitempty"`
}
