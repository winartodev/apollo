package repositories

import (
	"context"
	"database/sql"
)

type UserApplicationRepositoryItf interface {
	GetUserApplicationAccessDB(ctx context.Context, userID int64, applicationID int64, scopeID int64) (hasAccess bool, isAppActive bool, err error)
}

type UserApplicationRepository struct {
	DB *sql.DB
}

func NewUserApplicationRepository(repository UserApplicationRepository) UserApplicationRepositoryItf {
	return &UserApplicationRepository{
		DB: repository.DB,
	}
}

func (uar *UserApplicationRepository) GetUserApplicationAccessDB(ctx context.Context, userID int64, applicationID int64, scopeID int64) (hasAccess bool, isAppActive bool, err error) {
	err = uar.DB.QueryRowContext(ctx,
		GetUserApplicationAccessQuery,
		userID,
		applicationID,
		scopeID,
		applicationID,
	).Scan(&hasAccess, &isAppActive)
	if err != nil {
		return false, false, err
	}

	return hasAccess, isAppActive, nil
}
