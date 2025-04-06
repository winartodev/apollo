package entities

type GuardianUserAccessPermission struct {
	User        *GuardianUser        `json:"user"`
	Application *GuardianApplication `json:"application,omitempty"`
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
