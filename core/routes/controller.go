package routes

import (
	"github.com/winartodev/apollo/core/configs"
	applicationController "github.com/winartodev/apollo/modules/application/controllers"
	authController "github.com/winartodev/apollo/modules/auth/controllers"
	guardianController "github.com/winartodev/apollo/modules/guardian/controllers"
	userController "github.com/winartodev/apollo/modules/user/controllers"
)

type ControllerDependency struct {
	EnableOTPVerification bool
	SMTPClient            *configs.SMTPClient
	Twilio                *configs.TwilioClient
	Repository            *Repository
}

type Controller struct {
	AuthController         authController.AuthControllerItf
	GuardianController     guardianController.GuardianControllerItf
	UserController         userController.UserControllerItf
	VerificationController authController.VerificationControllerItf
}

func NewController(dependency ControllerDependency) *Controller {
	repository := dependency.Repository

	newUserController := userController.NewUserController(userController.UserController{
		UserRepository:            repository.UserRepository,
		UserRoleRepository:        repository.UserRoleRepository,
		UserApplicationRepository: repository.UserApplicationRepository,
	})

	newVerificationController := authController.NewVerificationController(authController.VerificationController{
		EnableOTPVerification:  dependency.EnableOTPVerification,
		SmtpClient:             dependency.SMTPClient,
		TwilioClient:           dependency.Twilio,
		VerificationRepository: repository.VerificationRepository,
	})

	newAuthController := authController.NewAuthController(authController.AuthController{
		VerificationController: newVerificationController,
		UserController:         newUserController,
	})

	newApplicationController := applicationController.NewApplicationController(applicationController.ApplicationController{
		ApplicationServiceRepository: repository.ApplicationServiceRepository,
	})

	newGuardianController := guardianController.NewGuardianController(guardianController.GuardianController{
		UserController:               newUserController,
		ApplicationController:        newApplicationController,
		GuardianPermissionRepository: repository.GuardianPermissionRepository,
	})

	return &Controller{
		UserController:         newUserController,
		VerificationController: newVerificationController,
		AuthController:         newAuthController,
		GuardianController:     newGuardianController,
	}
}
