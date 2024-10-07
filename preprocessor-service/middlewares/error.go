package middlewares

import (
	errType "github.com/anuragprog/notyoutube/preprocessor-service/types/errors"
	"github.com/gofiber/fiber/v2"
)

// will set the error status codes and messages according to the error caught
func GetErrorResponseHandlerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {

		err = c.Next()
		switch err := err.(type) {
		case errType.APIError:
			c.Status(err.StatusCode).JSON(fiber.Map{"message": err.Message})
		case *fiber.Error:
			c.Status(err.Code).JSON(fiber.Map{"message": err.Message})
		}

		return err
	}
}
