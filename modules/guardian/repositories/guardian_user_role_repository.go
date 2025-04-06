package repositories

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	guardianEntity "github.com/winartodev/apollo/modules/guardian/entities"
)

type GuardianUserRoleRepositoryItf interface {
	GetRoleByUserID(ctx context.Context, userID int64) (res *guardianEntity.GuardianRole, err error)
}

type GuardianUserRoleRepository struct {
	DB    *sql.DB
	Redis *redis.Client
}

func NewGuardianUserRoleRepository(repository GuardianUserRoleRepository) GuardianUserRoleRepositoryItf {
	return &GuardianUserRoleRepository{
		DB:    repository.DB,
		Redis: repository.Redis,
	}
}

func (r *GuardianUserRoleRepository) GetRoleByUserID(ctx context.Context, userID int64) (res *guardianEntity.GuardianRole, err error) {
	stmt, err := r.DB.PrepareContext(ctx, GetUserRoleByUserIDQuery)
	if err != nil {
		return nil, err
	}

	res = &guardianEntity.GuardianRole{}
	defer stmt.Close()
	err = stmt.QueryRowContext(ctx, userID).
		Scan(
			&res.ID,
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
