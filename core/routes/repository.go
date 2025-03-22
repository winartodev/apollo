package routes

import (
	"database/sql"
)

type RepositoryDependency struct {
	DB *sql.DB
}

type Repository struct {
}

func NewRepository(dependency RepositoryDependency) *Repository {
	return &Repository{}
}
