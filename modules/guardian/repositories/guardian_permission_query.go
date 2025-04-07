package repositories

const (
	GetServicePermissionByUserID = `
		SELECT
			grp.application_service_id AS service_id,
			ARRAY_AGG(gp.id) AS permission_ids,
			ARRAY_AGG(gp.slug) AS permission_slugs
		FROM guardian_user_roles as gur 
		JOIN guardian_role_permissions as grp 
			ON grp.role_id = gur.role_id
		JOIN guardian_permissions AS gp 
			ON gp.id = grp.permission_id
		WHERE gur.user_id = $1 
			AND gur.role_id = $2
			AND gur.application_id = $3
			AND grp.application_service_id = $4 
		GROUP BY 
			gur.user_id, 
			gur.role_id, 
			gur.application_id, 
			grp.application_service_id
		ORDER BY grp.application_service_id DESC
	`
)
