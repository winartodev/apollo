package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core"
	"github.com/winartodev/apollo/core/errors"
	"github.com/winartodev/apollo/core/helpers"
	"github.com/winartodev/apollo/core/middlewares"
	"github.com/winartodev/apollo/core/responses"
	applicationController "github.com/winartodev/apollo/modules/application/controllers"
	applicationEntity "github.com/winartodev/apollo/modules/application/entities"
)

const (
	serviceEndpoint = "%s%s/services"
	invalidParamID  = "param id must greater than 0"
)

type ServiceHandler struct {
	middlewares.Middleware
	ServiceController applicationController.ServiceControllerItf
}

func NewServiceHandler(handler ServiceHandler) ServiceHandler {
	return ServiceHandler{
		Middleware:        handler.Middleware,
		ServiceController: handler.ServiceController,
	}
}

func (sh *ServiceHandler) Create(ctx *fiber.Ctx) error {
	userID, err := helpers.GetUserIDFromFiberContext(ctx)
	if err != nil {
		return responses.FailedResponseV2(ctx, errors.AuthorizationErr.WithReason(err.Error()))
	}

	data := applicationEntity.Service{}
	err = ctx.BodyParser(&data)
	if err != nil {
		return responses.FailedResponseV2(ctx, errors.FailedParseRequestBodyErr.WithReason(err.Error()))
	}

	data.CreatedBy = userID
	data.UpdatedBy = userID

	context := ctx.Context()
	resp, errResp := sh.ServiceController.CreateService(context, data)
	if errResp != nil {
		return responses.FailedResponseV2(ctx, errResp)
	}

	return responses.SuccessResponseV2(ctx, fiber.StatusCreated, resp, nil)
}

func (sh *ServiceHandler) GetServiceByID(ctx *fiber.Ctx) error {
	_, err := helpers.GetUserIDFromFiberContext(ctx)
	if err != nil {
		return responses.FailedResponseV2(ctx, errors.AuthorizationErr.WithReason(err.Error()))
	}

	paramID, err := ctx.ParamsInt("id", 0)
	if err != nil {
		return responses.FailedResponseV2(ctx, errors.BadRequestErr.WithReason(err.Error()))
	}

	context := ctx.Context()
	if paramID <= 0 {
		return responses.FailedResponseV2(ctx, errors.BadRequestErr.WithReason(invalidParamID))
	}

	id := int64(paramID)
	resp, respErr := sh.ServiceController.GetServiceByID(context, id)
	if respErr != nil {
		return responses.FailedResponseV2(ctx, respErr)
	}

	return responses.SuccessResponseV2(ctx, fiber.StatusOK, resp, nil)

}

func (sh *ServiceHandler) GetServices(ctx *fiber.Ctx) error {
	_, err := helpers.GetUserIDFromFiberContext(ctx)
	if err != nil {
		return responses.FailedResponseV2(ctx, errors.AuthorizationErr.WithReason(err.Error()))
	}

	paginate := helpers.Paginate{
		ValidOrderOptions: applicationEntity.AllowedServiceFields,
	}

	paginateData, errResp := paginate.NewFromContext(ctx)
	if errResp != nil {
		return responses.FailedResponseV2(ctx, errResp)
	}

	context := ctx.Context()
	resp, total, errResp := sh.ServiceController.GetServices(context, *paginateData)
	if errResp != nil {
		return responses.FailedResponseV2(ctx, errResp)
	}

	return responses.SuccessResponseWithPaginate(ctx, fiber.StatusOK, resp, total, paginateData, nil)
}

