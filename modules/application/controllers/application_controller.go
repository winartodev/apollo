package controllers

import (
	"context"
	"github.com/winartodev/apollo/core/errors"
	"github.com/winartodev/apollo/core/helpers"
	applicationEntity "github.com/winartodev/apollo/modules/application/entities"
	applicationRepo "github.com/winartodev/apollo/modules/application/repository"
)

type ApplicationControllerItf interface {
	GetApplications(ctx context.Context, filter helpers.Paginate) (res []applicationEntity.Application, total int64, err errors.Errors)
	GetApplicationByID(ctx context.Context, id int64) (res *applicationEntity.Application, err errors.Errors)
}

type ApplicationController struct {
	ApplicationRepo applicationRepo.ApplicationRepositoryItf
}

func NewApplicationController(controller ApplicationController) ApplicationControllerItf {
	return &ApplicationController{
		ApplicationRepo: controller.ApplicationRepo,
	}
}

func (ac *ApplicationController) GetApplications(ctx context.Context, filter helpers.Paginate) (res []applicationEntity.Application, total int64, err errors.Errors) {
	total, dbErr := ac.ApplicationRepo.GetTotalApplicationDB(ctx)
	if dbErr != nil {
		return nil, 0, errors.InternalServerErr(dbErr.Error())
	}

	res, dbErr = ac.ApplicationRepo.GetApplicationsDB(ctx, filter)
	if dbErr != nil {
		return nil, 0, errors.InternalServerErr(dbErr.Error())
	}

	return res, total, err
}

func (ac *ApplicationController) GetApplicationByID(ctx context.Context, id int64) (res *applicationEntity.Application, err errors.Errors) {
	res, dbErr := ac.ApplicationRepo.GetApplicationByIDDB(ctx, id)
	if dbErr != nil {
		return nil, errors.InternalServerErr(dbErr.Error())
	}

	if res == nil {
		return nil, errors.DataNotFoundErr
	}

	return res, err
}
