package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetHealthHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"message": "pong"})
	}
}
