package repository

const (
	CreateApplicationScopeQuery = "INSERT INTO application_scope (application_id, scope_id, is_active, created_by, updated_by, created_at, updated_at) VALUES %s RETURNING id"
)
