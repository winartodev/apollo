package routes

import (
	"github.com/winartodev/apollo/core/configs"
	applicationController "github.com/winartodev/apollo/modules/application/controllers"
	authController "github.com/winartodev/apollo/modules/auth/controllers"
	guardianController "github.com/winartodev/apollo/modules/guardian/controllers"
	userController "github.com/winartodev/apollo/modules/user/controllers"
)

type ControllerDependency struct {
	SMTPClient *configs.SMTPClient
	Twilio     *configs.TwilioClient
	Repository *Repository
}

type Controller struct {
	UserController         userController.UserControllerItf
	VerificationController authController.VerificationControllerItf
	AuthController         authController.AuthControllerItf
	GuardianController     guardianController.GuardianControllerItf
}

func NewController(dependency ControllerDependency) *Controller {
	repository := dependency.Repository

	newUserController := userController.NewUserController(userController.UserController{
		UserRepository: repository.UserRepository,
	})

	newVerificationController := authController.NewVerificationController(authController.VerificationController{
		SmtpClient:             dependency.SMTPClient,
		TwilioClient:           dependency.Twilio,
		VerificationRepository: repository.VerificationRepository,
	})

	newAuthController := authController.NewAuthController(authController.AuthController{
		VerificationController: newVerificationController,
		UserController:         newUserController,
	})

	newGuardianController := guardianController.NewGuardianController(guardianController.GuardianController{
		ApplicationController: applicationController.NewApplicationController(applicationController.ApplicationController{
			UserController:                   newUserController,
			UserApplicationRepository:        repository.UserApplicationRepository,
			UserApplicationServiceRepository: repository.UserApplicationServiceRepository,
		}),
	})

	return &Controller{
		UserController:         newUserController,
		VerificationController: newVerificationController,
		AuthController:         newAuthController,
		GuardianController:     newGuardianController,
	}
}
