package middlewares

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
)

func GetRequestIdMiddleware() echo.MiddlewareFunc {
	return middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			generatedUUID, err := uuid.NewV7()
			if err != nil {
				panic(err)
			}
			return generatedUUID.String()
		},
		TargetHeader: "X-Request-Id",
	})
}
