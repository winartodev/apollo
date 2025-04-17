package repositories

const (
	GetUserApplicationAccessQuery = `
		SELECT EXISTS (
			SELECT 1
			FROM user_applications AS usr_app
			WHERE usr_app.user_id = $1
				AND usr_app.application_id = $2
				AND usr_app.scope_id = $3
		) AS has_access,
		EXISTS (
			SELECT 1
			FROM applications AS app
			WHERE app.is_active AND app.id = $4
		) AS is_app_active
	`
)
