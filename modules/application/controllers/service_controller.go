package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/winartodev/apollo/core/errors"
	"github.com/winartodev/apollo/core/helpers"
	"github.com/winartodev/apollo/core/responses"
	applicationEntity "github.com/winartodev/apollo/modules/application/entities"
	applicationRepo "github.com/winartodev/apollo/modules/application/repository"
	"time"
)

const (
	serviceSlugAlreadyExists      = "service with slug %s already exists"
	serviceNotFoundWithID         = "service not found with id %d"
	duplicateServiceInSameRequest = "duplicate service with slug %s"
)

type ServiceControllerItf interface {
	BulkInsertServiceTx(ctx context.Context, tx *sql.Tx, applicationID int64, services []applicationEntity.Service) (res *responses.MutateResponseData, err errors.Errors)
	GetServiceBySlug(ctx context.Context, applicationID int64, slug string) (res *applicationEntity.Service, err errors.Errors)
	GetServiceByID(ctx context.Context, applicationID int64, serviceID int64) (res *applicationEntity.Service, err errors.Errors)
	GetServices(ctx context.Context, applicationID int64, paginate helpers.Paginate) (res []applicationEntity.Service, total int64, err errors.Errors)
	UpdateServiceByID(ctx context.Context, applicationID int64, serviceID int64, data applicationEntity.Service) (res *responses.MutateResponseData, err errors.Errors)
	ActivateServiceByID(ctx context.Context, applicationID int64, serviceID int64, isActive bool) (res *responses.MutateResponseData, err errors.Errors)
	ProcessValidationService(ctx context.Context, userID int64, applicationID int64, services []applicationEntity.Service) (successRowsData []applicationEntity.Service, failRowsData []map[string]interface{}, err errors.Errors)
}

type ServiceController struct {
	Tx          helpers.DBTransactItf
	ServiceRepo applicationRepo.ServiceRepositoryItf
}

func NewServiceController(controller ServiceController) ServiceControllerItf {
	return &ServiceController{
		Tx:          controller.Tx,
		ServiceRepo: controller.ServiceRepo,
	}
}

func (sc *ServiceController) BulkInsertServiceTx(ctx context.Context, tx *sql.Tx, applicationID int64, services []applicationEntity.Service) (res *responses.MutateResponseData, err errors.Errors) {
	userID, ctxErr := helpers.GetUserIDFromContext(ctx)
	if ctxErr != nil {
		return nil, errors.AuthorizationErr.WithReason(ctxErr.Error())
	}

	successRowsData, failRowsData, err := sc.ProcessValidationService(ctx, userID, applicationID, services)
	if err != nil {
		return nil, err
	}

	res = &responses.MutateResponseData{}
	if failRowsData != nil && len(failRowsData) > 0 {
		res.FailedRowsData = failRowsData
		res.FailData = len(failRowsData)

		return res, nil
	}

	var sqlErr error
	shouldCommit := false
	if tx == nil {
		tx, sqlErr = sc.Tx.BeginTx(ctx)
		if sqlErr != nil {
			return nil, errors.InternalServerErr(sqlErr.Error())
		}

		shouldCommit = true
		defer func() {
			if shouldCommit && sqlErr != nil {
				sc.Tx.Rollback()
			}
		}()
	}

	ids, sqlErr := sc.ServiceRepo.BulkInsertServiceTxDB(ctx, tx, successRowsData)
	if sqlErr != nil {
		return nil, errors.InternalServerErr(sqlErr.Error())
	}

	for i, _ := range successRowsData {
		successRowsData[i].ID = ids[i]
	}

	if shouldCommit {
		if sqlErr = sc.Tx.Commit(); sqlErr != nil {
			return nil, errors.InternalServerErr(sqlErr.Error())
		}
	}

	if successRowsData != nil && len(successRowsData) > 0 {
		res.SuccessRowsData = successRowsData
		res.SuccessData = len(successRowsData)
	}

	return res, nil
}

func (sc *ServiceController) GetServiceBySlug(ctx context.Context, applicationID int64, slug string) (res *applicationEntity.Service, err errors.Errors) {
	data, dbErr := sc.ServiceRepo.GetServiceBySlugDB(ctx, applicationID, slug)
	if dbErr != nil {
		return nil, errors.InternalServerErr(dbErr.Error())
	}

	return data, nil
}

func (sc *ServiceController) GetServiceByID(ctx context.Context, applicationID int64, serviceID int64) (res *applicationEntity.Service, err errors.Errors) {
	data, dbErr := sc.ServiceRepo.GetServiceByIDDB(ctx, applicationID, serviceID)
	if dbErr != nil {
		return nil, errors.InternalServerErr(dbErr.Error())
	}

	return data, nil
}

