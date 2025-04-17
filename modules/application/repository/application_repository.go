package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/winartodev/apollo/core/helpers"
	applicationEntity "github.com/winartodev/apollo/modules/application/entities"
)

type ApplicationRepositoryItf interface {
	CreateApplicationDB(ctx context.Context, data applicationEntity.Application) (id int64, err error)
	UpdateApplicationByIDDB(ctx context.Context, id int64, data applicationEntity.Application) (err error)
	GetApplicationsDB(ctx context.Context, filter helpers.Paginate) (res []applicationEntity.Application, err error)
	GetApplicationByIDDB(ctx context.Context, id int64) (res *applicationEntity.Application, err error)
	GetApplicationBySlugDB(ctx context.Context, slug string) (res *applicationEntity.Application, err error)
	GetTotalApplicationDB(ctx context.Context) (total int64, err error)
}

type ApplicationRepository struct {
	DB *sql.DB
}

func NewApplicationRepository(repository ApplicationRepository) ApplicationRepositoryItf {
	return &ApplicationRepository{
		DB: repository.DB,
	}
}

func (ar *ApplicationRepository) CreateApplicationDB(ctx context.Context, data applicationEntity.Application) (id int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (ar *ApplicationRepository) GetApplicationsDB(ctx context.Context, filter helpers.Paginate) (res []applicationEntity.Application, err error) {
	query := fmt.Sprintf("%s %s", GetApplicationQuery, filter.BuildSQLQuery())
	rows, err := ar.DB.QueryContext(ctx,
		query,
		*filter.Limit,
		*filter.Offset)
	if err != nil {
		return nil, err
	}

	res = make([]applicationEntity.Application, 0)
	for rows.Next() {
		var createdAtUnix int64
		var updatedAtUnix int64
		var data applicationEntity.Application
		err = rows.Scan(
			&data.ID,
			&data.Slug,
			&data.Name,
			&data.IsActive,
			&data.CreatedBy,
			&data.UpdatedBy,
			&createdAtUnix,
			&updatedAtUnix,
		)
		if err != nil {
			return nil, err
		}

		data.CreatedAt = helpers.FormatUnixTime(createdAtUnix)
		data.UpdatedAt = helpers.FormatUnixTime(updatedAtUnix)

		res = append(res, data)
	}

	return res, nil
}

func (ar *ApplicationRepository) GetApplicationByIDDB(ctx context.Context, id int64) (res *applicationEntity.Application, err error) {
	var createdAtUnix int64
	var updatedAtUnix int64
	var data applicationEntity.Application
	err = ar.DB.QueryRowContext(ctx,
		GetApplicationByIDQuery,
		id,
	).Scan(
		&data.ID,
		&data.Slug,
		&data.Name,
		&data.IsActive,
		&data.CreatedBy,
		&data.UpdatedBy,
		&createdAtUnix,
		&updatedAtUnix,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	data.CreatedAt = helpers.FormatUnixTime(createdAtUnix)
	data.UpdatedAt = helpers.FormatUnixTime(updatedAtUnix)

	return &data, nil
}

func (ar *ApplicationRepository) GetApplicationBySlugDB(ctx context.Context, slug string) (res *applicationEntity.Application, err error) {
	//TODO implement me
	panic("implement me")
}

func (ar *ApplicationRepository) UpdateApplicationByIDDB(ctx context.Context, id int64, data applicationEntity.Application) (err error) {
	//TODO implement me
	panic("implement me")
}

func (ar *ApplicationRepository) GetTotalApplicationDB(ctx context.Context) (total int64, err error) {
	err = ar.DB.QueryRowContext(ctx, CountApplicationQuery).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
