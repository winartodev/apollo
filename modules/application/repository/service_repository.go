package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/winartodev/apollo/core/helpers"
	applicationEntity "github.com/winartodev/apollo/modules/application/entities"
)

type ServiceRepositoryItf interface {
	CreateServiceDB(ctx context.Context, data applicationEntity.Service) (id *int64, err error)
	GetServiceBySlugDB(ctx context.Context, slug string) (res *applicationEntity.Service, err error)
	GetServiceByIDDB(ctx context.Context, id int64) (res *applicationEntity.Service, err error)
	GetServicesDB(ctx context.Context, filter helpers.Paginate) (res []applicationEntity.Service, err error)
	GetTotalServiceDB(ctx context.Context) (res int64, err error)
	UpdateServiceByIDDB(ctx context.Context, id int64, data applicationEntity.Service) (err error)
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
		&res.ApplicationID,
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

func (sr *ServiceRepository) GetServiceByIDDB(ctx context.Context, id int64) (res *applicationEntity.Service, err error) {
	var createdAtUnix int64
	var updatedAtUnix int64
	res = &applicationEntity.Service{}

	err = sr.DB.QueryRowContext(ctx,
		GetServiceByIDQuery,
		id,
	).Scan(
		&res.ID,
		&res.ApplicationID,
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

func (sr *ServiceRepository) GetServicesDB(ctx context.Context, filter helpers.Paginate) (res []applicationEntity.Service, err error) {
	query := fmt.Sprintf("%s %s", GetServicesQuery, filter.BuildSQLQuery())
	rows, err := sr.DB.QueryContext(ctx,
		query,
		*filter.Limit,
		*filter.Offset,
	)
	if err != nil {
		return nil, err
	}

	var createdAtUnix int64
	var updatedAtUnix int64
	res = make([]applicationEntity.Service, 0)
	for rows.Next() {
		data := applicationEntity.Service{}
		err = rows.Scan(
			&data.ID,
			&data.ApplicationID,
			&data.Slug,
			&data.Name,
			&data.Description,
			&data.IsActive,
			&data.CreatedBy,
			&data.UpdatedBy,
			&createdAtUnix,
			&updatedAtUnix)
		if err != nil {
			return nil, err
		}

		data.CreatedAt = helpers.FormatUnixTime(createdAtUnix)
		data.UpdatedAt = helpers.FormatUnixTime(updatedAtUnix)

		res = append(res, data)
	}

	return res, nil
}

func (sr *ServiceRepository) GetTotalServiceDB(ctx context.Context) (res int64, err error) {
	err = sr.DB.QueryRowContext(ctx, CountServiceQuery).Scan(&res)
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (sr *ServiceRepository) UpdateServiceByIDDB(ctx context.Context, id int64, data applicationEntity.Service) (err error) {
	updatedAtUnix := data.UpdatedAt.Unix()

	stmt, err := sr.DB.PrepareContext(ctx, UpdateServiceQuery)
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.ExecContext(ctx,
		data.Slug,
		data.Name,
		data.Description,
		data.IsActive,
		data.UpdatedBy,
		updatedAtUnix,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}
