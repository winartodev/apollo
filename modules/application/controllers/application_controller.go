package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/winartodev/apollo/core/errors"
	"github.com/winartodev/apollo/core/helpers"
	"github.com/winartodev/apollo/core/responses"
	"github.com/winartodev/apollo/modules/application"
	applicationEntity "github.com/winartodev/apollo/modules/application/entities"
	applicationRepo "github.com/winartodev/apollo/modules/application/repository"
	userEntity "github.com/winartodev/apollo/modules/user/entities"
	userRepo "github.com/winartodev/apollo/modules/user/repositories"
	"time"
)

const (
	applicationSlugAlreadyExists = "application with slug %s already exists"
	applicationNotFoundWithID    = "application not found with id %d"
)

type ApplicationControllerItf interface {
	CreateApplication(ctx context.Context, data applicationEntity.Application) (res *responses.MutateResponseData, err errors.Errors)
	GetApplications(ctx context.Context, filter helpers.Paginate) (res []applicationEntity.Application, total int64, err errors.Errors)
	GetApplicationByID(ctx context.Context, id int64) (res *applicationEntity.Application, err errors.Errors)
	GetApplicationBySlug(ctx context.Context, slug string) (res *applicationEntity.Application, err errors.Errors)
}

type ApplicationController struct {
	Tx                   helpers.DBTransactItf
	ApplicationRepo      applicationRepo.ApplicationRepositoryItf
	ApplicationScopeRepo applicationRepo.ApplicationScopeRepositoryItf
	UserApplicationRepo  userRepo.UserApplicationRepositoryItf
	ServiceController    ServiceControllerItf
}

func NewApplicationController(controller ApplicationController) ApplicationControllerItf {
	return &ApplicationController{
		Tx:                   controller.Tx,
		ApplicationRepo:      controller.ApplicationRepo,
		ApplicationScopeRepo: controller.ApplicationScopeRepo,
		UserApplicationRepo:  controller.UserApplicationRepo,
		ServiceController:    controller.ServiceController,
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
		return nil, errors.DataNotFoundErr.WithReason(fmt.Sprintf(applicationNotFoundWithID, id))
	}

	return res, err
}

func (ac *ApplicationController) GetApplicationBySlug(ctx context.Context, slug string) (res *applicationEntity.Application, err errors.Errors) {
	if err = helpers.IsValidSlug(slug); err != nil {
		return nil, err
	}

	res, dbErr := ac.ApplicationRepo.GetApplicationBySlugDB(ctx, slug)
	if dbErr != nil {
		return nil, errors.InternalServerErr(dbErr.Error())
	}

	if res == nil {
		return nil, errors.DataNotFoundErr
	}

	return res, nil
}

func (ac *ApplicationController) CreateApplication(ctx context.Context, data applicationEntity.Application) (res *responses.MutateResponseData, errs errors.Errors) {
	res = &responses.MutateResponseData{}
	var err error

	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, errors.AuthorizationErr.WithReason(err.Error())
	}

	errs = applicationEntity.GenerateApplicationSlug(&data)
	if errs != nil {
		return nil, errs
	}

	dataExists, errs := ac.GetApplicationBySlug(ctx, data.Slug)
	if errs != nil && errs != errors.DataNotFoundErr {
		return nil, errs
	}

	if dataExists != nil {
		return nil, errors.DataAlreadyExistsErr(fmt.Sprintf(applicationSlugAlreadyExists, data.Slug))
	}

	now := time.Now()
	data.CreatedBy = userID
	data.UpdatedBy = userID
	data.CreatedAt = &now
	data.UpdatedAt = &now

	tx, err := ac.Tx.BeginTx(ctx)
	if err != nil {
		return nil, errors.InternalServerErr(err.Error())
	}

	defer func() {
		if err != nil {
			ac.Tx.Rollback()
		}
	}()

	var applicationID *int64
	applicationID, err = ac.ApplicationRepo.CreateApplicationTxDB(ctx, tx, data)
	if err != nil {
		return nil, errors.InternalServerErr(err.Error())
	}

	data.ID = *applicationID
	if data.Scopes == nil && len(data.Scopes) <= 0 {
		data.Scopes = []application.Scope{application.Internal}
	}

	var applicationScopes []applicationEntity.ApplicationScope
	for _, scope := range data.Scopes {
		applicationScopes = append(applicationScopes, applicationEntity.ApplicationScope{
			ApplicationID: *applicationID,
			ScopeID:       scope.ToInt64(),
			IsActive:      true,
			CreatedBy:     userID,
			UpdatedBy:     userID,
			CreatedAt:     &now,
			UpdatedAt:     &now,
		})
	}

	var lastScopeID *int64
	lastScopeID, err = ac.ApplicationScopeRepo.BulkCreateApplicationScopeTxDB(ctx, tx, applicationScopes)
	if err != nil {
		return nil, errors.InternalServerErr(err.Error())
	}

	if lastScopeID != nil {
		var applicationScopeIDs []int64
		for i := *lastScopeID; i < (*lastScopeID + int64(len(data.Scopes))); i++ {
			applicationScopeIDs = append(applicationScopeIDs, i)
		}

		var userApplications []userEntity.UserApplication
		for _, applicationScopeID := range applicationScopeIDs {
			userApplications = append(userApplications, userEntity.UserApplication{
				UserID:             userID,
				ApplicationScopeID: applicationScopeID,
			})
		}

		_, err = ac.UserApplicationRepo.BulkInsertUserApplicationTxDB(ctx, tx, userApplications)
		if err != nil {
			return nil, errors.InternalServerErr(err.Error())
		}
	}

	if data.Services != nil {
		var serviceData *responses.MutateResponseData
		serviceData, errs = ac.ServiceController.BulkInsertServiceTx(ctx, tx, *applicationID, data.Services)
		if errs != nil {
			err = errs.ToError()
			return nil, errors.InternalServerErr(err.Error())
		}

		if serviceData.SuccessRowsData != nil {
			var serviceDataByte []byte
			serviceDataByte, err = json.Marshal(serviceData.SuccessRowsData)
			if err != nil {
				return nil, errors.InternalServerErr(err.Error())
			}

			var services []applicationEntity.Service
			err = json.Unmarshal(serviceDataByte, &services)
			if err != nil {
				return nil, errors.InternalServerErr(err.Error())
			}

			data.Services = services
		}

		if serviceData.FailedRowsData != nil {
			res.FailData = serviceData.FailData
			res.FailedRowsData = serviceData.FailedRowsData

			err = errors.BadRequestErr.ToError()

			return res, nil
		}
	}

	if err = ac.Tx.Commit(); err != nil {
		return nil, errors.InternalServerErr(err.Error())
	}

	res.SuccessData = 1
	res.SuccessRowsData = data.ToResponse()

	return res, nil
}
