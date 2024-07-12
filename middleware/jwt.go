package middleware

import (
	"BookAuthor_ManyToMany/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"Unauthorized": "Missing or malformed JWT",
			})
		}

		token, err := utils.ParseToken(tokenString)
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"Unauthorized": "Invalid or expired JWT",
			})
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Locals("usersId", claims["userId"])
		return c.Next()
	}
}
