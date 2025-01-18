package middlewares

import "github.com/labstack/echo/v4"

func GetNotImplementedMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return echo.ErrNotImplemented
		}
	}
}
