package repositories

import (
	"context"
	"database/sql"
	"github.com/winartodev/apollo/modules/user/entities"
	"time"
)

type UserApplicationRepositoryItf interface {
	CreateUserApplicationDB(ctx context.Context, app *entities.UserApplication) (err error)
	GetUserApplicationsByUserIDDB(ctx context.Context, userID int64) (res []entities.UserApplicationResponse, err error)
	GetUserApplicationByUserIDAndApplicationSlugDB(ctx context.Context, userID int64, applicationSlug string) (res *entities.UserApplicationResponse, err error)
}

type UserApplicationRepository struct {
	DB *sql.DB
}

func NewUserApplicationRepository(repository UserApplicationRepository) UserApplicationRepositoryItf {
	return &UserApplicationRepository{
		DB: repository.DB,
	}
}

func (ur *UserApplicationRepository) CreateUserApplicationDB(ctx context.Context, data *entities.UserApplication) (err error) {
	now := time.Now()
	createdAtUnix := now.Unix()
	updatedAtUnix := now.Unix()

	tx, err := ur.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, InsertUserApplicationQueryDB)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmt.ExecContext(ctx,
		data.UserID,
		data.ApplicationID,
		data.CreatedAt,
		createdAtUnix,
		updatedAtUnix,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

func (ur *UserApplicationRepository) GetUserApplicationsByUserIDDB(ctx context.Context, userID int64) (res []entities.UserApplicationResponse, err error) {
	stmt, err := ur.DB.PrepareContext(ctx, GetUserApplicationsByUserIDQueryDB)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, err
	}

	res = make([]entities.UserApplicationResponse, 0)
	for rows.Next() {
		var app entities.UserApplicationResponse

		err = rows.Scan(
			&app.ID,
			&app.Slug,
			&app.Name,
			&app.IsActive,
		)

		if err != nil {
			return nil, err
		}

		res = append(res, app)
	}

	return res, err
}

func (ur *UserApplicationRepository) GetUserApplicationByUserIDAndApplicationSlugDB(ctx context.Context, userID int64, applicationSlug string) (res *entities.UserApplicationResponse, err error) {
	stmt, err := ur.DB.PrepareContext(ctx, GetUserApplicationByUserIDAndApplicationSlugQuery)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	res = &entities.UserApplicationResponse{}
	err = stmt.QueryRowContext(ctx, userID, applicationSlug).Scan(
		&res.ID,
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

	return res, nil
}
