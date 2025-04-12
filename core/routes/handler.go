package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core"
	"github.com/winartodev/apollo/core/middlewares"
	applicationHandler "github.com/winartodev/apollo/modules/application/handlers"
	authHandler "github.com/winartodev/apollo/modules/auth/handlers"
	userHandler "github.com/winartodev/apollo/modules/user/handlers"
	"time"
)

type HandlerDependency struct {
	Controller *Controller
}

type Handler struct {
	AuthHandler               authHandler.AuthHandler
	AuthApolloInternalHandler authHandler.AuthApolloInternalHandler
	UserHandler               userHandler.UserHandler
	ApplicationHandler        applicationHandler.ApplicationHandler
	ApplicationServiceHandler applicationHandler.ApplicationServiceHandler
}

func NewHandler(dependency HandlerDependency) *Handler {
	controller := dependency.Controller

	middleware := middlewares.Middleware{
		UserController: controller.UserController,
	}

	newAuthHandler := authHandler.NewAuthHandler(authHandler.AuthHandler{
		Middleware:             middleware,
		VerificationController: controller.VerificationController,
		AuthController:         controller.AuthController,
	})

	newAuthApolloInternalHandler := authHandler.NewAuthApolloInternalHandler(authHandler.AuthApolloInternalHandler{
		Middleware:            middleware,
		AuthController:        controller.AuthController,
		ApplicationController: controller.ApplicationController,
		UserController:        controller.UserController,
	})

	newUserHandler := userHandler.NewUserHandler(userHandler.UserHandler{
		Middleware:     middleware,
		UserController: controller.UserController,
	})

	newApplicationHandler := applicationHandler.NewApplicationHandler(applicationHandler.ApplicationHandler{
		Middleware:            middleware,
		ApplicationController: controller.ApplicationController,
		ServiceController:     controller.ServiceController,
	})

	newApolloApplicationServiceHandler := applicationHandler.NewApplicationServiceHandler(applicationHandler.ApplicationServiceHandler{
		Middleware:            middleware,
		ApplicationController: controller.ApplicationController,
		ServiceController:     controller.ServiceController,
	})

	return &Handler{
		AuthHandler:               newAuthHandler,
		AuthApolloInternalHandler: newAuthApolloInternalHandler,
		UserHandler:               newUserHandler,
		ApplicationHandler:        newApplicationHandler,
		ApplicationServiceHandler: newApolloApplicationServiceHandler,
	}
}

type RegisterHandlerItf interface {
	Register(router fiber.Router) error
}

func GetRegisters(handler *Handler) []RegisterHandlerItf {
	return []RegisterHandlerItf{
		&handler.AuthHandler,
		&handler.AuthApolloInternalHandler,
		&handler.UserHandler,
		&handler.ApplicationHandler,
		&handler.ApplicationServiceHandler,
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
