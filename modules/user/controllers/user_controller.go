package controllers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/winartodev/apollo/core"
	"github.com/winartodev/apollo/core/helpers"
	userEntity "github.com/winartodev/apollo/modules/user/entities"
	userRepo "github.com/winartodev/apollo/modules/user/repositories"
)

var (
	ErrorUserNotFound = errors.New("user not found")
)

type UserControllerItf interface {
	CreateUser(ctx context.Context, data userEntity.User) (id *int64, err error)
	UpdateRefreshToken(ctx context.Context, force bool, id int64, refreshToken *string) (err error)
	GetUserByID(ctx context.Context, id int64) (res *userEntity.User, err error)
	GetUserByEmail(ctx context.Context, email string) (res *userEntity.User, err error)
	GetPasswordByEmail(ctx context.Context, email string) (res *string, err error)
	GetRefreshTokenByID(ctx context.Context, id int64) (res *string, err error)
	ValidateUserIsExists(ctx context.Context, data *userEntity.User) (err error)
	GetUserRoleByID(ctx context.Context, id int64) (res *userEntity.UserRole, err error)
}

type UserController struct {
	UserRepository userRepo.UserRepositoryItf
}

func NewUserController(controller UserController) UserControllerItf {
	return &UserController{
		UserRepository: controller.UserRepository,
	}
}

func (uc *UserController) CreateUser(ctx context.Context, data userEntity.User) (id *int64, err error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	passwordHash, err := helpers.HashPassword(*data.Password)
	if err != nil {
		return nil, err
	}

	data.UUID = newUUID.String()
	data.Password = &passwordHash

	res, err := uc.UserRepository.CreateUserDB(ctx, &data)
	if err != nil {
		return nil, err
	}

	return &res, err
}

func (uc *UserController) GetUserByID(ctx context.Context, id int64) (res *userEntity.User, err error) {
	res, err = uc.UserRepository.GetUserByIDDB(ctx, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if res == nil {
		return nil, ErrorUserNotFound
	}

	return res, nil
}

func (uc *UserController) GetUserByEmail(ctx context.Context, email string) (res *userEntity.User, err error) {
	res, err = uc.UserRepository.GetUserByEmailDB(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if res == nil {
		return nil, ErrorUserNotFound
	}

	return res, nil
}

func (uc *UserController) GetPasswordByEmail(ctx context.Context, email string) (res *string, err error) {
	res, err = uc.UserRepository.GetUserPasswordByEmailDB(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if res == nil {
		return nil, ErrorUserNotFound
	}

	return res, nil
}

func (uc *UserController) GetRefreshTokenByID(ctx context.Context, id int64) (res *string, err error) {
	res, err = uc.UserRepository.GetRefreshTokenByIDDB(ctx, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if res == nil {
		return nil, ErrorUserNotFound
	}

	return res, nil
}

func (uc *UserController) GetUserRoleByID(ctx context.Context, id int64) (res *userEntity.UserRole, err error) {
	return uc.UserRepository.GetUserRoleByIDDB(ctx, id)
}

func (uc *UserController) UpdateRefreshToken(ctx context.Context, force bool, id int64, refreshToken *string) (err error) {
	if !force {
		exists, err := uc.UserRepository.IsRefreshTokenExistByIDDB(ctx, id)
		if err != nil {
			return err
		}

		if exists {
			return core.ErrRefreshTokenExists
		}
	}

	err = uc.UserRepository.UpdateRefreshTokenByIDDB(ctx, id, refreshToken)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserController) ValidateUserIsExists(ctx context.Context, data *userEntity.User) (err error) {
	res, err := uc.UserRepository.IsUserExistsDB(ctx, &userEntity.UserUniqueField{
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		Username:    data.Username,
	})
	if err != nil {
		return err
	}

	if res == nil {
		return errors.New("result is nil")
	}

	if res.IsEmailExists {
		return errors.New("email is exists")
	}

	if res.IsPhoneExists {
		return errors.New("phone is exists")
	}

	if res.IsUsernameExists {
		return errors.New("username is exists")
	}

	return nil
}
