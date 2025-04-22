package entities

type UserApplication struct {
	UserID             int64 `json:"user_id"`
	ApplicationScopeID int64 `json:"application_scope_id"`
}
