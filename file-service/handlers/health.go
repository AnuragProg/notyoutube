package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func HealthHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Status(http.StatusOK).JSON(map[string]interface{}{
			"message": "pong",
		})
		return nil
	}
}
