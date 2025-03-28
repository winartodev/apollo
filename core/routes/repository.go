package routes

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	authRepo "github.com/winartodev/apollo/modules/auth/repositories"
	userRepo "github.com/winartodev/apollo/modules/user/repositories"
)

type RepositoryDependency struct {
	DB    *sql.DB
	Redis *redis.Client
}

type Repository struct {
	UserRepository         userRepo.UserRepositoryItf
	VerificationRepository authRepo.VerificationRepositoryItf
}

func NewRepository(dependency RepositoryDependency) *Repository {
	newVerificationRepo := authRepo.NewVerificationRepository(authRepo.VerificationRepository{
		Redis: dependency.Redis,
	})
	newUserRepository := userRepo.NewUserRepository(dependency.DB)

	return &Repository{
		VerificationRepository: newVerificationRepo,
		UserRepository:         newUserRepository,
	}
}
