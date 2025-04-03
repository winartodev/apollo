package repositories

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	"github.com/winartodev/apollo/modules/application/entities"
)

type UserApplicationRepositoryItf interface {
	GetUserApplicationsByUserIDDB(ctx context.Context, userID int64) (res []entities.UserApplicationResponse, err error)
	GetUserApplicationByUserIDAndApplicationSlugDB(ctx context.Context, userID int64, applicationSlug string) (res *entities.UserApplicationResponse, err error)
}

type UserApplicationRepository struct {
	DB    *sql.DB
	Redis *redis.Client
}

func NewUserApplicationRepository(repository UserApplicationRepository) UserApplicationRepositoryItf {
	return &UserApplicationRepository{
		DB:    repository.DB,
		Redis: repository.Redis,
	}
}

func (r *UserApplicationRepository) GetUserApplicationsByUserIDDB(ctx context.Context, userID int64) (res []entities.UserApplicationResponse, err error) {

	stmt, err := r.DB.PrepareContext(ctx, GetUserApplicationsByUserIDQueryDB)
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

func (r *UserApplicationRepository) GetUserApplicationByUserIDAndApplicationSlugDB(ctx context.Context, userID int64, applicationSlug string) (res *entities.UserApplicationResponse, err error) {
	stmt, err := r.DB.PrepareContext(ctx, GetUserApplicationByUserIDAndApplicationSlugQuery)
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
