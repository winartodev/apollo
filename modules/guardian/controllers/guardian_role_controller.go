package controllers

import (
	"context"
	"fmt"
	guardianRole "github.com/winartodev/apollo/modules/guardian/entities"
	guardianRepo "github.com/winartodev/apollo/modules/guardian/repositories"
	"time"
)

const (
	errorRoleAlreadyExists = "role with name %s in application id %d already exists"
)

type GuardianRoleControllerItf interface {
	CreateRole(ctx context.Context, data guardianRole.GuardianRole) (res *guardianRole.GuardianRole, err error)
	GetRoles(ctx context.Context, appID int64) (res []guardianRole.GuardianRole, err error)
	GetRoleBySlug(ctx context.Context, appID int64, slug string) (res *guardianRole.GuardianRole, err error)
	GetRoleByID(ctx context.Context, appID int64, id int64) (res *guardianRole.GuardianRole, err error)
}

type GuardianRoleController struct {
	GuardianRoleRepo guardianRepo.GuardianRoleRepositoryItf
}

func NewGuardianRoleController(controller GuardianRoleController) GuardianRoleControllerItf {
	return &GuardianRoleController{
		GuardianRoleRepo: controller.GuardianRoleRepo,
	}
}
func (grc *GuardianRoleController) CreateRole(ctx context.Context, data guardianRole.GuardianRole) (res *guardianRole.GuardianRole, err error) {
	now := time.Now()

	data.Slug = data.GenerateSlug()
	data.CreatedAt = &now
	data.UpdatedAt = &now

	existsData, err := grc.GetRoleBySlug(ctx, data.ApplicationID, data.Slug)
	if err != nil {
		return nil, err
	}

	if existsData != nil {
		return nil, fmt.Errorf(errorRoleAlreadyExists, data.Name, data.ApplicationID)
	}

	id, err := grc.GuardianRoleRepo.CreateRole(ctx, data)
	if err != nil {
		return nil, err
	}

	data.ID = id

	return &data, nil
}

func (grc *GuardianRoleController) GetRoles(ctx context.Context, appID int64) (res []guardianRole.GuardianRole, err error) {
	return grc.GuardianRoleRepo.GetRoles(ctx, appID, nil)
}

func (grc *GuardianRoleController) GetRoleBySlug(ctx context.Context, appID int64, slug string) (res *guardianRole.GuardianRole, err error) {
	return grc.GuardianRoleRepo.GetRoleBySlug(ctx, appID, slug)
}

func (grc *GuardianRoleController) GetRoleByID(ctx context.Context, appID int64, id int64) (res *guardianRole.GuardianRole, err error) {
	return grc.GuardianRoleRepo.GetRoleByID(ctx, appID, id)
}
