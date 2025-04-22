package entities

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core/helpers"
	"github.com/winartodev/apollo/modules/application"
)

type SignUpRequest struct {
	Email       string `json:"email" form:"email"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
	Username    string `json:"username" form:"username"`
	Password    string `json:"password" form:"password"`
}

func (sur *SignUpRequest) BuildFromValue(ctx *fiber.Ctx) (res *SignUpRequest, err error) {
	email, err := helpers.GetFormValue(ctx, "email", true)
	if err != nil {
		return nil, err
	}

	phoneNumber, err := helpers.GetFormValue(ctx, "phone_number", true)
	if err != nil {
		return nil, err
	}

	username, err := helpers.GetFormValue(ctx, "username", true)
	if err != nil {
		return nil, err
	}

	password, err := helpers.GetFormValue(ctx, "password", true)
	if err != nil {
		return nil, err
	}

	return &SignUpRequest{
		Email:       email,
		PhoneNumber: phoneNumber,
		Username:    username,
		Password:    password,
	}, nil
}

type SignInRequest struct {
	Email       string `json:"email" form:"email"`
	Password    string `json:"password" form:"password"`
	Application *application.Access
}

func (sir *SignInRequest) BuildFromValue(ctx *fiber.Ctx) (res *SignInRequest, err error) {
	email, err := helpers.GetFormValue(ctx, "email", true)
	if err != nil {
		return nil, err
	}

	password, err := helpers.GetFormValue(ctx, "password", true)
	if err != nil {
		return nil, err
	}

	return &SignInRequest{
		Email:    email,
		Password: password,
	}, nil
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	Application  *application.Access
}
