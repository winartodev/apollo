package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core"
	"github.com/winartodev/apollo/core/errors"
	"github.com/winartodev/apollo/core/helpers"
	"github.com/winartodev/apollo/core/middlewares"
	"github.com/winartodev/apollo/core/responses"
	"github.com/winartodev/apollo/modules/application"
	applicationController "github.com/winartodev/apollo/modules/application/controllers"
	applicationEntity "github.com/winartodev/apollo/modules/application/entities"
	"reflect"
)

const (
	applicationsEndpoint = "%s%s/apollo/applications"
)

type ApplicationHandler struct {
	middlewares.Middleware
	ApplicationController applicationController.ApplicationControllerItf
	ServiceController     applicationController.ServiceControllerItf
}

func NewApplicationHandler(handler ApplicationHandler) ApplicationHandler {
	return ApplicationHandler{
		Middleware:            handler.Middleware,
		ApplicationController: handler.ApplicationController,
		ServiceController:     handler.ServiceController,
	}
}

func (ah *ApplicationHandler) GetApplications(ctx *fiber.Ctx) error {
	_, err := helpers.GetUserIDFromFiberContext(ctx)
	if err != nil {
		return responses.FailedResponseWithError(ctx, errors.AuthorizationErr.WithReason(err.Error()))
	}

	access, ctxErr := helpers.GetApplicationAccessFromFiberContext(ctx)
	if ctxErr != nil {
		return responses.FailedResponseWithError(ctx, errors.AuthorizationErr.WithReason(ctxErr.Error()))
	}

	if access == nil || !reflect.DeepEqual(access, application.ApolloInternal) {
		return responses.FailedResponseWithError(ctx, errors.UserApplicationHasNotAccess)
	}

	paginate := helpers.Paginate{
		ValidOrderOptions: applicationEntity.AllowedApplicationFields,
	}

	paginateData, errResp := paginate.NewFromContext(ctx)
	if errResp != nil {
		return responses.FailedResponseWithError(ctx, errResp)
	}

	context := ctx.Context()
	resp, total, errResp := ah.ApplicationController.GetApplications(context, *paginateData)
	if errResp != nil {
		return responses.FailedResponseWithError(ctx, errResp)
	}

	return responses.SuccessResponseWithPaginate(ctx, fiber.StatusOK, resp, total, paginateData, nil)
}

func (ah *ApplicationHandler) CreateApplication(ctx *fiber.Ctx) error {
	_, err := helpers.GetUserIDFromFiberContext(ctx)
	if err != nil {
		return responses.FailedResponseWithError(ctx, errors.AuthorizationErr.WithReason(err.Error()))
	}

	access, ctxErr := helpers.GetApplicationAccessFromFiberContext(ctx)
	if ctxErr != nil {
		return responses.FailedResponseWithError(ctx, errors.AuthorizationErr.WithReason(ctxErr.Error()))
	}

	if access == nil || !reflect.DeepEqual(access, application.ApolloInternal) {
		return responses.FailedResponseWithError(ctx, errors.UserApplicationHasNotAccess)
	}

	data := applicationEntity.Application{}
	err = ctx.BodyParser(&data)
	if err != nil {
		return responses.FailedResponseWithError(ctx, errors.FailedParseRequestBodyErr.WithReason(err.Error()))
	}

	context := ctx.Context()
	resp, errResp := ah.ApplicationController.CreateApplication(context, data)
	if errResp != nil {
		return responses.FailedResponseWithError(ctx, errResp)
	}

	if resp.FailedRowsData != nil {
		return responses.FailedResponseV2(ctx, fiber.StatusBadRequest, resp)
	}

	return responses.SuccessResponseV2(ctx, fiber.StatusCreated, resp, nil)
}

func (ah *ApplicationHandler) GetApplication(ctx *fiber.Ctx) error {
	_, err := helpers.GetUserIDFromFiberContext(ctx)
	if err != nil {
		return responses.FailedResponseWithError(ctx, errors.AuthorizationErr.WithReason(err.Error()))
	}

	access, ctxErr := helpers.GetApplicationAccessFromFiberContext(ctx)
	if ctxErr != nil {
		return responses.FailedResponseWithError(ctx, errors.AuthorizationErr.WithReason(ctxErr.Error()))
	}

	if access == nil || !reflect.DeepEqual(access, application.ApolloInternal) {
		return responses.FailedResponseWithError(ctx, errors.UserApplicationHasNotAccess)
	}

	applicationID, err := ctx.ParamsInt("application_id", 0)
	if err != nil {
		return responses.FailedResponseWithError(ctx, errors.BadRequestErr.WithReason(err.Error()))
	}

	context := ctx.Context()
	applicationData, errResp := ah.ApplicationController.GetApplicationByID(context, int64(applicationID))
	if errResp != nil {
		return responses.FailedResponseWithError(ctx, errResp)
	}

	if applicationData == nil {
		return responses.FailedResponseWithError(ctx, errors.DataNotFoundErr)
	}

	return responses.SuccessResponseV2(ctx, fiber.StatusOK, applicationData, nil)
}

func (ah *ApplicationHandler) Current(ctx *fiber.Ctx) error {
	_, err := helpers.GetUserIDFromFiberContext(ctx)
	if err != nil {
		return responses.FailedResponseWithError(ctx, errors.AuthorizationErr.WithReason(err.Error()))
	}

	access, ctxErr := helpers.GetApplicationAccessFromFiberContext(ctx)
	if ctxErr != nil {
		return responses.FailedResponseWithError(ctx, errors.AuthorizationErr.WithReason(ctxErr.Error()))
	}

	if access == nil || !reflect.DeepEqual(access, application.ApolloInternal) {
		return responses.FailedResponseWithError(ctx, errors.UserApplicationHasNotAccess)
	}

	paginate := helpers.Paginate{
		ValidOrderOptions: applicationEntity.AllowedApplicationFields,
	}

	paginateData := paginate.BuildDefault()

	context := ctx.Context()
	applicationData, errResp := ah.ApplicationController.GetApplicationByID(context, access.ID.ToInt64())
	if errResp != nil {
		return responses.FailedResponseWithError(ctx, errResp)
	}

	if applicationData == nil {
		return responses.FailedResponseWithError(ctx, errors.DataNotFoundErr)
	}

	servicesData, total, errResp := ah.ServiceController.GetServices(context, applicationData.ID, *paginateData)
	if errResp != nil {
		return responses.FailedResponseWithError(ctx, errResp)
	}

	res := applicationEntity.ApplicationService{
		Application: *applicationData,
		Services:    servicesData,
	}

	return responses.SuccessResponseWithPaginate(ctx, fiber.StatusOK, res, total, paginateData, nil)
}

func (ah *ApplicationHandler) Register(router fiber.Router) error {
	internal := router.Group(fmt.Sprintf(applicationsEndpoint, core.V1, core.AccessInternal), ah.HandleInternalAccess())
	internal.Post("/", ah.CreateApplication)
	internal.Get("/", ah.GetApplications)
	internal.Get("/current", ah.Current)
	internal.Get("/:application_id", ah.GetApplication)

	return nil
}
