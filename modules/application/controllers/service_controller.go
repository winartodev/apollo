package controllers

import (
	"context"
	"fmt"
	applicaitonEntity "github.com/winartodev/apollo/modules/application/entities"
	applicationRepo "github.com/winartodev/apollo/modules/application/repository"
	"time"
)

type ServiceControllerItf interface {
	CreateService(ctx context.Context, data applicaitonEntity.Service) (res *applicaitonEntity.Service, affectedRows int64, err error)
	GetServiceBySlug(ctx context.Context, slug string) (res *applicaitonEntity.Service, err error)
}

type ServiceController struct {
	ServiceRepo applicationRepo.ServiceRepositoryItf
}

func NewServiceController(controller ServiceController) ServiceControllerItf {
	return &ServiceController{
		ServiceRepo: controller.ServiceRepo,
	}
}

func (sc *ServiceController) CreateService(ctx context.Context, data applicaitonEntity.Service) (res *applicaitonEntity.Service, affectedRows int64, err error) {
	err = applicaitonEntity.GenerateServiceSlug(&data)
	if err != nil {
		return nil, 0, err
	}

	dataExists, err := sc.GetServiceBySlug(ctx, data.Slug)
	if err != nil {
		return nil, 0, err
	}

	if dataExists != nil {
		return nil, 0, fmt.Errorf("slug %s already exists", data.Slug)
	}

	now := time.Now()
	data.CreatedAt = &now
	data.UpdatedAt = &now

	id, err := sc.ServiceRepo.CreateServiceDB(ctx, data)
	if err != nil {
		return nil, 0, err
	}

	data.ID = *id

	return &data, 1, nil
}

func (sc *ServiceController) GetServiceBySlug(ctx context.Context, slug string) (res *applicaitonEntity.Service, err error) {
	return sc.ServiceRepo.GetServiceBySlugDB(ctx, slug)
}