func (sh *ServiceHandler) Deactivate(ctx *fiber.Ctx) error {
	_, err := helpers.GetUserIDFromFiberContext(ctx)
	if err != nil {
		return responses.FailedResponseV2(ctx, errors.AuthorizationErr.WithReason(err.Error()))
	}

	serviceID, err := ctx.ParamsInt("id", 0)
	if err != nil {
		return responses.FailedResponseV2(ctx, errors.BadRequestErr.WithReason(err.Error()))
	}

	context := ctx.Context()
	if serviceID <= 0 {
		return responses.FailedResponseV2(ctx, errors.BadRequestErr.WithReason(invalidParamID))
	}

	successData, respErr := sh.ServiceController.ActivateServiceByID(context, int64(serviceID), false)
	if respErr != nil {
		return responses.FailedResponseV2(ctx, respErr)
	}

	resp := responses.UpdateResponseData{
		Message:     fmt.Sprintf("service id %d is already inactive", serviceID),
		TotalData:   1,
		SuccessData: 1,
		FailData:    0,
		SuccessRowsData: []interface{}{
			successData,
		},
	}

	return responses.SuccessResponseV2(ctx, fiber.StatusOK, resp, nil)
}

func (sh *ServiceHandler) Activate(ctx *fiber.Ctx) error {
	_, err := helpers.GetUserIDFromFiberContext(ctx)
	if err != nil {
		return responses.FailedResponseV2(ctx, errors.AuthorizationErr.WithReason(err.Error()))
	}

	serviceID, err := ctx.ParamsInt("id", 0)
	if err != nil {
		return responses.FailedResponseV2(ctx, errors.BadRequestErr.WithReason(err.Error()))
	}

	context := ctx.Context()
	if serviceID <= 0 {
		return responses.FailedResponseV2(ctx, errors.BadRequestErr.WithReason(invalidParamID))
	}

	successData, respErr := sh.ServiceController.ActivateServiceByID(context, int64(serviceID), true)
	if respErr != nil {
		return responses.FailedResponseV2(ctx, respErr)
	}

	resp := responses.UpdateResponseData{
		Message:     fmt.Sprintf("service id %d is already active", serviceID),
		TotalData:   1,
		SuccessData: 1,
		FailData:    0,
		SuccessRowsData: []interface{}{
			successData,
		},
	}

	return responses.SuccessResponseV2(ctx, fiber.StatusOK, resp, nil)
}

func (sh *ServiceHandler) Update(ctx *fiber.Ctx) error {
	_, err := helpers.GetUserIDFromFiberContext(ctx)
	if err != nil {
		return responses.FailedResponseV2(ctx, errors.AuthorizationErr.WithReason(err.Error()))
	}

	serviceID, err := ctx.ParamsInt("id", 0)
	if err != nil {
		return responses.FailedResponseV2(ctx, errors.BadRequestErr.WithReason(err.Error()))
	}

	context := ctx.Context()
	if serviceID <= 0 {
		return responses.FailedResponseV2(ctx, errors.BadRequestErr.WithReason(invalidParamID))
	}

	data := applicationEntity.Service{}
	err = ctx.BodyParser(&data)
	if err != nil {
		return responses.FailedResponseV2(ctx, errors.FailedParseRequestBodyErr.WithReason(err.Error()))
	}

	successData, respErr := sh.ServiceController.UpdateServiceByID(context, int64(serviceID), data)
	if respErr != nil {
		return responses.FailedResponseV2(ctx, respErr)
	}

	resp := responses.UpdateResponseData{
		TotalData:   1,
		SuccessData: 1,
		FailData:    0,
		SuccessRowsData: []interface{}{
			successData,
		},
	}

	return responses.SuccessResponseV2(ctx, fiber.StatusOK, resp, nil)
}

func (sh *ServiceHandler) Register(router fiber.Router) error {
	internal := router.Group(fmt.Sprintf(serviceEndpoint, core.V1, core.AccessInternal), sh.HandleInternalAccess())
	internal.Post("/", sh.Create)
	internal.Get("/:id", sh.GetServiceByID)
	internal.Get("/", sh.GetServices)
	internal.Patch("/:id", sh.Update)
	internal.Patch("/:id/activate", sh.Activate)
	internal.Patch("/:id/deactivate", sh.Deactivate)
	return nil
}
