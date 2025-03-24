package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core"
	authHandler "github.com/winartodev/apollo/modules/auth/handlers"
	userHandler "github.com/winartodev/apollo/modules/user/handlers"
	"time"
)

type HandlerDependency struct {
	Controller *Controller
}

type Handler struct {
	AuthHandler authHandler.AuthHandler
	UserHandler userHandler.UserHandler
}

func NewHandler(dependency HandlerDependency) *Handler {
	controller := dependency.Controller

	newAuthHandler := authHandler.NewAuthHandler(authHandler.AuthHandler{
		AuthController: controller.AuthController,
	})

	newUserHandler := userHandler.NewUserHandler(userHandler.UserHandler{
		UserController: controller.UserController,
	})

	return &Handler{
		AuthHandler: newAuthHandler,
		UserHandler: newUserHandler,
	}
}

type RegisterHandlerItf interface {
	Register(router fiber.Router) error
}

func GetRegisters(handler *Handler) []RegisterHandlerItf {
	return []RegisterHandlerItf{
		&handler.AuthHandler,
		&handler.UserHandler,
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
