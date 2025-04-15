package controllers

import (
	"context"
	"fmt"
	"github.com/winartodev/apollo/core/errors"
	"github.com/winartodev/apollo/core/helpers"
	applicaitonEntity "github.com/winartodev/apollo/modules/application/entities"
	applicationRepo "github.com/winartodev/apollo/modules/application/repository"
	"time"
)

const (
	slugAlreadyExists     = "slug %s already exists"
	serviceNotFoundWithID = "service not found with id %d"
)

type ServiceControllerItf interface {
	CreateService(ctx context.Context, data applicaitonEntity.Service) (res *applicaitonEntity.Service, err errors.Errors)
	GetServiceBySlug(ctx context.Context, slug string) (res *applicaitonEntity.Service, err errors.Errors)
	GetServiceByID(ctx context.Context, id int64) (res *applicaitonEntity.Service, err errors.Errors)
	GetServices(ctx context.Context, paginate helpers.Paginate) (res []applicaitonEntity.Service, total int64, err errors.Errors)
	UpdateServiceByID(ctx context.Context, id int64, data applicaitonEntity.Service) (res *applicaitonEntity.Service, err errors.Errors)
	ActivateServiceByID(ctx context.Context, id int64, isActive bool) (res *applicaitonEntity.Service, err errors.Errors)
}

type ServiceController struct {
	ServiceRepo applicationRepo.ServiceRepositoryItf
}

func NewServiceController(controller ServiceController) ServiceControllerItf {
	return &ServiceController{
		ServiceRepo: controller.ServiceRepo,
	}
}

func (sc *ServiceController) CreateService(ctx context.Context, data applicaitonEntity.Service) (res *applicaitonEntity.Service, err errors.Errors) {
	err = applicaitonEntity.GenerateServiceSlug(&data)
	if err != nil {
		return nil, err
	}

	dataExists, err := sc.GetServiceBySlug(ctx, data.Slug)
	if err != nil {
		return nil, err
	}

	if dataExists != nil {
		return nil, errors.DataAlreadyExistsErr(fmt.Sprintf(slugAlreadyExists, data.Slug))
	}

	now := time.Now()
	data.CreatedAt = &now
	data.UpdatedAt = &now

	id, dbErr := sc.ServiceRepo.CreateServiceDB(ctx, data)
	if dbErr != nil {
		return nil, errors.InternalServerErr(dbErr.Error())
	}

	data.ID = *id

	return &data, nil
}

func (sc *ServiceController) GetServiceBySlug(ctx context.Context, slug string) (res *applicaitonEntity.Service, err errors.Errors) {
	data, dbErr := sc.ServiceRepo.GetServiceBySlugDB(ctx, slug)
	if dbErr != nil {
		return nil, errors.InternalServerErr(dbErr.Error())
	}

	return data, nil
}

func (sc *ServiceController) GetServiceByID(ctx context.Context, id int64) (res *applicaitonEntity.Service, err errors.Errors) {
	data, dbErr := sc.ServiceRepo.GetServiceByIDDB(ctx, id)
	if dbErr != nil {
		return nil, errors.InternalServerErr(dbErr.Error())
	}

	return data, nil
}

func (sc *ServiceController) GetServices(ctx context.Context, paginate helpers.Paginate) (res []applicaitonEntity.Service, total int64, err errors.Errors) {
	total, dbErr := sc.ServiceRepo.GetTotalServiceDB(ctx)
	if dbErr != nil {
		return nil, 0, errors.InternalServerErr(dbErr.Error())
	}

	data, dbErr := sc.ServiceRepo.GetServicesDB(ctx, paginate)
	if dbErr != nil {
		return nil, 0, errors.InternalServerErr(dbErr.Error())
	}

	return data, total, nil
}

func (sc *ServiceController) UpdateServiceByID(ctx context.Context, id int64, data applicaitonEntity.Service) (res *applicaitonEntity.Service, err errors.Errors) {
	now := time.Now()
	userID, ctxErr := helpers.GetUserIDFromContext(ctx)
	if ctxErr != nil {
		return nil, errors.InvalidUserID.WithReason(ctxErr.Error())
	}

	if userID == 0 {
		return nil, errors.InvalidUserID
	}

	oldData, respErr := sc.GetServiceByID(ctx, id)
	if respErr != nil {
		return nil, respErr
	}

	if oldData == nil {
		return nil, errors.DataNotFoundErr.WithReason(fmt.Sprintf(serviceNotFoundWithID, id))
	}

	updateData := *oldData

	err = applicaitonEntity.GenerateServiceSlug(&data)
	if err != nil {
		return nil, err
	}

	dataExists, err := sc.GetServiceBySlug(ctx, data.Slug)
	if err != nil {
		return nil, err
	}

	if dataExists != nil && dataExists.ID != id {
		return nil, errors.DataAlreadyExistsErr(fmt.Sprintf(slugAlreadyExists, data.Slug))
	}

	updateData.Name = data.Name
	updateData.Slug = data.Slug
	updateData.Description = data.Description
	updateData.UpdatedBy = userID
	updateData.UpdatedAt = &now

	dbErr := sc.ServiceRepo.UpdateServiceByIDDB(ctx, id, updateData)
	if dbErr != nil {
		return nil, errors.InternalServerErr(dbErr.Error())
	}

	return &updateData, nil
}

func (sc *ServiceController) ActivateServiceByID(ctx context.Context, id int64, isActive bool) (res *applicaitonEntity.Service, err errors.Errors) {
	now := time.Now()
	userID, ctxErr := helpers.GetUserIDFromContext(ctx)
	if ctxErr != nil {
		return nil, errors.InvalidUserID.WithReason(ctxErr.Error())
	}

	oldData, respErr := sc.GetServiceByID(ctx, id)
	if respErr != nil {
		return nil, respErr
	}

	if oldData == nil {
		return nil, errors.DataNotFoundErr.WithReason(fmt.Sprintf(serviceNotFoundWithID, id))
	}

	updateData := *oldData
	updateData.IsActive = isActive
	updateData.UpdatedBy = userID
	updateData.UpdatedAt = &now

	dbErr := sc.ServiceRepo.UpdateServiceByIDDB(ctx, id, updateData)
	if dbErr != nil {
		return nil, errors.InternalServerErr(dbErr.Error())
	}

	return &updateData, nil
}
