package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core"
	"github.com/winartodev/apollo/core/helpers"
	"github.com/winartodev/apollo/core/middlewares"
	"github.com/winartodev/apollo/core/responses"
	authController "github.com/winartodev/apollo/modules/auth/controllers"
	authEntity "github.com/winartodev/apollo/modules/auth/entities"
)

type AuthHandler struct {
	AuthController authController.AuthControllerItf
}

func NewAuthHandler(handler AuthHandler) AuthHandler {
	return AuthHandler{
		AuthController: handler.AuthController,
	}
}

func (h *AuthHandler) SignIn(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := authEntity.SignInRequest{}
	data, err := req.BuildFromValue(ctx)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusBadRequest, "Invalid sign-up request", err)
	}

	res, err := h.AuthController.SignIn(context, data)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusInternalServerError, "Failed to sign in", err)
	}

	return responses.SuccessResponse(ctx, fiber.StatusOK, "Success", res, nil)
}

func (h *AuthHandler) SignUp(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := authEntity.SignUpRequest{}
	data, err := req.BuildFromValue(ctx)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusBadRequest, "Invalid sign-up request", err)
	}

	_, err = h.AuthController.SignUp(context, data)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusInternalServerError, "Failed to create user account", err)
	}

	return responses.SuccessResponse(ctx, fiber.StatusCreated, "User account created successfully", nil, nil)
}

func (h *AuthHandler) SignOut(ctx *fiber.Ctx) error {
	context := ctx.Context()

	id, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusBadRequest, "Failed Sign Out", err)
	}

	if id > 0 {
		_, err = h.AuthController.SignOut(context, id)
		if err != nil {
			return responses.FailedResponse(ctx, fiber.StatusInternalServerError, "Failed Sign Out", err)
		}
	}

	return responses.SuccessResponse(ctx, fiber.StatusOK, "Success", nil, nil)
}

func (h *AuthHandler) RefreshToken(ctx *fiber.Ctx) error {
	context := ctx.Context()
	req := authEntity.RefreshTokenRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusBadRequest, "Failed Refresh Token", err)
	}

	res, err := h.AuthController.RefreshToken(context, req.RefreshToken)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusInternalServerError, "Failed to refresh token", err)
	}

	return responses.SuccessResponse(ctx, fiber.StatusOK, "Success", res, nil)
}

func (h *AuthHandler) Register(router fiber.Router) error {
	v1 := router.Group(core.V1)

	auth := v1.Group("/auth")
	auth.Post("/sign-in", h.SignIn)
	auth.Post("/sign-up", h.SignUp)
	auth.Post("/refresh", h.RefreshToken)

	userAuth := v1.Group("/users/auth", middlewares.HandlePublicAccess())
	userAuth.Post("/sign-out", h.SignOut)

	return nil
}
