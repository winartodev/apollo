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
	"github.com/winartodev/apollo/modules/application/controllers"
	applicationEntity "github.com/winartodev/apollo/modules/application/entities"
	"reflect"
)

const (
	applicationsEndpoint = "%s%s/applications"
)

type ApplicationHandler struct {
	middlewares.Middleware
	ApplicationController controllers.ApplicationControllerItf
}

func NewApplicationHandler(handler ApplicationHandler) ApplicationHandler {
	return ApplicationHandler{
		Middleware:            handler.Middleware,
		ApplicationController: handler.ApplicationController,
	}
}

func (ah *ApplicationHandler) GetApplications(ctx *fiber.Ctx) error {
	_, err := helpers.GetUserIDFromFiberContext(ctx)
	if err != nil {
		return responses.FailedResponseV2(ctx, errors.AuthorizationErr.WithReason(err.Error()))
	}

	access, ctxErr := helpers.GetApplicationAccessFromFiberContext(ctx)
	if ctxErr != nil {
		return responses.FailedResponseV2(ctx, errors.AuthorizationErr.WithReason(ctxErr.Error()))
	}

	if access == nil || !reflect.DeepEqual(access, application.ApolloInternal) {
		return responses.FailedResponseV2(ctx, errors.UserApplicationHasNotAccess)
	}

	paginate := helpers.Paginate{
		ValidOrderOptions: applicationEntity.AllowedApplicationFields,
	}

	paginateData, errResp := paginate.NewFromContext(ctx)
	if errResp != nil {
		return responses.FailedResponseV2(ctx, errResp)
	}

	context := ctx.Context()
	resp, total, errResp := ah.ApplicationController.GetApplications(context, *paginateData)
	if errResp != nil {
		return responses.FailedResponseV2(ctx, errResp)
	}

	return responses.SuccessResponseWithPaginate(ctx, fiber.StatusOK, resp, total, paginateData, nil)
}

func (ah *ApplicationHandler) Register(router fiber.Router) error {
	internal := router.Group(fmt.Sprintf(applicationsEndpoint, core.V1, core.AccessInternal), ah.HandleInternalAccess())
	internal.Get("/", ah.GetApplications)

	return nil
}
