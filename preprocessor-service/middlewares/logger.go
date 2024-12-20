package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	loggerRepo "github.com/anuragprog/notyoutube/preprocessor-service/repository/logger"
	errType "github.com/anuragprog/notyoutube/preprocessor-service/types/errors"
	loggerType "github.com/anuragprog/notyoutube/preprocessor-service/types/logger"
)

func GetLoggerMiddleware(appLogger loggerRepo.Logger) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			requestStart := time.Now()
			err = next(c)
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

					err = echo.NewHTTPError(http.StatusInternalServerError, "something went wrong")
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

			case *echo.HTTPError:
				severity := loggerType.API_ERROR_SEVERITY_MINOR
				if err.Code >= 500 {
					severity = loggerType.API_ERROR_SEVERITY_MAJOR
				}
				apiErr := errType.IntoAPIError(err, err.Code, fmt.Sprintf("%v", err.Message))
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

}
