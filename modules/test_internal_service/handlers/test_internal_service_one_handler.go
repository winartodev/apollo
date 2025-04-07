package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core"
	"github.com/winartodev/apollo/core/middlewares"
	"github.com/winartodev/apollo/core/responses"
	enums2 "github.com/winartodev/apollo/modules/application/enums"
)

var (
	middlewareHandleInternalAccess1 = &middlewares.InternalAccessConfig{
		Application:        enums2.ApolloInternal,
		ApplicationService: enums2.TestInternalServices1,
	}

	middlewareHandleInternalAccess2 = &middlewares.InternalAccessConfig{
		Application:        enums2.ApolloInternal,
		ApplicationService: enums2.TestInternalServices2,
	}

	middlewareHandleInternalAccess3 = &middlewares.InternalAccessConfig{
		Application:        enums2.ApolloInternal,
		ApplicationService: enums2.TestInternalServices3,
	}
)

type TestInternalServiceHandler struct {
	middlewares.Middleware
}

func NewTestInternalServiceHandler(handler TestInternalServiceHandler) TestInternalServiceHandler {
	return TestInternalServiceHandler{
		handler.Middleware,
	}
}

func (h *TestInternalServiceHandler) HandleRequest(ctx *fiber.Ctx) error {
	return responses.SuccessResponse(ctx, fiber.StatusOK, "Success", "ok", nil)
}

func (h *TestInternalServiceHandler) Register(router fiber.Router) error {

	v1 := router.Group(core.V1)
	internal := v1.Group(core.AccessInternal)
	internalServiceOne := internal.Group("/test-internal-service-1", h.HandleInternalAccess(middlewareHandleInternalAccess1))
	internalServiceOne.Get("/", h.HandleRequest)
	internalServiceOne.Post("/", h.HandleRequest)
	internalServiceOne.Put("/", h.HandleRequest)
	internalServiceOne.Delete("/", h.HandleRequest)

	internalServiceTwo := internal.Group("/test-internal-service-2", h.HandleInternalAccess(middlewareHandleInternalAccess2))
	internalServiceTwo.Get("/", h.HandleRequest)
	internalServiceTwo.Post("/", h.HandleRequest)
	internalServiceTwo.Put("/", h.HandleRequest)
	internalServiceTwo.Delete("/", h.HandleRequest)

	internalServiceThree := internal.Group("/test-internal-service-3", h.HandleInternalAccess(middlewareHandleInternalAccess3))
	internalServiceThree.Get("/", h.HandleRequest)
	internalServiceThree.Post("/", h.HandleRequest)
	internalServiceThree.Put("/", h.HandleRequest)
	internalServiceThree.Delete("/", h.HandleRequest)
	return nil
}
