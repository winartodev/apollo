package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core"
	"github.com/winartodev/apollo/core/middlewares"
	authHandler "github.com/winartodev/apollo/modules/auth/handlers"
	internalServiceTest "github.com/winartodev/apollo/modules/test_internal_service/handlers"
	userHandler "github.com/winartodev/apollo/modules/user/handlers"
	"time"
)

type HandlerDependency struct {
	Controller *Controller
}

type Handler struct {
	AuthHandler         authHandler.AuthHandler
	UserHandler         userHandler.UserHandler
	InternalServiceTest internalServiceTest.TestInternalServiceHandler
}

func NewHandler(dependency HandlerDependency) *Handler {
	controller := dependency.Controller

	middleware := middlewares.Middleware{
		UserController:     controller.UserController,
		GuardianController: controller.GuardianController,
	}

	newAuthHandler := authHandler.NewAuthHandler(authHandler.AuthHandler{
		Middleware:             middleware,
		VerificationController: controller.VerificationController,
		AuthController:         controller.AuthController,
	})

	newUserHandler := userHandler.NewUserHandler(userHandler.UserHandler{
		Middleware:     middleware,
		UserController: controller.UserController,
	})

	newTest := internalServiceTest.NewTestInternalServiceHandler(internalServiceTest.TestInternalServiceHandler{
		Middleware: middleware,
	})

	return &Handler{
		AuthHandler:         newAuthHandler,
		UserHandler:         newUserHandler,
		InternalServiceTest: newTest,
	}
}

type RegisterHandlerItf interface {
	Register(router fiber.Router) error
}

func GetRegisters(handler *Handler) []RegisterHandlerItf {
	return []RegisterHandlerItf{
		&handler.AuthHandler,
		&handler.UserHandler,
		&handler.InternalServiceTest,
	}
}

func RegisterHandler(router fiber.Router, handler *Handler) error {
	api := router.Group(core.API)

	for _, register := range GetRegisters(handler) {
		err := register.Register(api)
		if err != nil {
			return fmt.Errorf("failed to register handler: %w", err)
		}
	}

	healthHandler := HealthHandler{
		StartTime: time.Now(),
	}

	api.Get("/healthz", healthHandler.HealthZ)

	return nil
}
