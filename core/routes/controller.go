package routes

import (
	authController "github.com/winartodev/apollo/modules/auth/controllers"
	userController "github.com/winartodev/apollo/modules/user/controllers"
)

type ControllerDependency struct {
	Repository *Repository
}

type Controller struct {
	UserController userController.UserControllerItf
	AuthController authController.AuthControllerItf
}

func NewController(dependency ControllerDependency) *Controller {
	repository := dependency.Repository

	newUserController := userController.NewUserController(userController.UserController{
		UserRepository: repository.UserRepository,
	})

	newAuthController := authController.NewAuthController(authController.AuthController{
		UserController: newUserController,
	})

	return &Controller{
		UserController: newUserController,
		AuthController: newAuthController,
	}
}
