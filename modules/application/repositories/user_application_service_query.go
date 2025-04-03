package repositories

const (
	GetUserApplicationServices = `
		SELECT 
			uas.user_id as user_id, 
			app.id as app_id,
			app.slug as app_slug,
			app_service.id as app_service_id,
			app_service.scope as app_service_scope,
			app_service.slug as app_service_slug,
			app_service.name as app_service_name
		FROM user_application_services AS uas
		JOIN application_services as app_service 
			ON app_service.id = uas.application_service_id
		JOIN applications AS app 
			ON app.id = app_service.application_id
		WHERE uas.user_id = $1 
			AND app_service.application_id = $2
	`

	GetUserApplicationServiceBySlugQuery = `
		SELECT 
			uas.user_id as user_id, 
			app.id as app_id,
			app.slug as app_slug,
			app_service.id as app_service_id,
			app_service.scope as app_service_scope,
			app_service.slug as app_service_slug,
			app_service.name as app_service_name
		FROM user_application_services AS uas
		JOIN application_services as app_service 
			ON app_service.id = uas.application_service_id
		JOIN applications AS app 
			ON app.id = app_service.application_id
		WHERE uas.user_id = $1 
			AND app_service.application_id = $2
			AND app_service.slug = $3
	`
)
