package repository

import (
	"context"
	"database/sql"
	applicationEntity "github.com/winartodev/apollo/modules/application/entities"
)

type ServiceRepositoryItf interface {
	CreateServiceDB(ctx context.Context, data applicationEntity.Service) (id *int64, err error)
	GetServiceBySlugDB(ctx context.Context, slug string) (res *applicationEntity.Service, err error)
}

type ServiceRepository struct {
	DB *sql.DB
}

func NewServiceRepository(repository ServiceRepository) ServiceRepositoryItf {
	return &ServiceRepository{
		DB: repository.DB,
	}
}

func (sr *ServiceRepository) CreateServiceDB(ctx context.Context, data applicationEntity.Service) (id *int64, err error) {
	stmt, err := sr.DB.PrepareContext(ctx, InsertIntoServiceQuery)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	createdAtUnix := data.CreatedAt.Unix()
	updatedAtUnix := data.UpdatedAt.Unix()
	err = stmt.QueryRowContext(ctx,
		data.Slug,
		data.Name,
		data.Description,
		data.IsActive,
		data.CreatedBy,
		data.UpdatedBy,
		createdAtUnix,
		updatedAtUnix,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (sr *ServiceRepository) GetServiceBySlugDB(ctx context.Context, slug string) (res *applicationEntity.Service, err error) {
	var createdAtUnix int64
	var updatedAtUnix int64

	res = &applicationEntity.Service{}
	err = sr.DB.QueryRowContext(ctx,
		GetServiceBySlugQuery,
		slug,
	).Scan(
		&res.ID,
		&res.Slug,
		&res.Name,
		&res.Description,
		&res.IsActive,
		&res.CreatedBy,
		&res.UpdatedBy,
		&createdAtUnix,
		&updatedAtUnix,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return res, nil
}
