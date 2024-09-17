package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/anuragprog/notyoutube/file-service/utils/errors"
	"github.com/anuragprog/notyoutube/file-service/utils/log"
)

func GetLoggerMiddleware(appLogger log.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {

		requestStart := time.Now()
		err := c.Next()
		latency := time.Since(requestStart)

		switch err := err.(type) {

		case errors.APIError:
			severity := log.API_ERROR_SEVERITY_MINOR
			if err.StatusCode >= 500 {
				severity = log.API_ERROR_SEVERITY_MAJOR
			}
			errorLog := log.NewAPIErrorLog(
				c,
				requestStart,
				latency,
				err,
				severity,
				map[string]interface{}{},
			)
			go appLogger.LogAPIError(errorLog)

		case *fiber.Error:
			severity := log.API_ERROR_SEVERITY_MINOR
			if err.Code >= 500 {
				severity = log.API_ERROR_SEVERITY_MAJOR
			}
			apiErr := errors.IntoAPIError(err, err.Code, err.Message)
			errorLog := log.NewAPIErrorLog(
				c,
				requestStart,
				latency,
				apiErr,
				severity,
				map[string]interface{}{},
			)
			go appLogger.LogAPIError(errorLog)

		case nil:
			apiLog := log.NewAPIInfoLog(
				c,
				requestStart,
				latency,
				map[string]interface{}{},
			)
			go appLogger.LogAPIInfo(apiLog)

		default:

		}

		return err
	}
}
