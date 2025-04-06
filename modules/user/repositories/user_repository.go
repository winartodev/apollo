package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/winartodev/apollo/core/helpers"
	"github.com/winartodev/apollo/modules/user/entities"
	"time"
)

type UserRepositoryItf interface {
	CreateUserDB(ctx context.Context, user *entities.User) (id int64, err error)
	UpdateRefreshTokenByIDDB(ctx context.Context, id int64, refreshToken *string) (err error)
	UpdatePasswordByIDDB(ctx context.Context, id int64, password *string) error
	GetUserByIDDB(ctx context.Context, id int64) (res *entities.User, err error)
	GetUserByEmailDB(ctx context.Context, email string) (res *entities.User, err error)
	GetRefreshTokenByIDDB(ctx context.Context, id int64) (res *string, err error)
	GetUserPasswordByEmailDB(ctx context.Context, email string) (res *string, err error)
	GetUserRoleByIDDB(ctx context.Context, id int64) (res *entities.UserRole, err error)
	IsRefreshTokenExistByIDDB(ctx context.Context, id int64) (exists bool, err error)
	IsUserExistsDB(ctx context.Context, data *entities.UserUniqueField) (res *entities.UserUniqueFieldExists, err error)
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepositoryItf {
	return &UserRepository{
		DB: db,
	}
}

func (ur *UserRepository) CreateUserDB(ctx context.Context, user *entities.User) (id int64, err error) {
	now := time.Now()
	createdAtUnix := now.Unix()
	updatedAtUnix := now.Unix()

	tx, err := ur.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	stmt, err := tx.PrepareContext(ctx, InsertUserDBQuery)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = stmt.QueryRowContext(ctx,
		user.UUID,
		user.Email,
		user.PhoneNumber,
		user.Username,
		user.FirstName,
		user.LastName,
		user.ProfilePicture,
		user.Password,
		user.RefreshToken,
		user.IsEmailVerified,
		user.IsPhoneVerified,
		createdAtUnix,
		updatedAtUnix,
	).Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return id, err
}

func (ur *UserRepository) GetUserByIDDB(ctx context.Context, id int64) (res *entities.User, err error) {
	query := fmt.Sprintf("%s WHERE id = $1", GetUserQueryDB)

	var lastLoginUnix int64
	var createdAtUnix int64
	var updatedAtUnix int64

	res = &entities.User{}
	err = ur.DB.QueryRowContext(ctx, query, id).
		Scan(
			&res.ID,
			&res.UUID,
			&res.Email,
			&res.PhoneNumber,
			&res.Username,
			&res.FirstName,
			&res.LastName,
			&res.ProfilePicture,
			&res.IsEmailVerified,
			&res.IsPhoneVerified,
			&lastLoginUnix,
			&createdAtUnix,
			&updatedAtUnix,
		)

	if err != nil {
		return nil, err
	}

	res.LastLogin = helpers.FormatUnixTime(lastLoginUnix)
	res.CreatedAt = helpers.FormatUnixTime(createdAtUnix)
	res.UpdatedAt = helpers.FormatUnixTime(updatedAtUnix)

	return res, err
}

func (ur *UserRepository) GetUserByEmailDB(ctx context.Context, email string) (res *entities.User, err error) {
	query := fmt.Sprintf("%s WHERE email = $1", GetUserQueryDB)

	var lastLoginUnix int64
	var createdAtUnix int64
	var updatedAtUnix int64

	res = &entities.User{}
	err = ur.DB.QueryRowContext(ctx,
		query,
		email,
	).Scan(
		&res.ID,
		&res.UUID,
		&res.Email,
		&res.PhoneNumber,
		&res.Username,
		&res.FirstName,
		&res.LastName,
		&res.ProfilePicture,
		&res.IsEmailVerified,
		&res.IsPhoneVerified,
		&lastLoginUnix,
		&createdAtUnix,
		&updatedAtUnix,
	)
	if err != nil {
		return nil, err
	}

	res.LastLogin = helpers.FormatUnixTime(lastLoginUnix)
	res.CreatedAt = helpers.FormatUnixTime(createdAtUnix)
	res.UpdatedAt = helpers.FormatUnixTime(updatedAtUnix)

	return res, nil
}

func (ur *UserRepository) UpdateRefreshTokenByIDDB(ctx context.Context, id int64, refreshToken *string) (err error) {
	tx, err := ur.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, UpdateRefreshTokenByIDDBQuery)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmt.ExecContext(ctx, refreshToken, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) UpdatePasswordByIDDB(ctx context.Context, id int64, password *string) error {
	return nil
}

func (ur *UserRepository) GetRefreshTokenByIDDB(ctx context.Context, id int64) (res *string, err error) {
	err = ur.DB.QueryRowContext(ctx, GetRefreshTokenByIDDBQuery,
		id,
	).Scan(&res)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (ur *UserRepository) GetUserPasswordByEmailDB(ctx context.Context, email string) (res *string, err error) {
	err = ur.DB.QueryRowContext(ctx, GetUserPasswordByEmailDBQuery,
		email,
	).Scan(&res)

	if err != nil {
		return nil, err
	}

	return res, err
}

func (ur *UserRepository) IsRefreshTokenExistByIDDB(ctx context.Context, id int64) (exists bool, err error) {
	err = ur.DB.QueryRowContext(ctx, IsRefreshTokenIsExistsDBQuery,
		id,
	).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (ur *UserRepository) IsUserExistsDB(ctx context.Context, data *entities.UserUniqueField) (res *entities.UserUniqueFieldExists, err error) {
	res = &entities.UserUniqueFieldExists{}

	err = ur.DB.QueryRowContext(ctx, IsUserExistDBQuery,
		&data.Username,
		&data.Email,
		&data.PhoneNumber,
	).Scan(
		&res.IsUsernameExists,
		&res.IsEmailExists,
		&res.IsPhoneExists,
	)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (ur *UserRepository) GetUserRoleByIDDB(ctx context.Context, id int64) (res *entities.UserRole, err error) {
	stmt, err := ur.DB.PrepareContext(ctx, GetUserRoleByIDQuery)
	if err != nil {
		return nil, err
	}

	res = &entities.UserRole{}
	defer stmt.Close()
	err = stmt.QueryRowContext(ctx, id).
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
