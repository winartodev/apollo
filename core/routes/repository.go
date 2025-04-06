package routes

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	applicationRepo "github.com/winartodev/apollo/modules/application/repositories"
	authRepo "github.com/winartodev/apollo/modules/auth/repositories"
	guardianRepo "github.com/winartodev/apollo/modules/guardian/repositories"
	userRepo "github.com/winartodev/apollo/modules/user/repositories"
)

type RepositoryDependency struct {
	DB    *sql.DB
	Redis *redis.Client
}

type Repository struct {
	UserRepository                   userRepo.UserRepositoryItf
	VerificationRepository           authRepo.VerificationRepositoryItf
	UserApplicationRepository        applicationRepo.UserApplicationRepositoryItf
	UserApplicationServiceRepository applicationRepo.UserApplicationServiceRepositoryItf
	GuardianUserRoleRepo             guardianRepo.GuardianUserRoleRepositoryItf
}

func NewRepository(dependency RepositoryDependency) *Repository {
	newVerificationRepo := authRepo.NewVerificationRepository(authRepo.VerificationRepository{
		Redis: dependency.Redis,
	})

	newUserRepository := userRepo.NewUserRepository(dependency.DB)

	newUserApplicationRepo := applicationRepo.NewUserApplicationRepository(applicationRepo.UserApplicationRepository{
		DB:    dependency.DB,
		Redis: dependency.Redis,
	})

	newUserApplicationServiceRepo := applicationRepo.NewUserApplicationService(applicationRepo.UserApplicationServiceRepository{
		DB:    dependency.DB,
		Redis: dependency.Redis,
	})

	newGuardianUserRoleRepo := guardianRepo.NewGuardianUserRoleRepository(guardianRepo.GuardianUserRoleRepository{
		DB:    dependency.DB,
		Redis: dependency.Redis,
	})

	return &Repository{
		VerificationRepository:           newVerificationRepo,
		UserRepository:                   newUserRepository,
		UserApplicationRepository:        newUserApplicationRepo,
		UserApplicationServiceRepository: newUserApplicationServiceRepo,
		GuardianUserRoleRepo:             newGuardianUserRoleRepo,
	}
}