func (sc *ServiceController) GetServices(ctx context.Context, applicationID int64, paginate helpers.Paginate) (res []applicationEntity.Service, total int64, err errors.Errors) {
	total, dbErr := sc.ServiceRepo.GetTotalServiceDB(ctx, applicationID)
	if dbErr != nil {
		return nil, 0, errors.InternalServerErr(dbErr.Error())
	}

	data, dbErr := sc.ServiceRepo.GetServicesDB(ctx, applicationID, paginate)
	if dbErr != nil {
		return nil, 0, errors.InternalServerErr(dbErr.Error())
	}

	return data, total, nil
}

func (sc *ServiceController) UpdateServiceByID(ctx context.Context, applicationID int64, serviceID int64, data applicationEntity.Service) (res *responses.MutateResponseData, err errors.Errors) {
	userID, ctxErr := helpers.GetUserIDFromContext(ctx)
	if ctxErr != nil {
		return nil, errors.InvalidUserID.WithReason(ctxErr.Error())
	}

	if userID == 0 {
		return nil, errors.InvalidUserID
	}

	oldData, respErr := sc.GetServiceByID(ctx, applicationID, serviceID)
	if respErr != nil {
		return nil, respErr
	}

	if oldData == nil {
		return nil, errors.DataNotFoundErr.WithReason(fmt.Sprintf(serviceNotFoundWithID, serviceID))
	}

	updateData := *oldData
	updateData.Name = data.Name
	updateData.Description = data.Description

	var services []applicationEntity.Service
	services = append(services, updateData)
	successRowsData, failRowsData, err := sc.ProcessValidationService(ctx, userID, applicationID, services)
	if err != nil {
		return nil, err
	}

	res = &responses.MutateResponseData{}
	if failRowsData != nil && len(failRowsData) > 0 {
		res.FailedRowsData = failRowsData
		res.FailData = len(failRowsData)

		return res, nil
	}

	if oldData.Slug == successRowsData[0].Slug {
		return res, nil
	}

	dbErr := sc.ServiceRepo.UpdateServiceByIDDB(ctx, applicationID, serviceID, successRowsData[0])
	if dbErr != nil {
		return nil, errors.InternalServerErr(dbErr.Error())
	}

	if successRowsData != nil && len(successRowsData) > 0 {
		res.SuccessRowsData = successRowsData
		res.SuccessData = len(successRowsData)
	}

	return res, nil
}

func (sc *ServiceController) ActivateServiceByID(ctx context.Context, applicationID int64, serviceID int64, isActive bool) (res *responses.MutateResponseData, err errors.Errors) {
	now := time.Now()
	userID, ctxErr := helpers.GetUserIDFromContext(ctx)
	if ctxErr != nil {
		return nil, errors.InvalidUserID.WithReason(ctxErr.Error())
	}

	oldData, respErr := sc.GetServiceByID(ctx, applicationID, serviceID)
	if respErr != nil {
		return nil, respErr
	}

	if oldData == nil {
		return nil, errors.DataNotFoundErr.WithReason(fmt.Sprintf(serviceNotFoundWithID, serviceID))
	}

	updateData := *oldData
	updateData.IsActive = isActive
	updateData.UpdatedBy = userID
	updateData.UpdatedAt = &now

	dbErr := sc.ServiceRepo.UpdateServiceByIDDB(ctx, applicationID, serviceID, updateData)
	if dbErr != nil {
		return nil, errors.InternalServerErr(dbErr.Error())
	}

	return &responses.MutateResponseData{
		SuccessData:     1,
		SuccessRowsData: updateData,
	}, nil
}

func (sc *ServiceController) ProcessValidationService(ctx context.Context, userID int64, applicationID int64, services []applicationEntity.Service) (successRowsData []applicationEntity.Service, failRowsData []map[string]interface{}, err errors.Errors) {
	now := time.Now()
	seenSlugs := make(map[string]bool)

	for i, _ := range services {
		var service = services[i]
		var serviceErrors []string

		service.ApplicationID = &applicationID

		if err := applicationEntity.GenerateServiceSlug(&service); err != nil {
			serviceErrors = append(serviceErrors, err.ToString())
		}

		if seenSlugs[service.Slug] {
			serviceErrors = append(serviceErrors, fmt.Sprintf(duplicateServiceInSameRequest, service.Slug))
			failRowsData = append(failRowsData, map[string]interface{}{
				"data":   service,
				"errors": serviceErrors,
			})
			continue
		}

		seenSlugs[service.Slug] = true

		if dataExists, err := sc.GetServiceBySlug(ctx, applicationID, service.Slug); err != nil {
			return nil, nil, err
		} else if dataExists != nil && dataExists.ID != service.ID {
			serviceErrors = append(serviceErrors, fmt.Sprintf(serviceSlugAlreadyExists, service.Slug))
		}

		if len(serviceErrors) > 0 {
			failRowsData = append(failRowsData, map[string]interface{}{
				"data":   service,
				"errors": serviceErrors,
			})
			continue
		}

		service.CreatedBy = userID
		service.UpdatedBy = userID
		service.CreatedAt = &now
		service.UpdatedAt = &now

		successRowsData = append(successRowsData, service)
	}

	return successRowsData, failRowsData, nil
}
