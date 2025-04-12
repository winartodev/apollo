package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core"
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
		return responses.FailedResponse(ctx, fiber.StatusUnauthorized, "Authentication required", err)
	}

	data := entities.Service{}
	err = ctx.BodyParser(&data)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusBadRequest, "Failed to parse request body", err)
	}

	data.CreatedBy = id
	data.UpdatedBy = id

	context := ctx.Context()
	resp, _, err := sh.ServiceController.CreateService(context, data)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusInternalServerError, "Failed to create service", err)
	}

	return responses.SuccessResponse(ctx, fiber.StatusCreated, "Your request has been processed successfully", resp, nil)
}

func (sh *ServiceHandler) Register(router fiber.Router) error {
	internal := router.Group(fmt.Sprintf(serviceEndpoint, core.V1, core.AccessInternal), sh.HandleInternalAccess())
	internal.Post("/", sh.Create)
	return nil
}
