package repositories

import (
	"database/sql"
	"github.com/winartodev/apollo/modules/application/entities"
	"golang.org/x/net/context"
)

type ApplicationServiceRepositoryItf interface {
	GetApplicationServiceBySlugDB(ctx context.Context, slug string) (res *entities.ApplicationService, err error)
}

type ApplicationServiceRepository struct {
	DB *sql.DB
}

func NewApplicationServiceRepository(repository ApplicationServiceRepository) ApplicationServiceRepositoryItf {
	return &ApplicationServiceRepository{
		DB: repository.DB,
	}
}

func (r *ApplicationServiceRepository) GetApplicationServiceBySlugDB(ctx context.Context, slug string) (res *entities.ApplicationService, err error) {
	stmt, err := r.DB.PrepareContext(ctx, GetApplicationServiceBySlug)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	res = &entities.ApplicationService{}
	err = stmt.QueryRowContext(ctx, slug).
		Scan(
			&res.ID,
			&res.ApplicationID,
			&res.Scope,
			&res.Slug,
			&res.Name,
			&res.IsActive,
		)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return res, err
}
