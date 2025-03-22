package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core"
)

type HandlerDependency struct {
	Controller *Controller
}

type Handler struct {
}

func NewHandler(dependency HandlerDependency) *Handler {
	_ = dependency.Controller

	return &Handler{}
}

type RegisterHandlerItf interface {
	Register(router fiber.Router) error
}

func GetRegisters(handler *Handler) []RegisterHandlerItf {
	return []RegisterHandlerItf{}
}

func RegisterHandler(router fiber.Router, handler *Handler) error {
	api := router.Group(core.API)
	for _, register := range GetRegisters(handler) {
		err := register.Register(api)
		if err != nil {
			return fmt.Errorf("failed to register handler: %w", err)
		}
	}

	return nil
}
