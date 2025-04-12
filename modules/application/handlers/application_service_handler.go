package handlers

import (
	"encoding/json"
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
	applicationServiceEndpoint = applicationsEndpoint + "/:application_id/services"
)

type ApplicationServiceHandler struct {
	middlewares.Middleware
	ApplicationController applicationController.ApplicationControllerItf
	ServiceController     applicationController.ServiceControllerItf
}

func NewApplicationServiceHandler(handler ApplicationServiceHandler) ApplicationServiceHandler {
	return ApplicationServiceHandler{
		Middleware:            handler.Middleware,
		ApplicationController: handler.ApplicationController,
		ServiceController:     handler.ServiceController,
	}
}

func (ash *ApplicationServiceHandler) getParamApplicationID(ctx *fiber.Ctx) (res int64, err error) {
	applicationID, ctxErr := ctx.ParamsInt("application_id", 0)
	if ctxErr != nil {
		return 0, ctxErr
	}

	if applicationID <= 0 {
		return 0, fmt.Errorf(errors.ReasonInvalidParamID, "application_id")
	}

	return int64(applicationID), nil
}

func (ash *ApplicationServiceHandler) getParamServiceID(ctx *fiber.Ctx) (res int64, err error) {
	serviceID, ctxErr := ctx.ParamsInt("service_id", 0)
	if ctxErr != nil {
		return 0, ctxErr
	}

	if serviceID <= 0 {
		return 0, fmt.Errorf(errors.ReasonInvalidParamID, "service_id")
	}

	return int64(serviceID), nil
}

func (ash *ApplicationServiceHandler) validateUserAccess(ctx *fiber.Ctx) (err errors.Errors) {
	_, ctxErr := helpers.GetUserIDFromFiberContext(ctx)
	if ctxErr != nil {
		return errors.AuthorizationErr.WithReason(ctxErr.Error())
	}

	access, ctxErr := helpers.GetApplicationAccessFromFiberContext(ctx)
	if ctxErr != nil {
		return errors.AuthorizationErr.WithReason(ctxErr.Error())
	}

	if access == nil || !reflect.DeepEqual(access, application.ApolloInternal) {
		return errors.UserApplicationHasNotAccess
	}

	_, ctxErr = helpers.GetUserIDFromFiberContext(ctx)
	if ctxErr != nil {
		return errors.AuthorizationErr.WithReason(ctxErr.Error())
	}

	return nil
}

func (ash *ApplicationServiceHandler) CreateApplicationService(ctx *fiber.Ctx) error {
	if err := ash.validateUserAccess(ctx); err != nil {
		return responses.FailedResponseWithError(ctx, err)
	}

	applicationID, err := ash.getParamApplicationID(ctx)
	if err != nil {
		return responses.FailedResponseWithError(ctx, errors.BadRequestErr.WithReason(err.Error()))
	}

	var data = make(map[string]interface{})
	ctxErr := ctx.BodyParser(&data)
	if ctxErr != nil {
		return responses.FailedResponseWithError(ctx, errors.FailedParseRequestBodyErr.WithReason(ctxErr.Error()))
	}

	var services []applicationEntity.Service
	dataByte, err := json.Marshal(data["services"])
	if err != nil {
		return responses.FailedResponseWithError(ctx, errors.BadRequestErr.WithReason(err.Error()))
	}

	err = json.Unmarshal(dataByte, &services)
	if err != nil {
		return responses.FailedResponseWithError(ctx, errors.BadRequestErr.WithReason(err.Error()))
	}

	context := ctx.Context()
	resp, errResp := ash.ServiceController.BulkInsertServiceTx(context, nil, applicationID, services)
	if errResp != nil {
		return responses.FailedResponseWithError(ctx, errResp)
	}

	if resp.FailedRowsData != nil {
		return responses.FailedResponseV2(ctx, fiber.StatusBadRequest, resp)
	}

	return responses.SuccessResponseV2(ctx, fiber.StatusOK, resp, nil)
}

func (ash *ApplicationServiceHandler) GetApplicationServices(ctx *fiber.Ctx) error {
	if err := ash.validateUserAccess(ctx); err != nil {
		return responses.FailedResponseWithError(ctx, err)
	}

	applicationID, err := ash.getParamApplicationID(ctx)
	if err != nil {
		return responses.FailedResponseWithError(ctx, errors.BadRequestErr.WithReason(err.Error()))
	}

	paginate := helpers.Paginate{
		ValidOrderOptions: applicationEntity.AllowedServiceFields,
	}

	paginateData, errResp := paginate.NewFromContext(ctx)
	if errResp != nil {
		return responses.FailedResponseWithError(ctx, errResp)
	}

	context := ctx.Context()
	servicesData, total, errResp := ash.ServiceController.GetServices(context, applicationID, *paginateData)
	if errResp != nil {
		return responses.FailedResponseWithError(ctx, errResp)
	}

	return responses.SuccessResponseWithPaginate(ctx, fiber.StatusOK, servicesData, total, paginateData, nil)
}

