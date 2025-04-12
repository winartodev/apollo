package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core"
	"github.com/winartodev/apollo/core/helpers"
	"github.com/winartodev/apollo/core/middlewares"
	"github.com/winartodev/apollo/core/responses"
	authController "github.com/winartodev/apollo/modules/auth/controllers"
	"github.com/winartodev/apollo/modules/auth/emums"
	authEntity "github.com/winartodev/apollo/modules/auth/entities"
)

type AuthHandler struct {
	middlewares.Middleware
	VerificationController authController.VerificationControllerItf
	AuthController         authController.AuthControllerItf
}

func NewAuthHandler(handler AuthHandler) AuthHandler {
	return AuthHandler{
		Middleware:             handler.Middleware,
		VerificationController: handler.VerificationController,
		AuthController:         handler.AuthController,
	}
}

func (h *AuthHandler) SignIn(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := authEntity.SignInRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusBadRequest, "Failed to sign in", err)
	}

	res, err := h.AuthController.SignIn(context, &req)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusInternalServerError, "Failed to sign in", err)
	}

	return responses.SuccessResponse(ctx, fiber.StatusOK, "Success", res, nil)
}

func (h *AuthHandler) SignUp(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := authEntity.SignUpRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusBadRequest, "Failed to create user account", err)
	}

	res, err := h.AuthController.SignUp(context, &req)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusInternalServerError, "Failed to create user account", err)
	}

	return responses.SuccessResponse(ctx, fiber.StatusCreated, "User account created successfully", res, nil)
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

func (h *AuthHandler) GenerateEmailOTP(ctx *fiber.Ctx) error {
	context := ctx.Context()

	email := ctx.Query("email", "")
	if email == "" {
		return responses.FailedResponse(ctx, fiber.StatusBadRequest, "Invalid email address", nil)
	}

	err := h.VerificationController.CreateOTP(context, emums.VerificationEmail, email)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusInternalServerError, "Failed to create otp code", err)
	}

	return responses.SuccessResponse(ctx, fiber.StatusOK, "Success", "OTP send successfully", nil)
}

func (h *AuthHandler) ValidateEmailOTP(ctx *fiber.Ctx) error {
	context := ctx.Context()

	otp := ctx.Query("otp", "")
	email := ctx.Query("email", "")
	if otp == "" && email == "" {
		return responses.FailedResponse(ctx, fiber.StatusBadRequest, "Failed Validate OTP", nil)
	}

	err := h.VerificationController.VerifyOTP(context, emums.VerificationEmail, email, otp)
	if err != nil && !errors.Is(err, authController.ErrorOTPAlreadyVerified) {
		return responses.FailedResponse(ctx, fiber.StatusInternalServerError, "Failed to create otp code", err)
	}

	if errors.Is(err, authController.ErrorOTPAlreadyVerified) {
		data := authController.ErrorOTPAlreadyVerified.Error()
		return responses.SuccessResponse(ctx, fiber.StatusOK, "Success", data, nil)
	}

	return responses.SuccessResponse(ctx, fiber.StatusOK, "Success", "validate OTP successfully", nil)
}

func (h *AuthHandler) ResendEmailOTP(ctx *fiber.Ctx) error {
	context := ctx.Context()

	email := ctx.Query("email", "")
	if email == "" {
		return responses.FailedResponse(ctx, fiber.StatusBadRequest, "Failed Resend OTP", nil)
	}

	err := h.VerificationController.ResendOTP(context, emums.VerificationEmail, email)
	if err != nil && !errors.Is(err, authController.ErrorOTPAlreadyVerified) {
		return responses.FailedResponse(ctx, fiber.StatusInternalServerError, "Failed to create otp code", err)
	}

	if errors.Is(err, authController.ErrorOTPAlreadyVerified) {
		data := authController.ErrorOTPAlreadyVerified.Error()
		return responses.SuccessResponse(ctx, fiber.StatusOK, "Success", data, nil)
	}

	return responses.SuccessResponse(ctx, fiber.StatusOK, "Success", "resend OTP successfully", nil)
}

