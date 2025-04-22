package repository

import (
	"context"
	"database/sql"
	"fmt"
	applicationEntity "github.com/winartodev/apollo/modules/application/entities"
	"strings"
)

type ApplicationScopeRepositoryItf interface {
	BulkCreateApplicationScopeTxDB(ctx context.Context, tx *sql.Tx, data []applicationEntity.ApplicationScope) (id *int64, err error)
}

type ApplicationScopeRepository struct {
	DB *sql.DB
}

func NewApplicationScopeRepository(repository ApplicationScopeRepository) ApplicationScopeRepositoryItf {
	return &ApplicationScopeRepository{
		DB: repository.DB,
	}
}

func (ar *ApplicationScopeRepository) BulkCreateApplicationScopeTxDB(ctx context.Context, tx *sql.Tx, data []applicationEntity.ApplicationScope) (id *int64, err error) {
	if data == nil || len(data) <= 0 {
		return nil, nil
	}

	valueStrings := make([]string, 0, len(data))
	valueArgs := make([]interface{}, 0, len(data)*10)

	for i, scope := range data {
		pos := i * 7
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			pos+1, pos+2, pos+3, pos+4, pos+5, pos+6, pos+7))
		valueArgs = append(valueArgs, scope.ApplicationID)
		valueArgs = append(valueArgs, scope.ScopeID)
		valueArgs = append(valueArgs, scope.IsActive)
		valueArgs = append(valueArgs, scope.CreatedBy)
		valueArgs = append(valueArgs, scope.UpdatedBy)
		valueArgs = append(valueArgs, scope.CreatedAt.Unix())
		valueArgs = append(valueArgs, scope.UpdatedAt.Unix())
	}

	query := fmt.Sprintf(CreateApplicationScopeQuery, strings.Join(valueStrings, ","))

	var stmt *sql.Stmt
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, query)
	} else {
		stmt, err = ar.DB.PrepareContext(ctx, query)
	}

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, valueArgs...).Scan(&id)
	if err != nil {
		return nil, err
	}

	return id, nil
}
