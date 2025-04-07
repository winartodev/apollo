package middlewares

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core/helpers"
	"github.com/winartodev/apollo/core/responses"
	applicationEnum "github.com/winartodev/apollo/modules/application/enums"
	guardianController "github.com/winartodev/apollo/modules/guardian/controllers"
	userController "github.com/winartodev/apollo/modules/user/controllers"
	"strings"
)

const (
	protected = "protected"
	internal  = "internal"
)

var (
	errorInvalidPublicAccess   = errors.New("public resource cannot be accessed due to invalid request")
	errorInvalidInternalAccess = errors.New("internal resource cannot be accessed due to invalid request")
	errorInvalidToken          = errors.New("provided token is invalid or expired")
	errorFailedInstanceJWT     = errors.New("failed to create instance JWT")
	errorMissingToken          = errors.New("authentication token is missing or improperly formatted. Expected 'Bearer <token>'")
	errorUserNotFound          = errors.New("user not found")
)

type Middleware struct {
	UserController     userController.UserControllerItf
	GuardianController guardianController.GuardianControllerItf
}

func (m *Middleware) HandlePublicAccess() fiber.Handler {
	return func(c *fiber.Ctx) error {
		access := getAccessFromPath(c)
		if access == internal || access == protected {
			return responses.FailedResponse(c, fiber.StatusForbidden, "Access Denied", errorInvalidPublicAccess)
		}

		var token string
		if isAuthHeaderExists(c, &token) {
			claim, err := verifyAuthHeader(token)
			if err != nil {
				return responses.FailedResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
			}

			c.Locals("id", claim.ID)
			c.Locals("username", claim.Username)
			c.Locals("email", claim.Email)
		}

		return c.Next()
	}
}

type InternalAccessConfig struct {
	Application        applicationEnum.ApplicationEnum
	ApplicationService applicationEnum.ApplicationServiceEnum
}

func (m *Middleware) HandleInternalAccess(config *InternalAccessConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		access := getAccessFromPath(c)
		if access != internal {
			return responses.FailedResponse(c, fiber.StatusForbidden, "Access Denied", errorInvalidInternalAccess)
		}

		var token string
		if isAuthHeaderExists(c, &token) {
			claim, err := verifyAuthHeader(token)
			if err != nil {
				return responses.FailedResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
			}

			context := c.Context()
			user, err := m.UserController.GetUserByID(context, claim.ID)
			if err != nil {
				return responses.FailedResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
			}

			if user == nil {
				return responses.FailedResponse(c, fiber.StatusUnauthorized, "Unauthorized", errorUserNotFound)
			}

			if config != nil {
				permissionGranted, err := m.GuardianController.CheckUserPermissionToInternalApp(context, user.ID, config.Application, config.ApplicationService, c.Method())
				if err != nil {
					return responses.FailedResponse(c, fiber.StatusInternalServerError, "Unauthorized", err)
				}

				if !permissionGranted {
					return responses.FailedResponse(c, fiber.StatusForbidden, "Unauthorized", errorInvalidInternalAccess)
				}
			}

			c.Locals("id", claim.ID)
			c.Locals("username", claim.Username)
			c.Locals("email", claim.Email)
		} else {
			return responses.FailedResponse(c, fiber.StatusForbidden, "Access Denied", errorMissingToken)
		}

		return c.Next()
	}
}

func isAuthHeaderExists(c *fiber.Ctx, token *string) bool {
	authHeader := c.Get("Authorization")

	if len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
		*token = authHeader[7:]
		return true
	}

	return false
}

func getAccessFromPath(ctx *fiber.Ctx) string {
	paths := strings.Split(ctx.Path(), "/")
	if len(paths) > 3 {
		return paths[3]
	}

	return ""
}

func verifyAuthHeader(token string) (res *helpers.JWTClaims, err error) {
	jwt, err := helpers.NewJWT()
	if err != nil {
		return nil, errorFailedInstanceJWT
	}

	claims, isValid, err := jwt.VerifyToken(jwt.AccessToken.SecretKey, token)
	if err != nil {
		return nil, err
	}

	if !isValid {
		return nil, errorInvalidToken
	}

	claimByte, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}

	var claim helpers.JWTClaims
	err = json.Unmarshal(claimByte, &claim)
	if err != nil {
		return nil, err
	}

	return &claim, nil
}
