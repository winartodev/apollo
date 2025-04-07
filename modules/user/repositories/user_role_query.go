package repositories

const (
	GetUserRoleByIDQuery = `
		SELECT
			gr.id AS role_id,
			gr.name AS role_name,
			gr.slug AS role_slug
		FROM guardian_user_roles as gur 
		JOIN guardian_roles as gr
			ON gr.id = gur.role_id
		WHERE gur.user_id = $1
			AND gur.application_id = $2
	`
)