func (h *AuthHandler) GeneratePhoneOTP(ctx *fiber.Ctx) error {
	context := ctx.Context()

	phone := ctx.Query("phone", "")
	if phone == "" {
		return responses.FailedResponse(ctx, fiber.StatusBadRequest, "Invalid phone number", nil)
	}

	newPhone, err := helpers.FormatIndonesianPhoneNumber(phone)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusBadRequest, "Invalid phone number", nil)
	}

	err = h.VerificationController.CreateOTP(context, emums.VerificationPhone, newPhone)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusInternalServerError, "Failed to create otp code", err)
	}

	return responses.SuccessResponse(ctx, fiber.StatusOK, "Success", "OTP send successfully", nil)
}

func (h *AuthHandler) ValidatePhoneOTP(ctx *fiber.Ctx) error {
	context := ctx.Context()

	otp := ctx.Query("otp", "")
	phone := ctx.Query("phone", "")
	if otp == "" && phone == "" {
		return responses.FailedResponse(ctx, fiber.StatusBadRequest, "Failed Validate OTP", nil)
	}

	newPhone, err := helpers.FormatIndonesianPhoneNumber(phone)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusBadRequest, "Invalid phone number", nil)
	}

	err = h.VerificationController.VerifyOTP(context, emums.VerificationPhone, newPhone, otp)
	if err != nil && !errors.Is(err, authController.ErrorOTPAlreadyVerified) {
		return responses.FailedResponse(ctx, fiber.StatusInternalServerError, "Failed to create otp code", err)
	}

	if errors.Is(err, authController.ErrorOTPAlreadyVerified) {
		data := authController.ErrorOTPAlreadyVerified.Error()
		return responses.SuccessResponse(ctx, fiber.StatusOK, "Success", data, nil)
	}

	return responses.SuccessResponse(ctx, fiber.StatusOK, "Success", "validate OTP successfully", nil)
}

func (h *AuthHandler) ResendPhoneOTP(ctx *fiber.Ctx) error {
	context := ctx.Context()

	phone := ctx.Query("phone", "")
	if phone == "" {
		return responses.FailedResponse(ctx, fiber.StatusBadRequest, "Failed Resend OTP", nil)
	}

	newPhone, err := helpers.FormatIndonesianPhoneNumber(phone)
	if err != nil {
		return responses.FailedResponse(ctx, fiber.StatusBadRequest, "Invalid phone number", nil)
	}

	err = h.VerificationController.ResendOTP(context, emums.VerificationPhone, newPhone)
	if err != nil && !errors.Is(err, authController.ErrorOTPAlreadyVerified) {
		return responses.FailedResponse(ctx, fiber.StatusInternalServerError, "Failed to create otp code", err)
	}

	if errors.Is(err, authController.ErrorOTPAlreadyVerified) {
		data := authController.ErrorOTPAlreadyVerified.Error()
		return responses.SuccessResponse(ctx, fiber.StatusOK, "Success", data, nil)
	}

	return responses.SuccessResponse(ctx, fiber.StatusOK, "Success", "resend OTP successfully", nil)
}

func (h *AuthHandler) Register(router fiber.Router) error {
	v1 := router.Group(core.V1)

	auth := v1.Group("/auth")

	auth.Post("/sign-in", h.SignIn)
	auth.Post("/sign-up", h.SignUp)
	auth.Post("/refresh", h.RefreshToken)

	otp := auth.Group("/otp")
	otp.Post("/email", h.GenerateEmailOTP)
	otp.Post("/email/validate", h.ValidateEmailOTP)
	otp.Post("/email/resend", h.ResendEmailOTP)

	otp.Post("/phone", h.GeneratePhoneOTP)
	otp.Post("/phone/validate", h.ValidatePhoneOTP)
	otp.Post("/phone/resend", h.ResendPhoneOTP)

	userAuth := v1.Group("/users/auth", h.HandlePublicAccess())
	userAuth.Post("/sign-out", h.SignOut)

	return nil
}
