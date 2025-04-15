package repository

const (
	InsertIntoServiceQuery = `
		INSERT INTO services (slug, name, description, is_active,created_by, updated_by, created_at, updated_at) 
		VALUES ( $1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;
	`

	GetServiceBySlugQuery = `
		SELECT 
		    id, 
		    slug, 
		    name, 
		    description, 
		    is_active, 
		    created_by, 
		    updated_by, 
		    created_at, 
		    updated_at 
		FROM services 
		WHERE slug = $1;
	`

	GetServiceByIDQuery = `
		SELECT 
		    id, 
		    slug, 
		    name, 
		    description, 
		    is_active, 
		    created_by, 
		    updated_by, 
		    created_at, 
		    updated_at 
		FROM services 
		WHERE id = $1;
	`

	GetServicesQuery = `
		SELECT 
		    id, 
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

	CountServiceQuery = `SELECT COUNT(*) FROM services;`
)
