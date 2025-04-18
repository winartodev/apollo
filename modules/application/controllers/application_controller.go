package controllers

import (
	"context"
	"fmt"
	"github.com/winartodev/apollo/core/errors"
	"github.com/winartodev/apollo/core/helpers"
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
	CreateApplication(ctx context.Context, data applicationEntity.Application) (res *applicationEntity.ApplicationResponse, err errors.Errors)
	GetApplications(ctx context.Context, filter helpers.Paginate) (res []applicationEntity.Application, total int64, err errors.Errors)
	GetApplicationByID(ctx context.Context, id int64) (res *applicationEntity.Application, err errors.Errors)
	GetApplicationBySlug(ctx context.Context, slug string) (res *applicationEntity.Application, err errors.Errors)
}

type ApplicationController struct {
	Tx                   helpers.DBTransactItf
	ApplicationRepo      applicationRepo.ApplicationRepositoryItf
	ApplicationScopeRepo applicationRepo.ApplicationScopeRepositoryItf
	UserApplicationRepo  userRepo.UserApplicationRepositoryItf
}

func NewApplicationController(controller ApplicationController) ApplicationControllerItf {
	return &ApplicationController{
		Tx:                   controller.Tx,
		ApplicationRepo:      controller.ApplicationRepo,
		ApplicationScopeRepo: controller.ApplicationScopeRepo,
		UserApplicationRepo:  controller.UserApplicationRepo,
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

func (ac *ApplicationController) CreateApplication(ctx context.Context, data applicationEntity.Application) (res *applicationEntity.ApplicationResponse, err errors.Errors) {
	userID, ctxErr := helpers.GetUserIDFromContext(ctx)
	if ctxErr != nil {
		return nil, errors.AuthorizationErr.WithReason(ctxErr.Error())
	}

	err = applicationEntity.GenerateApplicationSlug(&data)
	if err != nil {
		return nil, err
	}

	dataExists, err := ac.GetApplicationBySlug(ctx, data.Slug)
	if err != nil && err != errors.DataNotFoundErr {
		return nil, err
	}

	if dataExists != nil {
		return nil, errors.DataAlreadyExistsErr(fmt.Sprintf(applicationSlugAlreadyExists, data.Slug))
	}

	now := time.Now()
	data.CreatedBy = userID
	data.UpdatedBy = userID
	data.CreatedAt = &now
	data.UpdatedAt = &now

	// TODO: USE TRANSACTION WHEN CREATING NEW APPLICATIONS
	tx, dbErr := ac.Tx.BeginTx(ctx)
	if dbErr != nil {
		return nil, errors.InternalServerErr(dbErr.Error())
	}

	defer func() {
		if dbErr != nil {
			ac.Tx.Rollback()
		}
	}()

	var applicationID *int64
	applicationID, dbErr = ac.ApplicationRepo.CreateApplicationTxDB(ctx, tx, data)
	if dbErr != nil {
		return nil, errors.InternalServerErr(dbErr.Error())
	}

	data.ID = *applicationID
	if data.Scopes != nil && len(data.Scopes) > 0 {
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
		lastScopeID, dbErr = ac.ApplicationScopeRepo.BulkCreateApplicationScopeTxDB(ctx, tx, applicationScopes)
		if dbErr != nil {
			return nil, errors.InternalServerErr(dbErr.Error())
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

			_, dbErr = ac.UserApplicationRepo.BulkInsertUserApplicationTxDB(ctx, tx, userApplications)
			if dbErr != nil {
				return nil, errors.InternalServerErr(dbErr.Error())
			}
		}
	}

	if dbErr = ac.Tx.Commit(); dbErr != nil {
		return nil, errors.InternalServerErr(dbErr.Error())
	}

	res = data.ToResponse()
	return res, nil
}
