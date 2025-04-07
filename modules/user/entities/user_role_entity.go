package entities

type UserRole struct {
	UserID    int64 `json:"user_id"`
	RoleID    int64 `json:"role_id"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

type UserRoleResponse struct {
	RoleID int64  `json:"role_id"`
	Slug   string `json:"slug"`
	Name   string `json:"name"`
}
