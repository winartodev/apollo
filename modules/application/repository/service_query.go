package repository

const (
	InsertIntoServiceQuery = `
		INSERT INTO services (slug, name, description, is_active, created_by, updated_by, created_at, updated_at, application_id) 
		VALUES %s RETURNING id;
	`

	GetServiceBySlugQuery = `
		SELECT 
		    id, 
		    application_id,
		    slug, 
		    name, 
		    description, 
		    is_active, 
		    created_by, 
		    updated_by, 
		    created_at, 
		    updated_at 
		FROM services 
		WHERE slug = $1 AND application_id =$2;
	`

	GetServiceByIDQuery = `
		SELECT 
		    id, 
			application_id,
		    slug, 
		    name, 
		    description, 
		    is_active, 
		    created_by, 
		    updated_by, 
		    created_at, 
		    updated_at 
		FROM services 
		WHERE id = $1 AND application_id = $2;
	`

	GetServicesQuery = `
		SELECT 
		    id, 
			application_id,
		    slug, 
		    name, 
		    description, 
		    is_active, 
		    created_by, 
		    updated_by, 
		    created_at, 
		    updated_at 
		FROM services
	`

	CountServiceQuery = `SELECT COUNT(*) FROM services WHERE application_id = $1;`

	UpdateServiceQuery = `
		UPDATE services 
		SET slug = $1, 
			name = $2,
			description = $3,
			is_active = $4,
			updated_by = $5,
			updated_at = $6
		WHERE id = $7 AND application_id = $8;
	`
)
