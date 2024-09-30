package middlewares

import (
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	loggerRepo "github.com/anuragprog/notyoutube/file-service/repository/logger"
	errType "github.com/anuragprog/notyoutube/file-service/types/errors"
	loggerType "github.com/anuragprog/notyoutube/file-service/types/logger"
)

func GetLoggerMiddleware(appLogger loggerRepo.Logger) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {

		requestStart := time.Now()
		err = c.Next()
		latency := time.Since(requestStart)

		// Panic logger
		defer func() {
			if r := recover(); r != nil {
				// log the as critical error
				severity := loggerType.API_ERROR_SEVERITY_CRITICAL
				apiErr := errType.IntoAPIError(errors.New("panic occurred"), http.StatusInternalServerError, "panic occurred")
				errorLog := loggerType.NewAPIErrorLog(
					c,
					requestStart,
					latency,
					apiErr,
					severity,
					map[string]interface{}{},
				)
				go appLogger.LogAPIError(errorLog)

				err = fiber.NewError(http.StatusInternalServerError, "something went wrong")
			}
		}()

		switch err := err.(type) {
		case errType.APIError:
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
			apiErr := errType.IntoAPIError(err, err.Code, err.Message)
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

		// error handler middleware to be called before this and setup of status codes and messages 
		// to be done before logger middleware in order to propagate and extract correct logs
		err = nil
		return
	}
}