func (ash *ApplicationServiceHandler) GetApplicationService(ctx *fiber.Ctx) error {
	if err := ash.validateUserAccess(ctx); err != nil {
		return responses.FailedResponseWithError(ctx, err)
	}

	applicationID, err := ash.getParamApplicationID(ctx)
	if err != nil {
		return responses.FailedResponseWithError(ctx, errors.BadRequestErr.WithReason(err.Error()))
	}

	serviceID, err := ash.getParamServiceID(ctx)
	if err != nil {
		return responses.FailedResponseWithError(ctx, errors.BadRequestErr.WithReason(err.Error()))
	}

	context := ctx.Context()
	servicesData, errResp := ash.ServiceController.GetServiceByID(context, applicationID, serviceID)
	if errResp != nil {
		return responses.FailedResponseWithError(ctx, errResp)
	}

	return responses.SuccessResponseV2(ctx, fiber.StatusOK, servicesData, nil)
}

func (ash *ApplicationServiceHandler) DeactivateService(ctx *fiber.Ctx) error {
	if err := ash.validateUserAccess(ctx); err != nil {
		return responses.FailedResponseWithError(ctx, err)
	}

	applicationID, err := ash.getParamApplicationID(ctx)
	if err != nil {
		return responses.FailedResponseWithError(ctx, errors.BadRequestErr.WithReason(err.Error()))
	}

	serviceID, err := ash.getParamServiceID(ctx)
	if err != nil {
		return responses.FailedResponseWithError(ctx, errors.BadRequestErr.WithReason(err.Error()))
	}

	context := ctx.Context()
	resp, respErr := ash.ServiceController.ActivateServiceByID(context, applicationID, serviceID, false)
	if respErr != nil {
		return responses.FailedResponseWithError(ctx, respErr)
	}

	return responses.SuccessResponseV2(ctx, fiber.StatusOK, resp, nil)
}

func (ash *ApplicationServiceHandler) ActivateService(ctx *fiber.Ctx) error {
	if err := ash.validateUserAccess(ctx); err != nil {
		return responses.FailedResponseWithError(ctx, err)
	}

	applicationID, err := ash.getParamApplicationID(ctx)
	if err != nil {
		return responses.FailedResponseWithError(ctx, errors.BadRequestErr.WithReason(err.Error()))
	}

	serviceID, err := ash.getParamServiceID(ctx)
	if err != nil {
		return responses.FailedResponseWithError(ctx, errors.BadRequestErr.WithReason(err.Error()))
	}

	context := ctx.Context()
	resp, respErr := ash.ServiceController.ActivateServiceByID(context, applicationID, serviceID, true)
	if respErr != nil {
		return responses.FailedResponseWithError(ctx, respErr)
	}

	return responses.SuccessResponseV2(ctx, fiber.StatusOK, resp, nil)
}

func (ash *ApplicationServiceHandler) UpdateService(ctx *fiber.Ctx) error {
	if err := ash.validateUserAccess(ctx); err != nil {
		return responses.FailedResponseWithError(ctx, err)
	}

	applicationID, err := ash.getParamApplicationID(ctx)
	if err != nil {
		return responses.FailedResponseWithError(ctx, errors.BadRequestErr.WithReason(err.Error()))
	}

	serviceID, err := ash.getParamServiceID(ctx)
	if err != nil {
		return responses.FailedResponseWithError(ctx, errors.BadRequestErr.WithReason(err.Error()))
	}

	data := applicationEntity.Service{}
	err = ctx.BodyParser(&data)
	if err != nil {
		return responses.FailedResponseWithError(ctx, errors.FailedParseRequestBodyErr.WithReason(err.Error()))
	}

	context := ctx.Context()
	resp, respErr := ash.ServiceController.UpdateServiceByID(context, applicationID, serviceID, data)
	if respErr != nil {
		return responses.FailedResponseWithError(ctx, respErr)
	}

	return responses.SuccessResponseV2(ctx, fiber.StatusOK, resp, nil)
}

func (ash *ApplicationServiceHandler) Register(router fiber.Router) error {
	internal := router.Group(fmt.Sprintf(applicationServiceEndpoint, core.V1, core.AccessInternal), ash.HandleInternalAccess())
	internal.Post("/", ash.CreateApplicationService)
	internal.Get("/", ash.GetApplicationServices)
	internal.Get("/:service_id", ash.GetApplicationService)
	internal.Put("/:service_id", ash.UpdateService)
	internal.Post("/:service_id/activate", ash.ActivateService)
	internal.Post("/:service_id/deactivate", ash.DeactivateService)

	return nil
}
