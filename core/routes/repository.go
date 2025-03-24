package routes

import (
	"database/sql"
	userRepo "github.com/winartodev/apollo/modules/user/repositories"
)

type RepositoryDependency struct {
	DB *sql.DB
}

type Repository struct {
	UserRepository userRepo.UserRepositoryItf
}

func NewRepository(dependency RepositoryDependency) *Repository {
	newUserRepository := userRepo.NewUserRepository(dependency.DB)

	return &Repository{
		UserRepository: newUserRepository,
	}
}
