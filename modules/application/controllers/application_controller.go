package controllers

import (
	"context"
	applicationEntity "github.com/winartodev/apollo/modules/application/entities"
	"github.com/winartodev/apollo/modules/application/enums"
	applicationRepo "github.com/winartodev/apollo/modules/application/repositories"
)

type ApplicationControllerItf interface {
	GetApplicationAccess(ctx context.Context, userID int64, applicationSlug string) (res *applicationEntity.ApplicationAccess, err error)
	GetApplicationService(ctx context.Context, userID int64, applicationID int64, appServiceSlug string) (res *applicationEntity.ApplicationService, err error)
}

type ApplicationController struct {
	UserApplicationRepository        applicationRepo.UserApplicationRepositoryItf
	UserApplicationServiceRepository applicationRepo.UserApplicationServiceRepositoryItf
}

func NewApplicationController(controller ApplicationController) ApplicationControllerItf {
	return &ApplicationController{
		UserApplicationRepository:        controller.UserApplicationRepository,
		UserApplicationServiceRepository: controller.UserApplicationServiceRepository,
	}
}

func (c *ApplicationController) GetApplicationAccess(ctx context.Context, userID int64, applicationSlug string) (res *applicationEntity.ApplicationAccess, err error) {
	res = &applicationEntity.ApplicationAccess{}

	userApplication, err := c.UserApplicationRepository.GetUserApplicationByUserIDAndApplicationSlugDB(ctx, userID, applicationSlug)
	if err != nil {
		return nil, err
	}

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

func (c *ApplicationController) GetApplicationService(ctx context.Context, userID int64, applicationID int64, appServiceSlug string) (res *applicationEntity.ApplicationService, err error) {
	appService, err := c.UserApplicationServiceRepository.GetApplicationServiceAccessDB(ctx, userID, applicationID, appServiceSlug)
	if err != nil {
		return nil, err
	}

	if appService == nil {
		return nil, nil
	}

	res = &applicationEntity.ApplicationService{
		ID:    appService.AppServiceID,
		Slug:  appService.AppServiceSlug,
		Scope: enums.ApplicationScope(appService.AppServiceScope).ToString(),
		Name:  appService.AppServiceName,
	}

	return res, err
}
