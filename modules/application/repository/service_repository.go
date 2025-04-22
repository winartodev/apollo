package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/winartodev/apollo/core/helpers"
	applicationEntity "github.com/winartodev/apollo/modules/application/entities"
	"strings"
)

type ServiceRepositoryItf interface {
	BulkInsertServiceTxDB(ctx context.Context, tx *sql.Tx, data []applicationEntity.Service) (ids []int64, err error)
	GetServiceBySlugDB(ctx context.Context, applicationID int64, slug string) (res *applicationEntity.Service, err error)
	GetServiceByIDDB(ctx context.Context, applicationID int64, serviceID int64) (res *applicationEntity.Service, err error)
	GetServicesDB(ctx context.Context, applicationID int64, filter helpers.Paginate) (res []applicationEntity.Service, err error)
	GetTotalServiceDB(ctx context.Context, applicationID int64) (res int64, err error)
	UpdateServiceByIDDB(ctx context.Context, applicationID int64, serviceID int64, data applicationEntity.Service) (err error)
}

type ServiceRepository struct {
	DB *sql.DB
}

func NewServiceRepository(repository ServiceRepository) ServiceRepositoryItf {
	return &ServiceRepository{
		DB: repository.DB,
	}
}

func (sr *ServiceRepository) BulkInsertServiceTxDB(ctx context.Context, tx *sql.Tx, data []applicationEntity.Service) (ids []int64, err error) {
	if data == nil || len(data) <= 0 {
		return nil, nil
	}

	valueStrings := make([]string, 0, len(data))
	valueArgs := make([]interface{}, 0, len(data)*10)

	for i, service := range data {
		pos := i * 9
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			pos+1, pos+2, pos+3, pos+4, pos+5, pos+6, pos+7, pos+8, pos+9))
		valueArgs = append(valueArgs, service.Slug)
		valueArgs = append(valueArgs, service.Name)
		valueArgs = append(valueArgs, service.Description)
		valueArgs = append(valueArgs, service.IsActive)
		valueArgs = append(valueArgs, service.CreatedBy)
		valueArgs = append(valueArgs, service.UpdatedBy)
		valueArgs = append(valueArgs, service.CreatedAt.Unix())
		valueArgs = append(valueArgs, service.UpdatedAt.Unix())
		valueArgs = append(valueArgs, *service.ApplicationID)
	}

	query := fmt.Sprintf(InsertIntoServiceQuery, strings.Join(valueStrings, ","))
	var stmt *sql.Stmt
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, query)
	} else {
		stmt, err = sr.DB.PrepareContext(ctx, query)
	}

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, valueArgs...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (sr *ServiceRepository) GetServiceBySlugDB(ctx context.Context, applicationID int64, slug string) (res *applicationEntity.Service, err error) {
	var createdAtUnix int64
	var updatedAtUnix int64

	res = &applicationEntity.Service{}
	err = sr.DB.QueryRowContext(ctx,
		GetServiceBySlugQuery,
		slug,
		applicationID,
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

func (sr *ServiceRepository) GetServiceByIDDB(ctx context.Context, applicationID int64, serviceID int64) (res *applicationEntity.Service, err error) {
	var createdAtUnix int64
	var updatedAtUnix int64
	res = &applicationEntity.Service{}

	err = sr.DB.QueryRowContext(ctx,
		GetServiceByIDQuery,
		serviceID,
		applicationID,
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

func (sr *ServiceRepository) GetServicesDB(ctx context.Context, applicationID int64, filter helpers.Paginate) (res []applicationEntity.Service, err error) {
	query := fmt.Sprintf("%s WHERE application_id = $1 ORDER BY id ASC LIMIT $2 OFFSET $3", GetServicesQuery)
	rows, err := sr.DB.QueryContext(ctx,
		query,
		applicationID,
		*filter.Limit,
		*filter.Offset,
	)
	if err != nil {
		return nil, err
	}

	var createdAtUnix int64
	var updatedAtUnix int64
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

func (sr *ServiceRepository) GetTotalServiceDB(ctx context.Context, applicationID int64) (res int64, err error) {
	err = sr.DB.QueryRowContext(ctx, CountServiceQuery, applicationID).Scan(&res)
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (sr *ServiceRepository) UpdateServiceByIDDB(ctx context.Context, applicationID int64, serviceID int64, data applicationEntity.Service) (err error) {
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
		serviceID,
		applicationID,
	)
	if err != nil {
		return err
	}

	return nil
}
