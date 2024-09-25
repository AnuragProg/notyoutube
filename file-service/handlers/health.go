package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func GetHealthHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"message": "pong"})
	}
}
