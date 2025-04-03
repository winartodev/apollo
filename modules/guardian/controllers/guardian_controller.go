package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	applicationController "github.com/winartodev/apollo/modules/application/controllers"
)

type GuardianControllerItf interface {
	CheckUserPermission(ctx context.Context, userId int64, applicationSlug string, appServiceSlug string, httpMethod string) (permissionGranted bool, err error)
}

type GuardianController struct {
	ApplicationController applicationController.ApplicationControllerItf
}

func NewGuardianController(controller GuardianController) GuardianControllerItf {
	return &GuardianController{
		ApplicationController: controller.ApplicationController,
	}
}

func (c *GuardianController) CheckUserPermission(ctx context.Context, userID int64, appSlug string, appServiceSlug string, httpMethod string) (permissionGranted bool, err error) {
	applicationAccess, err := c.ApplicationController.GetApplicationAccess(ctx, userID, appSlug)
	if err != nil {
		return false, err
	}

	if applicationAccess.Applications == nil {
		return false, nil
	}

	application := &applicationAccess.Applications[0]
	appServices, err := c.ApplicationController.GetApplicationService(ctx, userID, application.ID, appServiceSlug)
	if err != nil {
		return false, err
	}

	if appServices == nil {
		return false, nil
	}

	application.Services = appServices

	marshaled, err := json.MarshalIndent(applicationAccess, "", "   ")
	if err != nil {
		return false, err
	}

	fmt.Println(string(marshaled))

	return true, nil
}
