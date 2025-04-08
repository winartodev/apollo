package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core"
	"github.com/winartodev/apollo/core/helpers"
	"github.com/winartodev/apollo/core/middlewares"
	"github.com/winartodev/apollo/core/responses"
	guardianController "github.com/winartodev/apollo/modules/guardian/controllers"
	guardianEntity "github.com/winartodev/apollo/modules/guardian/entities"
	"strconv"
)

type GuardianRoleHandler struct {
	middlewares.Middleware
	GuardianRole guardianController.GuardianRoleControllerItf
}

func NewGuardianRoleHandler(handler GuardianRoleHandler) GuardianRoleHandler {
	return GuardianRoleHandler{
		Middleware:   handler.Middleware,
		GuardianRole: handler.GuardianRole,
	}
}

func (h *GuardianRoleHandler) GetRoles(c *fiber.Ctx) error {
	ctx := c.Context()
	appID, err := helpers.GetAppIDFromContext(c)
	if err != nil {
		return responses.FailedResponse(
			c,
			fiber.StatusBadRequest,
			"Failed retrieve roles data",
			errors.New("application id not provided"),
		)
	}

	slug := c.Query("slug")
	if slug != "" {
		data, err := h.GuardianRole.GetRoleBySlug(ctx, appID, slug)
		if err != nil {
			return responses.FailedResponse(
				c,
				fiber.StatusBadRequest,
				"Failed retrieve role",
				errors.New("application id not provided"),
			)
		}

		return responses.SuccessResponse(
			c,
			fiber.StatusCreated,
			"Success retrieve role data",
			data,
			nil,
		)
	}

	data, err := h.GuardianRole.GetRoles(ctx, appID)
	if err != nil {
		return responses.FailedResponse(
			c,
			fiber.StatusInternalServerError,
			"Failed retrieve roles data",
			err,
		)
	}

	return responses.SuccessResponse(
		c,
		fiber.StatusOK,
		"Success retrieve roles data",
		data,
		nil,
	)
}

func (h *GuardianRoleHandler) AddRole(c *fiber.Ctx) error {
	ctx := c.Context()
	appID, err := helpers.GetAppIDFromContext(c)
	if err != nil {
		return responses.FailedResponse(
			c,
			fiber.StatusBadRequest,
			"Failed create role",
			errors.New("application id not provided"),
		)
	}

	req := guardianEntity.GuardianRole{}
	err = c.BodyParser(&req)
	if err != nil {
		return responses.FailedResponse(
			c,
			fiber.StatusBadRequest,
			"Failed create role",
			err,
		)
	}

	req.ApplicationID = appID
	res, err := h.GuardianRole.CreateRole(ctx, req)
	if err != nil {
		return responses.FailedResponse(
			c,
			fiber.StatusInternalServerError,
			"Failed create role",
			err,
		)
	}

	return responses.SuccessResponse(
		c,
		fiber.StatusCreated,
		"Success create role",
		res,
		nil,
	)
}

func (h *GuardianRoleHandler) GetRole(c *fiber.Ctx) error {
	ctx := c.Context()
	appID, err := helpers.GetAppIDFromContext(c)
	if err != nil {
		return responses.FailedResponse(
			c,
			fiber.StatusBadRequest,
			"Failed retrieve role data",
			errors.New("application id not provided"),
		)
	}

	paramID := c.Params("id")
	if paramID == "" {
		return responses.FailedResponse(
			c,
			fiber.StatusCreated,
			"Missing identifier",
			errors.New("role ID must be provided"),
		)
	}

	id, err := strconv.ParseInt(paramID, 10, 64)
	if err != nil {
		return responses.FailedResponse(
			c,
			fiber.StatusBadRequest,
			"Failed retrieve role data",
			err,
		)
	}

	data, err := h.GuardianRole.GetRoleByID(ctx, appID, id)
	if err != nil {
		return responses.FailedResponse(
			c,
			fiber.StatusBadRequest,
			"Failed retrieve role data",
			errors.New("application id not provided"),
		)
	}

	return responses.SuccessResponse(
		c,
		fiber.StatusCreated,
		"Success retrieve role data",
		data,
		nil,
	)

}

func (h *GuardianRoleHandler) Register(router fiber.Router) error {
	v1 := router.Group(core.V1)
	internal := v1.Group(core.AccessInternal, h.HandleInternalAccess(nil))
	guardian := internal.Group("/guardian")
	guardian.Get("/roles", h.GetRoles)
	guardian.Get("/roles/:id", h.GetRole)
	guardian.Post("/roles", h.AddRole)

	return nil
}
