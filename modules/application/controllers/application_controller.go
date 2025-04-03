package controllers

import (
	"context"
	applicationEntity "github.com/winartodev/apollo/modules/application/entities"
	"github.com/winartodev/apollo/modules/application/enums"
	applicationRepo "github.com/winartodev/apollo/modules/application/repositories"
	userController "github.com/winartodev/apollo/modules/user/controllers"
)

type ApplicationControllerItf interface {
	GetApplicationAccess(ctx context.Context, userID int64, applicationSlug string) (res *applicationEntity.ApplicationAccess, err error)
	GetApplicationService(ctx context.Context, userID int64, applicationID int64, appServiceSlug string) (res []applicationEntity.ApplicationService, err error)
}

type ApplicationController struct {
	UserController                   userController.UserControllerItf
	UserApplicationRepository        applicationRepo.UserApplicationRepositoryItf
	UserApplicationServiceRepository applicationRepo.UserApplicationServiceRepositoryItf
}

func NewApplicationController(controller ApplicationController) ApplicationControllerItf {
	return &ApplicationController{
		UserController:                   controller.UserController,
		UserApplicationRepository:        controller.UserApplicationRepository,
		UserApplicationServiceRepository: controller.UserApplicationServiceRepository,
	}
}

func (c *ApplicationController) GetApplicationAccess(ctx context.Context, userID int64, applicationSlug string) (res *applicationEntity.ApplicationAccess, err error) {
	res = &applicationEntity.ApplicationAccess{}

	userData, err := c.UserController.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if userData == nil {
		return nil, nil
	}

	userApplication, err := c.UserApplicationRepository.GetUserApplicationByUserIDAndApplicationSlugDB(ctx, userID, applicationSlug)
	if err != nil {
		return nil, err
	}

	res.User = *userData

	if userApplication != nil {
		application := applicationEntity.Application{
			ID:       userApplication.ID,
			Slug:     userApplication.Slug,
			Name:     userApplication.Name,
			IsActive: userApplication.IsActive,
		}

		res.Applications = append(res.Applications, application)
	}

	return res, err
}

func (c *ApplicationController) GetApplicationService(ctx context.Context, userID int64, applicationID int64, appServiceSlug string) (res []applicationEntity.ApplicationService, err error) {
	applicationServices, err := c.UserApplicationServiceRepository.GetApplicationServiceAccess(ctx, userID, applicationID)
	if err != nil {
		return nil, err
	}

	res = []applicationEntity.ApplicationService{}
	for _, appService := range applicationServices {
		var app = applicationEntity.ApplicationService{
			ID:    appService.AppServiceID,
			Slug:  appService.AppServiceSlug,
			Scope: enums.ApplicationScope(appService.AppServiceScope).ToString(),
			Name:  appService.AppServiceName,
		}

		res = append(res, app)
	}

	return res, err
}
