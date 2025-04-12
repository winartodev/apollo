package routes

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	"github.com/winartodev/apollo/core/helpers"
	applicationRepo "github.com/winartodev/apollo/modules/application/repository"
	authRepo "github.com/winartodev/apollo/modules/auth/repositories"
	userRepo "github.com/winartodev/apollo/modules/user/repositories"
)

type RepositoryDependency struct {
	DB    *sql.DB
	Redis *redis.Client
}

type Repository struct {
	UserRepository             userRepo.UserRepositoryItf
	UserApplicationRepository  userRepo.UserApplicationRepositoryItf
	VerificationRepository     authRepo.VerificationRepositoryItf
	ServiceRepository          applicationRepo.ServiceRepositoryItf
	ApplicationRepository      applicationRepo.ApplicationRepositoryItf
	ApplicationScopeRepository applicationRepo.ApplicationScopeRepositoryItf
	DBTransact                 helpers.DBTransactItf
}

func NewRepository(dependency RepositoryDependency) *Repository {
	dbtx := helpers.NewDBTransact(dependency.DB)

	newVerificationRepo := authRepo.NewVerificationRepository(authRepo.VerificationRepository{
		Redis: dependency.Redis,
	})

	newUserRepository := userRepo.NewUserRepository(dependency.DB)
	newUserApplicationRepository := userRepo.NewUserApplicationRepository(userRepo.UserApplicationRepository{
		DB: dependency.DB,
	})

	newServiceRepo := applicationRepo.NewServiceRepository(applicationRepo.ServiceRepository{
		DB: dependency.DB,
	})

	newApplicationRepo := applicationRepo.NewApplicationRepository(applicationRepo.ApplicationRepository{
		DB: dependency.DB,
	})

	newApplicationScopeRepo := applicationRepo.NewApplicationScopeRepository(applicationRepo.ApplicationScopeRepository{
		DB: dependency.DB,
	})

	return &Repository{
		VerificationRepository:     newVerificationRepo,
		UserRepository:             newUserRepository,
		UserApplicationRepository:  newUserApplicationRepository,
		ServiceRepository:          newServiceRepo,
		ApplicationRepository:      newApplicationRepo,
		ApplicationScopeRepository: newApplicationScopeRepo,
		DBTransact:                 dbtx,
	}
}
