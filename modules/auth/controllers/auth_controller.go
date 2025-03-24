package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/winartodev/apollo/core"
	"github.com/winartodev/apollo/core/helpers"
	authEntity "github.com/winartodev/apollo/modules/auth/entities"
	userController "github.com/winartodev/apollo/modules/user/controllers"
	userEntity "github.com/winartodev/apollo/modules/user/entities"
)

type AuthControllerItf interface {
	SignIn(ctx context.Context, data *authEntity.SignInRequest) (res *authEntity.AuthResponse, err error)
	SignUp(ctx context.Context, data *authEntity.SignUpRequest) (success bool, err error)
	SignOut(ctx context.Context, id int64) (success bool, err error)
	RefreshToken(ctx context.Context, providedRefreshToken string) (res *authEntity.AuthResponse, err error)
}

type AuthController struct {
	UserController userController.UserControllerItf
}

func NewAuthController(controller AuthController) AuthControllerItf {
	return &AuthController{
		UserController: controller.UserController,
	}
}

func (ac *AuthController) SignIn(ctx context.Context, data *authEntity.SignInRequest) (res *authEntity.AuthResponse, err error) {
	passwordHash, err := ac.UserController.GetPasswordByEmail(ctx, data.Email)
	if err != nil {
		return nil, err
	}

	verified := helpers.VerifyPassword(data.Password, *passwordHash)
	if !verified {
		return nil, errors.New("invalid password")
	}

	user, err := ac.UserController.GetUserByEmail(ctx, data.Email)
	if err != nil {
		return nil, err
	}

	jwt, err := helpers.NewJWT()
	if err != nil {
		return nil, err
	}

	token, err := jwt.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	err = ac.UserController.UpdateRefreshToken(ctx, false, user.ID, &token.RefreshToken)
	if err != nil && err != core.ErrRefreshTokenExists {
		log.Debug(err)
		return nil, err
	}

	if err == core.ErrRefreshTokenExists {
		token.RefreshToken = ""
	}

	res = &authEntity.AuthResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	return res, nil
}

func (ac *AuthController) SignUp(ctx context.Context, data *authEntity.SignUpRequest) (success bool, err error) {
	user := userEntity.User{
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		Username:    data.Username,
		Password:    &data.Password,
	}

	id, err := ac.UserController.CreateUser(ctx, user)
	if err != nil {
		return false, err
	}

	if id == nil || *id == 0 {
		return false, errors.New("user can't created")
	}

	return true, nil
}

func (ac *AuthController) SignOut(ctx context.Context, id int64) (success bool, err error) {
	// remove refresh_token from database by user id
	err = ac.UserController.UpdateRefreshToken(ctx, true, id, nil)
	if err != nil {
		return false, err
	}

	return true, err
}

func (ac *AuthController) RefreshToken(ctx context.Context, providedRefreshToken string) (res *authEntity.AuthResponse, err error) {
	if providedRefreshToken == "" {
		return nil, errors.New("refresh token is required")
	}

	jwt, err := helpers.NewJWT()
	if err != nil {
		return nil, errors.New("failed to initialize JWT instance")
	}

	claims, valid, err := jwt.VerifyToken(jwt.RefreshToken.SecretKey, providedRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify refresh token: %v", err)
	}

	if !valid {
		return nil, errors.New("invalid or expired refresh token")
	}

	userID, ok := claims["id"].(float64)
	if !ok {
		return nil, errors.New("invalid refresh token")
	}

	userRefreshToken, err := ac.UserController.GetRefreshTokenByID(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	if userRefreshToken == nil || *userRefreshToken != providedRefreshToken {
		return nil, errors.New("invalid refresh token")
	}

	user, err := ac.UserController.GetUserByID(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	token, err := jwt.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	err = ac.UserController.UpdateRefreshToken(ctx, true, user.ID, &token.RefreshToken)
	if err != nil {
		return nil, err
	}

	res = &authEntity.AuthResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	return res, err
}
