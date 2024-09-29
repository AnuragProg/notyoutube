package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/anuragprog/notyoutube/file-service/types/errors"
	loggerType "github.com/anuragprog/notyoutube/file-service/types/logger"
	loggerRepo "github.com/anuragprog/notyoutube/file-service/repository/logger"
)

func GetLoggerMiddleware(appLogger loggerRepo.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {

		requestStart := time.Now()
		err := c.Next()
		latency := time.Since(requestStart)

		switch err := err.(type) {

		case errors.APIError:
			severity := loggerType.API_ERROR_SEVERITY_MINOR
			if err.StatusCode >= 500 {
				severity = loggerType.API_ERROR_SEVERITY_MAJOR
			}
			errorLog := loggerType.NewAPIErrorLog(
				c,
				requestStart,
				latency,
				err,
				severity,
				map[string]interface{}{},
			)
			go appLogger.LogAPIError(errorLog)

		case *fiber.Error:
			severity := loggerType.API_ERROR_SEVERITY_MINOR
			if err.Code >= 500 {
				severity = loggerType.API_ERROR_SEVERITY_MAJOR
			}
			apiErr := errors.IntoAPIError(err, err.Code, err.Message)
			errorLog := loggerType.NewAPIErrorLog(
				c,
				requestStart,
				latency,
				apiErr,
				severity,
				map[string]interface{}{},
			)
			go appLogger.LogAPIError(errorLog)

		case nil:
			apiLog := loggerType.NewAPIInfoLog(
				c,
				requestStart,
				latency,
				map[string]interface{}{},
			)
			go appLogger.LogAPIInfo(apiLog)

		}

		return err
	}
}
