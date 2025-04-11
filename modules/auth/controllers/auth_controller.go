package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/winartodev/apollo/core"
	"github.com/winartodev/apollo/core/configs"
	"github.com/winartodev/apollo/core/helpers"
	authEnum "github.com/winartodev/apollo/modules/auth/emums"
	authEntity "github.com/winartodev/apollo/modules/auth/entities"
	userController "github.com/winartodev/apollo/modules/user/controllers"
	userEntity "github.com/winartodev/apollo/modules/user/entities"
)

type AuthControllerItf interface {
	SignIn(ctx context.Context, data *authEntity.SignInRequest) (res *authEntity.AuthResponse, err error)
	SignUp(ctx context.Context, data *authEntity.SignUpRequest) (res *userEntity.User, err error)
	SignOut(ctx context.Context, id int64) (success bool, err error)
	RefreshToken(ctx context.Context, providedRefreshToken string) (res *authEntity.AuthResponse, err error)
}

type AuthController struct {
	OTP                    *configs.OTP
	VerificationController VerificationControllerItf
	UserController         userController.UserControllerItf
}

func NewAuthController(controller AuthController) AuthControllerItf {
	return &AuthController{
		OTP:                    controller.OTP,
		VerificationController: controller.VerificationController,
		UserController:         controller.UserController,
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

func (ac *AuthController) SignUp(ctx context.Context, data *authEntity.SignUpRequest) (res *userEntity.User, err error) {
	newPhone, err := helpers.FormatIndonesianPhoneNumber(data.PhoneNumber)
	if err != nil {
		return nil, err
	}

	user := userEntity.User{
		Email:       data.Email,
		PhoneNumber: newPhone,
		Username:    data.Username,
		Password:    &data.Password,
	}

	err = ac.UserController.ValidateUserIsExists(ctx, &user)
	if err != nil {
		return nil, err
	}

	if ac.OTP.Enable {
		err = ac.CheckOTPVerificationOTP(ctx, &user)
		if err != nil {
			return nil, err
		}
	}

	newUser, err := ac.UserController.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	if newUser == nil {
		return nil, errors.New("user can't created")
	}

	err = ac.VerificationController.DeleteOTP(ctx, authEnum.VerificationPhone, user.PhoneNumber)
	if err != nil && err != ErrorOTPDataEmpty {
		return nil, err
	}

	err = ac.VerificationController.DeleteOTP(ctx, authEnum.VerificationEmail, user.Email)
	if err != nil && err != ErrorOTPDataEmpty {
		return nil, err
	}

	return newUser, nil
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

func (ac *AuthController) CheckOTPVerificationOTP(ctx context.Context, user *userEntity.User) (err error) {
	data := *user
	otpPhone, err := ac.VerificationController.GetOTP(ctx, authEnum.VerificationPhone, data.PhoneNumber)
	if err != nil {
		return err
	}

	otpEmail, err := ac.VerificationController.GetOTP(ctx, authEnum.VerificationEmail, data.Email)
	if err != nil {
		return err
	}

	if otpPhone == nil || !otpPhone.IsVerified {
		return errors.New("phone number is not verified")
	}

	if otpEmail == nil || !otpEmail.IsVerified {
		return errors.New("email is not verified")
	}

	data.IsEmailVerified = true
	data.IsPhoneVerified = true

	*user = data

	return nil
}
