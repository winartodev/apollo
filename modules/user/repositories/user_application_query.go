package repositories

const (
	GetUserApplicationAccessQuery = `
		SELECT 
			EXISTS (
				SELECT 1
				FROM user_applications AS usr_app
				JOIN application_scope AS app_scope 
					ON usr_app.application_scope_id = app_scope.id
				WHERE usr_app.user_id = $1
					AND app_scope.application_id = $2
					AND app_scope.scope_id = $3
			) as has_access, 
			EXISTS (
				SELECT 1
				FROM applications AS app
				WHERE app.is_active AND app.id = $4 
			) as is_app_active
	`

	InsertIntoUserApplication = `INSERT INTO user_applications (user_id, application_scope_id) VALUES %s`
)
