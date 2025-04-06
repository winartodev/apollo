package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	coreEnum "github.com/winartodev/apollo/core/enums"
	applicationController "github.com/winartodev/apollo/modules/application/controllers"
	guardianEntity "github.com/winartodev/apollo/modules/guardian/entities"
	userController "github.com/winartodev/apollo/modules/user/controllers"
)

type GuardianControllerItf interface {
	CheckUserPermissionToInternalApp(ctx context.Context, userId int64, application coreEnum.ApplicationEnum, applicationService coreEnum.ApplicationServiceEnum, httpMethod string) (permissionGranted bool, err error)
}

type GuardianController struct {
	UserController        userController.UserControllerItf
	ApplicationController applicationController.ApplicationControllerItf
}

func NewGuardianController(controller GuardianController) GuardianControllerItf {
	return &GuardianController{
		UserController:        controller.UserController,
		ApplicationController: controller.ApplicationController,
	}
}

func (c *GuardianController) CheckUserPermissionToInternalApp(ctx context.Context, userID int64, application coreEnum.ApplicationEnum, applicationService coreEnum.ApplicationServiceEnum, httpMethod string) (permissionGranted bool, err error) {
	// get user data is exists
	userData, err := c.UserController.GetUserByID(ctx, userID)
	if err != nil {
		return false, err
	}

	if userData == nil {
		return false, nil
	}

	// get user application access
	appSlug := application.ToSlug()
	applicationAccess, err := c.ApplicationController.GetApplicationAccess(ctx, userID, appSlug)
	if err != nil {
		return false, err
	}

	if applicationAccess.Applications == nil {
		return false, nil
	}

	applicationData := &applicationAccess.Applications[0]

	// get user application service access
	appServiceSlug := applicationService.ToSlug()
	appService, err := c.ApplicationController.GetApplicationService(ctx, userID, applicationData.ID, appServiceSlug)
	if err != nil {
		return false, err
	}

	if appService == nil {
		return false, nil
	}

	// get user role
	userRoleData, err := c.UserController.GetUserRoleByID(ctx, userID)
	if err != nil {
		return false, err
	}

	var guardianAccessPermission guardianEntity.GuardianUserAccessPermission
	result := guardianAccessPermission.Build(userData, userRoleData, applicationData, appService)

	marshaled, err := json.MarshalIndent(result, "", "   ")
	if err != nil {
		return false, err
	}

	fmt.Println(string(marshaled))

	return true, nil
}
