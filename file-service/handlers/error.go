package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/anuragprog/notyoutube/file-service/types/errors"
)

func GetErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {

		switch err := err.(type) {
		case errors.APIError:
			return c.Status(err.StatusCode).JSON(fiber.Map{"message": err.Message})
		case *fiber.Error:
			return c.Status(err.Code).JSON(fiber.Map{"message": err.Message})
		}
		return nil
	}
}
