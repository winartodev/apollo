package repositories

const (
	InsertIntoGuardianRoleQueryDB = `
		INSERT INTO guardian_roles
		    (
				application_id, 
				slug, 
				name, 
				description, 
				created_at, 
				updated_at
			) VALUES (
				  $1, 
				  $2, 
				  $3, 
				  $4, 
				  $5, 
				  $6
			) RETURNING id;
	`

	GetRolesQueryDB = `
		SELECT
		    id,
			application_id, 
			slug, 
			name, 
			description, 
			created_at, 
			updated_at
		FROM guardian_roles 
		WHERE application_id = $1  
	`

	GetRoleByIDQueryDB = `
		SELECT
		    id,
			application_id, 
			slug, 
			name, 
			description, 
			created_at, 
			updated_at
		FROM guardian_roles 
		WHERE application_id = $1
			AND id = $2
	`

	GetRoleBySlugQueryDB = `
		SELECT
			id,
			application_id, 
			slug, 
			name, 
			description, 
			created_at, 
			updated_at
		FROM guardian_roles 
		WHERE application_id = $1
			AND slug = $2
	`
)
