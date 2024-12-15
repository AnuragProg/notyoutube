package middlewares

import (
	errType "github.com/anuragprog/notyoutube/file-service/types/errors"
	"github.com/labstack/echo/v4"
)

// will set the error status codes and messages according to the error caught
func GetErrorResponseHandlerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {

			err = next(c)
			switch err := err.(type) {
			case errType.APIError:
				c.JSON(err.StatusCode, echo.Map{"message": err.Message})
			case *echo.HTTPError:
				c.JSON(err.Code, echo.Map{"message": err.Message})
			}

			return err
		}

	}
}
