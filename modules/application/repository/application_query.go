package repository

const (
	InsertApplicationQuery = `
		INSERT INTO applications (slug, name, is_active, created_by, updated_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id
	`

	GetApplicationQuery = `
		SELECT 
		    id, 
		    slug, 
		    name, 
		    is_active, 
		    created_by, 
		    updated_by, 
		    created_at, 
		    updated_at 
		FROM applications
	`

	GetApplicationByIDQuery = GetApplicationQuery + ` WHERE id = $1`

	GetApplicationBySlugQuery = GetApplicationQuery + ` WHERE slug = $1`

	CountApplicationQuery = `SELECT COUNT(*) FROM applications;`
)
