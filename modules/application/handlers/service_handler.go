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
	"github.com/winartodev/apollo/modules/application/entities"
)

const (
	serviceEndpoint = "%s%s/services"
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
	id, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		return responses.FailedResponseV2(ctx, errors.AuthorizationErr(err.Error()))
	}

	data := entities.Service{}
	err = ctx.BodyParser(&data)
	if err != nil {
		return responses.FailedResponseV2(ctx, errors.FailedParseRequestBodyErr(err.Error()))
	}

	data.CreatedBy = id
	data.UpdatedBy = id

	context := ctx.Context()
	resp, errResp := sh.ServiceController.CreateService(context, data)
	if errResp != nil {
		return responses.FailedResponseV2(ctx, errResp)
	}

	return responses.SuccessResponseV2(ctx, fiber.StatusCreated, resp, nil)
}

func (sh *ServiceHandler) Register(router fiber.Router) error {
	internal := router.Group(fmt.Sprintf(serviceEndpoint, core.V1, core.AccessInternal), sh.HandleInternalAccess())
	internal.Post("/", sh.Create)
	return nil
}
