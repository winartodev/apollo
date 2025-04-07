package repositories

import (
	"context"
	"database/sql"
	"github.com/winartodev/apollo/modules/user/entities"
)

type UserRoleRepositoryItf interface {
	GetUserRoleByIDDB(ctx context.Context, id int64, appID int64) (res *entities.UserRoleResponse, err error)
}

type UserRoleRepository struct {
	DB *sql.DB
}

func NewUserRoleRepository(repository UserRoleRepository) UserRoleRepositoryItf {
	return &UserRoleRepository{
		DB: repository.DB,
	}
}

func (ur *UserRoleRepository) GetUserRoleByIDDB(ctx context.Context, id int64, appID int64) (res *entities.UserRoleResponse, err error) {
	stmt, err := ur.DB.PrepareContext(ctx, GetUserRoleByIDQuery)
	if err != nil {
		return nil, err
	}

	res = &entities.UserRoleResponse{}
	defer stmt.Close()
	err = stmt.QueryRowContext(ctx, id, appID).
		Scan(
			&res.RoleID,
			&res.Slug,
			&res.Name,
		)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return res, err
}
