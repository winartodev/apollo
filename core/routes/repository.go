package routes

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	applicationRepo "github.com/winartodev/apollo/modules/application/repository"
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
	ServiceRepository      applicationRepo.ServiceRepositoryItf
}

func NewRepository(dependency RepositoryDependency) *Repository {
	newVerificationRepo := authRepo.NewVerificationRepository(authRepo.VerificationRepository{
		Redis: dependency.Redis,
	})

	newUserRepository := userRepo.NewUserRepository(dependency.DB)

	newServiceRepo := applicationRepo.NewServiceRepository(applicationRepo.ServiceRepository{
		DB: dependency.DB,
	})

	return &Repository{
		VerificationRepository: newVerificationRepo,
		UserRepository:         newUserRepository,
		ServiceRepository:      newServiceRepo,
	}
}
