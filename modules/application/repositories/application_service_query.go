package repositories

const (
	GetApplicationServiceBySlug = `
		SELECT 
			app_service.id,
			app_service.application_id,
			app_service.scope,
			app_service.slug,
			app_service.name,
			app_service.is_active
		FROM application_services as app_service
		WHERE 
			app_service.slug = $1
	`
)
