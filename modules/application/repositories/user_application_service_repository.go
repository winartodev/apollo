package repositories

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	applicationEntity "github.com/winartodev/apollo/modules/application/entities"
)

type UserApplicationServiceRepositoryItf interface {
	GetApplicationServiceAccess(ctx context.Context, appID int64, userID int64) (res []applicationEntity.UserApplicationServiceResponse, err error)
}

type UserApplicationServiceRepository struct {
	DB    *sql.DB
	Redis *redis.Client
}

func NewUserApplicationService(repository UserApplicationServiceRepository) UserApplicationServiceRepositoryItf {
	return &UserApplicationServiceRepository{
		DB:    repository.DB,
		Redis: repository.Redis,
	}
}

func (r *UserApplicationServiceRepository) GetApplicationServiceAccess(ctx context.Context, appID int64, userID int64) (res []applicationEntity.UserApplicationServiceResponse, err error) {
	stmt, err := r.DB.PrepareContext(ctx, GetUserApplicationServices)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID, appID)
	if err != nil {
		return nil, err
	}

	res = make([]applicationEntity.UserApplicationServiceResponse, 0)
	for rows.Next() {
		data := &applicationEntity.UserApplicationServiceResponse{}
		err = rows.Scan(
			&data.UserID,
			&data.AppID,
			&data.AppSlug,
			&data.AppServiceID,
			&data.AppServiceScope,
			&data.AppServiceSlug,
			&data.AppServiceName,
		)
		if err != nil {
			return nil, err
		}

		res = append(res, *data)
	}

	return res, nil
}
