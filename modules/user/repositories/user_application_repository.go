package repositories

import (
	"context"
	"database/sql"
	"fmt"
	userEntity "github.com/winartodev/apollo/modules/user/entities"
	"strings"
)

type UserApplicationRepositoryItf interface {
	GetUserApplicationAccessDB(ctx context.Context, userID int64, applicationID int64, scopeID int64) (hasAccess bool, isAppActive bool, err error)
	BulkInsertUserApplicationTxDB(ctx context.Context, tx *sql.Tx, data []userEntity.UserApplication) (res sql.Result, err error)
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

func (uar *UserApplicationRepository) BulkInsertUserApplicationTxDB(ctx context.Context, tx *sql.Tx, data []userEntity.UserApplication) (res sql.Result, err error) {
	if data == nil || len(data) <= 0 {
		return nil, nil
	}

	valueStrings := make([]string, 0, len(data))
	valueArgs := make([]interface{}, 0, len(data)*10)

	for i, userApp := range data {
		pos := i * 2
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", pos+1, pos+2))
		valueArgs = append(valueArgs, userApp.UserID)
		valueArgs = append(valueArgs, userApp.ApplicationScopeID)
	}

	query := fmt.Sprintf(InsertIntoUserApplication, strings.Join(valueStrings, ","))

	var stmt *sql.Stmt
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, query)
	} else {
		stmt, err = uar.DB.PrepareContext(ctx, query)
	}

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	res, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		return nil, err
	}

	return res, nil
}
