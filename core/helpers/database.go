package helpers

import (
	"context"
	"database/sql"
)

type DBTransactItf interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	Commit() error
	Rollback() error
	GetTx() *sql.Tx
}

type DBTransact struct {
	DB *sql.DB
	tx *sql.Tx
}

func (dtx *DBTransact) BeginTx(ctx context.Context) (*sql.Tx, error) {
	if dtx.tx != nil {
		return nil, sql.ErrTxDone
	}

	tx, err := dtx.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	dtx.tx = tx
	return tx, nil
}

func (dtx *DBTransact) Commit() error {
	if dtx.tx == nil {
		return sql.ErrTxDone
	}

	err := dtx.tx.Commit()
	dtx.tx = nil
	return err
}

func (dtx *DBTransact) Rollback() error {
	if dtx.tx == nil {
		return sql.ErrTxDone
	}

	err := dtx.tx.Rollback()
	dtx.tx = nil
	return err
}

func (dtx *DBTransact) GetTx() *sql.Tx {
	return dtx.tx
}

func NewDBTransact(db *sql.DB) DBTransactItf {
	return &DBTransact{
		DB: db,
	}
}
