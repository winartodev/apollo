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
	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		return responses.FailedResponseV2(ctx, errors.AuthorizationErr(err.Error()))
	}

	data := applicationEntity.Service{}
	err = ctx.BodyParser(&data)
	if err != nil {
		return responses.FailedResponseV2(ctx, errors.FailedParseRequestBodyErr(err.Error()))
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
	_, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		return responses.FailedResponseV2(ctx, errors.AuthorizationErr(err.Error()))
	}

	paramID, err := ctx.ParamsInt("id", 0)
	if err != nil {
		return responses.FailedResponseV2(ctx, errors.BadRequestWithReasonErr(err.Error()))
	}

	context := ctx.Context()
	if paramID > 0 {
		id := int64(paramID)
		resp, respErr := sh.ServiceController.GetServiceByID(context, id)
		if respErr != nil {
			return responses.FailedResponseV2(ctx, respErr)
		}

		return responses.SuccessResponseV2(ctx, fiber.StatusOK, resp, nil)
	}

	return responses.FailedResponseV2(ctx, errors.BadRequestWithReasonErr(invalidParamID))
}

func (sh *ServiceHandler) GetServices(ctx *fiber.Ctx) error {
	_, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		return responses.FailedResponseV2(ctx, errors.AuthorizationErr(err.Error()))
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

func (sh *ServiceHandler) Register(router fiber.Router) error {
	internal := router.Group(fmt.Sprintf(serviceEndpoint, core.V1, core.AccessInternal), sh.HandleInternalAccess())
	internal.Post("/", sh.Create)
	internal.Get("/:id", sh.GetServiceByID)
	internal.Get("/", sh.GetServices)
	return nil
}
