package repositories

const (
	GetUserApplicationsByUserIDQueryDB = `
		SELECT DISTINCT
			a.id AS app_id,
			a.slug AS app_slug,
			a.name AS app_name,
			a.is_active AS app_is_active
		FROM user_applications AS ua
		INNER JOIN users AS u 
			ON u.id = ua.user_id
		INNER JOIN applications AS a 
			ON a.id = ua.application_id
			AND a.is_active = true
		WHERE u.id = $1
	`

	GetUserApplicationByUserIDAndApplicationSlugQuery = `
		SELECT DISTINCT
			a.id AS app_id,
			a.slug AS app_slug,
			a.name AS app_name,
			a.is_active AS app_is_active
		FROM user_applications AS ua
		INNER JOIN users AS u 
			ON u.id = ua.user_id
		INNER JOIN applications AS a 
			ON a.id = ua.application_id
			AND a.is_active = true
		WHERE u.id = $1 AND a.slug = $2
	`
)
