package routes

import (
	"github.com/winartodev/apollo/core/configs"
	applicationController "github.com/winartodev/apollo/modules/application/controllers"
	authController "github.com/winartodev/apollo/modules/auth/controllers"
	userController "github.com/winartodev/apollo/modules/user/controllers"
)

type ControllerDependency struct {
	OTP        *configs.OTP
	SMTPClient *configs.SMTPClient
	Twilio     *configs.TwilioClient
	Repository *Repository
}

type Controller struct {
	UserController         userController.UserControllerItf
	VerificationController authController.VerificationControllerItf
	AuthController         authController.AuthControllerItf
	ServiceController      applicationController.ServiceControllerItf
	ApplicationController  applicationController.ApplicationControllerItf
}

func NewController(dependency ControllerDependency) *Controller {
	repository := dependency.Repository

	newUserController := userController.NewUserController(userController.UserController{
		UserApplicationRepository: repository.UserApplicationRepository,
		UserRepository:            repository.UserRepository,
	})

	newVerificationController := authController.NewVerificationController(authController.VerificationController{
		OTP:                    dependency.OTP,
		SmtpClient:             dependency.SMTPClient,
		TwilioClient:           dependency.Twilio,
		VerificationRepository: repository.VerificationRepository,
	})

	newAuthController := authController.NewAuthController(authController.AuthController{
		OTP:                    dependency.OTP,
		VerificationController: newVerificationController,
		UserController:         newUserController,
	})

	newServiceController := applicationController.NewServiceController(applicationController.ServiceController{
		Tx:          repository.DBTransact,
		ServiceRepo: repository.ServiceRepository,
	})

	newApplicationController := applicationController.NewApplicationController(applicationController.ApplicationController{
		Tx:                   repository.DBTransact,
		ApplicationRepo:      repository.ApplicationRepository,
		ApplicationScopeRepo: repository.ApplicationScopeRepository,
		UserApplicationRepo:  repository.UserApplicationRepository,
		ServiceController:    newServiceController,
	})

	return &Controller{
		UserController:         newUserController,
		VerificationController: newVerificationController,
		AuthController:         newAuthController,
		ServiceController:      newServiceController,
		ApplicationController:  newApplicationController,
	}
}
