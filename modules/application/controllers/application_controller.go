package controllers

import (
	"context"
	"errors"
	applicationEntity "github.com/winartodev/apollo/modules/application/entities"
	coreEnum "github.com/winartodev/apollo/modules/application/enums"
	applicationRepo "github.com/winartodev/apollo/modules/application/repositories"
)

var (
	ErrorApplicationServiceInactive = errors.New("application service inactive")
)

type ApplicationControllerItf interface {
	GetApplicationBySlug(ctx context.Context, slug string) (res *applicationEntity.Application, err error)
	GetApplicationServiceBySlug(ctx context.Context, slug string) (res *applicationEntity.ApplicationService, err error)
}

type ApplicationController struct {
	ApplicationServiceRepository applicationRepo.ApplicationServiceRepositoryItf
}

func NewApplicationController(controller ApplicationController) ApplicationControllerItf {
	return &ApplicationController{
		ApplicationServiceRepository: controller.ApplicationServiceRepository,
	}
}

func (ac *ApplicationController) GetApplicationBySlug(ctx context.Context, slug string) (res *applicationEntity.Application, err error) {
	return res, err
}

func (ac *ApplicationController) GetApplicationServiceBySlug(ctx context.Context, slug string) (res *applicationEntity.ApplicationService, err error) {
	res, err = ac.ApplicationServiceRepository.GetApplicationServiceBySlugDB(ctx, slug)
	if err != nil {
		return nil, err
	}

	if !res.IsActive {
		return nil, ErrorApplicationServiceInactive
	}

	res.Scope = coreEnum.ApplicationServiceEnum(res.Scope).ToSlug()

	return res, nil
}
