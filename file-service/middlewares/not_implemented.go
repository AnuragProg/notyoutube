package middlewares

import "github.com/gofiber/fiber/v2"

func GetNotImplementedMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}
