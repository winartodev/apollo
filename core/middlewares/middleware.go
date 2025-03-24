package middlewares

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core/helpers"
	"github.com/winartodev/apollo/core/responses"
	"strings"
)

const (
	protected = "protected"
	internal  = "internal"
)

var (
	errorInvalidPublicAccess = errors.New("public resource cannot be accessed due to invalid request")
	errorInvalidToken        = errors.New("provided token is invalid or expired")
	errorFailedInstanceJWT   = errors.New("failed to create instance JWT")
)

func HandlePublicAccess() fiber.Handler {
	return func(c *fiber.Ctx) error {
		access := getAccessFromPath(c)
		if access == internal || access == protected {
			return responses.FailedResponse(c, fiber.StatusForbidden, "Access Denied", errorInvalidPublicAccess)
		}

		var token string
		if isAuthHeaderExists(c, &token) {
			jwt, err := helpers.NewJWT()
			if err != nil {
				return responses.FailedResponse(c, fiber.StatusInternalServerError, "Access Denied", errorFailedInstanceJWT)
			}

			claims, isValid, err := jwt.VerifyToken(jwt.AccessToken.SecretKey, token)
			if err != nil {
				return responses.FailedResponse(c, fiber.StatusInternalServerError, "Authentication Failed", err)
			}

			if !isValid {
				return responses.FailedResponse(c, fiber.StatusInternalServerError, "Authentication Failed", errorInvalidToken)
			}

			if id, ok := claims["id"].(float64); ok {
				c.Locals("id", id)
			}

			if username, ok := claims["username"].(string); ok {
				c.Locals("username", username)
			}

			if email, ok := claims["email"].(string); ok {
				c.Locals("email", email)
			}
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
