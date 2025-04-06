package repositories

const (
	GetUserRoleByUserIDQuery = `
		SELECT 
			role.id, 
			role.slug, 
			role.name
		FROM guardian_user_roles as gur 
		INNER JOIN guardian_roles as role 
			ON role.id = gur.role_id
		WHERE gur.user_id = $1
	`
)
