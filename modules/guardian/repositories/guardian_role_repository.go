package repositories

import (
	"context"
	"database/sql"
	"github.com/winartodev/apollo/core/helpers"
	guardianRole "github.com/winartodev/apollo/modules/guardian/entities"
)

type GuardianRoleRepositoryItf interface {
	CreateRole(ctx context.Context, data guardianRole.GuardianRole) (id int64, err error)
	GetRoles(ctx context.Context, appID int64, paginate *helpers.Paginate) (res []guardianRole.GuardianRole, err error)
	GetRoleByID(ctx context.Context, appID int64, id int64) (res *guardianRole.GuardianRole, err error)
	GetRoleBySlug(ctx context.Context, appID int64, slug string) (data *guardianRole.GuardianRole, err error)
	UpdateRole(ctx context.Context, appID int64, id int64, data guardianRole.GuardianRole) (err error)
}

type GuardianRoleRepository struct {
	DB *sql.DB
}

func NewGuardianRoleRepository(repository GuardianRoleRepository) GuardianRoleRepositoryItf {
	return &GuardianRoleRepository{
		DB: repository.DB,
	}
}

func (r *GuardianRoleRepository) CreateRole(ctx context.Context, data guardianRole.GuardianRole) (id int64, err error) {
	stmt, err := r.DB.PrepareContext(ctx, InsertIntoGuardianRoleQueryDB)
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	err = stmt.QueryRowContext(ctx,
		data.ApplicationID,
		data.Slug,
		data.Name,
		data.Description,
		data.CreatedAt.Unix(),
		data.UpdatedAt.Unix(),
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *GuardianRoleRepository) GetRoles(ctx context.Context, appID int64, paginate *helpers.Paginate) (res []guardianRole.GuardianRole, err error) {
	rows, err := r.DB.QueryContext(ctx, GetRolesQueryDB, appID)
	if err != nil {
		return nil, err
	}

	var createdAtUnix int64
	var updatedAtUnix int64

	res = make([]guardianRole.GuardianRole, 0)
	defer rows.Close()
	for rows.Next() {
		var data guardianRole.GuardianRole
		err = rows.Scan(
			&data.ID,
			&data.ApplicationID,
			&data.Slug,
			&data.Name,
			&data.Description,
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

func (r *GuardianRoleRepository) GetRoleByID(ctx context.Context, appID int64, id int64) (res *guardianRole.GuardianRole, err error) {
	var createdAtUnix int64
	var updatedAtUnix int64

	res = &guardianRole.GuardianRole{}
	err = r.DB.QueryRowContext(ctx, GetRoleByIDQueryDB,
		appID,
		id,
	).Scan(
		&res.ID,
		&res.ApplicationID,
		&res.Slug,
		&res.Name,
		&res.Description,
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

func (r *GuardianRoleRepository) GetRoleBySlug(ctx context.Context, appID int64, slug string) (res *guardianRole.GuardianRole, err error) {
	var createdAtUnix int64
	var updatedAtUnix int64

	res = &guardianRole.GuardianRole{}
	err = r.DB.QueryRowContext(ctx, GetRoleBySlugQueryDB,
		appID,
		slug,
	).Scan(
		&res.ID,
		&res.ApplicationID,
		&res.Slug,
		&res.Name,
		&res.Description,
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

func (r *GuardianRoleRepository) UpdateRole(ctx context.Context, appID int64, id int64, data guardianRole.GuardianRole) (err error) {
	panic("un implemented")
}
