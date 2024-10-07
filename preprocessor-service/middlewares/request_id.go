package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
)

func GetRequestIdMiddleware() fiber.Handler {
	return requestid.New(requestid.Config{
		Header: "X-Request-Id", // although by default is X-Request-ID, but it is made canonical to X-Request-Id
		Generator: func() string {
			generatedUUID, err := uuid.NewV7()
			if err != nil {
				panic(err)
			}
			return generatedUUID.String()
		},
	})
}
