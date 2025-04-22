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
	authController "github.com/winartodev/apollo/modules/auth/controllers"
	authEntity "github.com/winartodev/apollo/modules/auth/entities"
	userController "github.com/winartodev/apollo/modules/user/controllers"
)

const (
	authApolloInternalEndpoint     = "%s%s/apollo/auth"
	authUserApolloInternalEndpoint = "%s%s/apollo/user/auth"
)

type AuthApolloInternalHandler struct {
	middlewares.Middleware
	ApplicationController applicationController.ApplicationControllerItf
	AuthController        authController.AuthControllerItf
	UserController        userController.UserControllerItf
}

func NewAuthApolloInternalHandler(handler AuthApolloInternalHandler) AuthApolloInternalHandler {
	return AuthApolloInternalHandler{
		Middleware:            handler.Middleware,
		ApplicationController: handler.ApplicationController,
		AuthController:        handler.AuthController,
		UserController:        handler.UserController,
	}
}

func (ah *AuthApolloInternalHandler) SignIn(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := authEntity.SignInRequest{
		Application: application.ApolloInternal,
	}

	err := ctx.BodyParser(&req)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusBadRequest, "Failed to sign in", err)
	}

	userData, err := ah.UserController.GetUserByEmail(context, req.Email)
	if err != nil && err != userController.ErrorUserNotFound {
		return responses.FailedResponse(ctx, fiber.StatusInternalServerError, "Failed to sign in", err)
	}

	if userData == nil {
		return responses.FailedResponseWithError(ctx, errors.DataNotFoundErr.WithReason("user data not found"))
	}

	respErr := ah.UserController.GetUserApplicationAccess(context, userData.ID, req.Application.ID.ToInt64(), req.Application.Scope.ToInt64())
	if respErr != nil {
		return responses.FailedResponseWithError(ctx, respErr)
	}

	res, err := ah.AuthController.SignIn(context, &req)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusInternalServerError, "Failed to sign in", err)
	}

	return responses.SuccessResponse(ctx, fiber.StatusOK, "Success", res, nil)
}

func (ah *AuthApolloInternalHandler) SignOut(ctx *fiber.Ctx) error {
	context := ctx.Context()

	id, err := helpers.GetUserIDFromFiberContext(ctx)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusBadRequest, "Failed Sign Out", err)
	}

	if id > 0 {
		_, err = ah.AuthController.SignOut(context, id)
		if err != nil {
			return responses.FailedResponse(ctx, fiber.StatusInternalServerError, "Failed Sign Out", err)
		}
	}

	return responses.SuccessResponse(ctx, fiber.StatusOK, "Success", nil, nil)
}

func (ah *AuthApolloInternalHandler) RefreshToken(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := authEntity.RefreshTokenRequest{
		Application: application.ApolloInternal,
	}

	err := ctx.BodyParser(&req)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusBadRequest, "Failed Refresh Token", err)
	}

	res, err := ah.AuthController.RefreshToken(context, req.RefreshToken)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusInternalServerError, "Failed to refresh token", err)
	}

	return responses.SuccessResponse(ctx, fiber.StatusOK, "Success", res, nil)
}

func (ah *AuthApolloInternalHandler) Register(router fiber.Router) error {
	internalNoAuth := router.Group(fmt.Sprintf(authApolloInternalEndpoint, core.V1, core.AccessInternal))
	internalNoAuth.Post("/sign-in", ah.SignIn)
	internalNoAuth.Post("/refresh", ah.RefreshToken)

	internalWithAuth := router.Group(fmt.Sprintf(authUserApolloInternalEndpoint, core.V1, core.AccessInternal), ah.HandleInternalAccess())
	internalWithAuth.Post("/sign-out", ah.SignOut)
	return nil
}
