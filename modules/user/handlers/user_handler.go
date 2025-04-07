package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core"
	"github.com/winartodev/apollo/core/helpers"
	"github.com/winartodev/apollo/core/middlewares"
	"github.com/winartodev/apollo/core/responses"
	enums2 "github.com/winartodev/apollo/modules/application/enums"
	userControler "github.com/winartodev/apollo/modules/user/controllers"
)

var (
	userNotLoggedIn = errors.New("not logged in")

	apolloInternalUserAccess = &middlewares.InternalAccessConfig{
		Application:        enums2.ApolloInternal,
		ApplicationService: enums2.TestInternalServices1,
	}
)

type UserHandler struct {
	middlewares.Middleware
	UserController userControler.UserControllerItf
}

func NewUserHandler(handler UserHandler) UserHandler {
	return UserHandler{
		Middleware:     handler.Middleware,
		UserController: handler.UserController,
	}
}

func (h *UserHandler) GetCurrentUser(ctx *fiber.Ctx) error {
	context := ctx.Context()
	id, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusBadRequest, "Failed Get Current User", userNotLoggedIn)
	}

	res, err := h.UserController.GetUserByID(context, id)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusInternalServerError, "Failed Get Current User", err)
	}

	return responses.SuccessResponse(ctx, fiber.StatusOK, "Success Get Current User", res, nil)
}

func (h *UserHandler) Register(router fiber.Router) error {
	v1 := router.Group(core.V1)
	internal := v1.Group(core.AccessInternal)
	user := internal.Group("/users", h.HandleInternalAccess(apolloInternalUserAccess))
	user.Get("/me", h.GetCurrentUser)

	return nil
}
