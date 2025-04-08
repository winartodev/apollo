package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	applicationController "github.com/winartodev/apollo/modules/application/controllers"
	applicationEnum "github.com/winartodev/apollo/modules/application/enums"
	guardianEntity "github.com/winartodev/apollo/modules/guardian/entities"
	guardianRepo "github.com/winartodev/apollo/modules/guardian/repositories"
	userController "github.com/winartodev/apollo/modules/user/controllers"
	"strings"
)

var (
	ErrorUserNotGetPermissionToTheService = errors.New("User not get permission to the service")
)

type GuardianControllerItf interface {
	CheckUserPermissionToInternalApp(ctx context.Context, userId int64, application applicationEnum.ApplicationEnum, applicationService applicationEnum.ApplicationServiceEnum, httpMethod string) (permissionGranted bool, err error)
	CheckApplicationServicePermission(ctx context.Context, userID int64, roleID int64, appID int64, appServiceSlug string) (res *guardianEntity.GuardianApplicationService, hasPermission bool, err error)
}

type GuardianController struct {
	UserController               userController.UserControllerItf
	ApplicationController        applicationController.ApplicationControllerItf
	GuardianPermissionRepository guardianRepo.GuardianPermissionRepositoryItf
}

func NewGuardianController(controller GuardianController) GuardianControllerItf {
	return &GuardianController{
		UserController:               controller.UserController,
		GuardianPermissionRepository: controller.GuardianPermissionRepository,
		ApplicationController:        controller.ApplicationController,
	}
}

func (c *GuardianController) CheckUserPermissionToInternalApp(ctx context.Context, userID int64, application applicationEnum.ApplicationEnum, applicationService applicationEnum.ApplicationServiceEnum, httpMethod string) (permissionGranted bool, err error) {
	// get user data is exists
	userData, err := c.UserController.GetUserByID(ctx, userID)
	if err != nil {
		return false, err
	}

	if userData == nil {
		return false, nil
	}

	// get user application access
	userApplication, err := c.UserController.GetUserApplicationByUserIDAndApplicationSlug(ctx, userID, application.ToSlug())
	if err != nil {
		return false, err
	}

	if userApplication == nil {
		return false, nil
	}

	// get user role
	userRoleData, err := c.UserController.GetUserRoleByID(ctx, userID, userApplication.ID)
	if err != nil {
		return false, err
	}

	// check service permission
	servicePermissions, hasServicePermission, err := c.CheckApplicationServicePermission(ctx,
		userID,
		userRoleData.RoleID,
		userApplication.ID,
		applicationService.ToSlug(),
	)
	if err != nil {
		return false, err
	}

	if !hasServicePermission {
		return false, ErrorUserNotGetPermissionToTheService
	}

	if !c.IsEligiblePermissionBaseOnHTTPMethod(servicePermissions.Permissions, strings.ToLower(httpMethod)) {
		return false, ErrorUserNotGetPermissionToTheService
	}

	var guardianAccessPermission guardianEntity.GuardianUserAccessPermission
	result := guardianAccessPermission.Build(userData, userRoleData, userApplication, servicePermissions)

	marshaled, err := json.MarshalIndent(result, "", "   ")
	if err != nil {
		return false, err
	}

	fmt.Println(string(marshaled))

	// save into redis
	redisKey := fmt.Sprintf("guardian_access_permission:%d:%s:%s", userID, application.ToSlug(), applicationService.ToSlug())
	fmt.Printf(redisKey)

	return true, nil
}

func (c *GuardianController) CheckApplicationServicePermission(ctx context.Context, userID int64, roleID int64, appID int64, appServiceSlug string) (res *guardianEntity.GuardianApplicationService, hasPermission bool, err error) {
	appService, err := c.ApplicationController.GetApplicationServiceBySlug(ctx, appServiceSlug)
	if err != nil {
		return nil, false, err
	}

	if appService == nil {
		return nil, false, nil
	}

	servicePermissions, err := c.GuardianPermissionRepository.GetApplicationServicePermissionsDB(ctx, userID, roleID, appID, appService.ID)
	if err != nil {
		return nil, false, err
	}

	if servicePermissions == nil {
		return nil, false, nil
	}

	res = &guardianEntity.GuardianApplicationService{
		ID:          appService.ID,
		Scope:       appService.Scope,
		Slug:        appService.Slug,
		Name:        appService.Name,
		Permissions: servicePermissions,
	}

	return res, len(servicePermissions) > 0, nil
}

func (c *GuardianController) IsEligiblePermissionBaseOnHTTPMethod(permissions []guardianEntity.GuardianPermission, httpMethod string) bool {
	for _, permission := range permissions {
		if permission.Slug == httpMethod {
			return true

		}
	}

	return false
}
