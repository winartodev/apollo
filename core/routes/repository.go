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
	ApplicationServiceRepository applicationRepo.ApplicationServiceRepositoryItf
	UserRepository               userRepo.UserRepositoryItf
	UserRoleRepository           userRepo.UserRoleRepositoryItf
	UserApplicationRepository    userRepo.UserApplicationRepositoryItf
	GuardianPermissionRepository guardianRepo.GuardianPermissionRepositoryItf
	VerificationRepository       authRepo.VerificationRepositoryItf
}

func NewRepository(dependency RepositoryDependency) *Repository {
	newVerificationRepo := authRepo.NewVerificationRepository(authRepo.VerificationRepository{
		Redis: dependency.Redis,
	})

	newUserRepository := userRepo.NewUserRepository(dependency.DB)
	userUserRoleRepository := userRepo.NewUserRoleRepository(userRepo.UserRoleRepository{
		DB: dependency.DB,
	})

	newUserApplicationRepo := userRepo.NewUserApplicationRepository(userRepo.UserApplicationRepository{
		DB: dependency.DB,
	})

	newGuardianPermissionRepo := guardianRepo.NewGuardianPermissionRepository(guardianRepo.GuardianPermissionRepository{
		DB: dependency.DB,
	})

	newApplicationServiceRepo := applicationRepo.NewApplicationServiceRepository(applicationRepo.ApplicationServiceRepository{
		DB: dependency.DB,
	})

	return &Repository{
		ApplicationServiceRepository: newApplicationServiceRepo,
		VerificationRepository:       newVerificationRepo,
		UserRepository:               newUserRepository,
		UserRoleRepository:           userUserRoleRepository,
		GuardianPermissionRepository: newGuardianPermissionRepo,
		UserApplicationRepository:    newUserApplicationRepo,
	}
}
