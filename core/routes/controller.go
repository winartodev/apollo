package routes

import (
	"github.com/winartodev/apollo/core/configs"
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
}

func NewController(dependency ControllerDependency) *Controller {
	repository := dependency.Repository

	newUserController := userController.NewUserController(userController.UserController{
		UserRepository: repository.UserRepository,
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

	return &Controller{
		UserController:         newUserController,
		VerificationController: newVerificationController,
		AuthController:         newAuthController,
	}
}
