package middlewares

import (
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	errorType "github.com/anuragprog/notyoutube/file-service/types/errors"
	loggerType "github.com/anuragprog/notyoutube/file-service/types/logger"
	loggerRepo "github.com/anuragprog/notyoutube/file-service/repository/logger"
)

func GetRecoverMiddleware(appLogger loggerRepo.Logger) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {

		requestStart := time.Now()
		err = c.Next()
		latency := time.Since(requestStart)

		// Catch panics
		defer func() {
			if r := recover(); r != nil {

				// log the as critical error
				severity := loggerType.API_ERROR_SEVERITY_CRITICAL
				apiErr := errorType.IntoAPIError(errors.New("panic occurred"), http.StatusInternalServerError, "panic occurred")
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

		// Return err if exist, else move to next handler
		return err
	}
}
