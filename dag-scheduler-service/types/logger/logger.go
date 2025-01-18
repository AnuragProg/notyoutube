package logger

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/anuragprog/notyoutube/dag-scheduler-service/types/errors"
)

type BaseAPILog struct {
	// non api related
	LogType   string    `json:"log_type"`
	Timestamp time.Time `json:"timestamp"`

	// request related
	Method         string              `json:"method"`
	Endpoint       string              `json:"endpoint"`
	ClientIP       string              `json:"client_ip"`
	RequestID      string              `json:"request_id"`
	TraceID        string              `json:"trace_id"`
	QueryParams    map[string][]string `json:"query_params"`    // Query parameters (if any)
	RequestHeaders map[string][]string `json:"request_headers"` // Relevant request headers

	// miscellaneous
	Miscellaneous map[string]interface{} `json:"miscellaneous"`
}

type APIInfoLog struct {
	BaseAPILog

	// response related
	LatencyInMs int64 `json:"latency_in_ms"`
	StatusCode  int   `json:"status_code"`
}

type APIDebugLog struct {
	BaseAPILog

	Message string `json:"message"`
}

type APIWarningLog struct {
	BaseAPILog
}

type APIErrorSeverity string

const (
	API_ERROR_SEVERITY_CRITICAL APIErrorSeverity = "critical"
	API_ERROR_SEVERITY_MAJOR    APIErrorSeverity = "major"
	API_ERROR_SEVERITY_MINOR    APIErrorSeverity = "minor"
)

type APIErrorLog struct {
	BaseAPILog

	// response related
	LatencyInMs int64 `json:"latency_in_ms"`
	StatusCode  int   `json:"status_code"`

	Severity     APIErrorSeverity `json:"severity"`
	StackTrace   string           `json:"stack_trace"`
	ErrorMessage string           `json:"error_message"`
}

func newAPIBaseLog(c echo.Context, timestamp time.Time, miscellaneous map[string]interface{}) BaseAPILog {
	method := c.Request().Method
	endpoint := ""
	if c.Request().URL != nil {
		endpoint = c.Request().URL.RawPath
	}
	clientIp := c.Request().RemoteAddr
	requestId := c.Response().Header().Get("x-request-id")
	traceId := c.Response().Header().Get("x-trace-id")

	queryParams := c.QueryParams()
	requestHeaders := c.Request().Header

	return BaseAPILog{
		LogType:   "api",
		Timestamp: timestamp,

		Method:         method,
		Endpoint:       endpoint,
		ClientIP:       clientIp,
		RequestID:      requestId,
		TraceID:        traceId,
		QueryParams:    queryParams,
		RequestHeaders: requestHeaders,

		Miscellaneous: miscellaneous,
	}
}

// This function should be called within the handler function. After creating the APIInfoLog object, you can safely use it elsewhere in your application.
func NewAPIInfoLog(c echo.Context, timestamp time.Time, latency time.Duration, miscellaneous map[string]interface{}) APIInfoLog {
	return APIInfoLog{
		BaseAPILog: newAPIBaseLog(c, timestamp, miscellaneous),

		LatencyInMs: latency.Milliseconds(),
		StatusCode:  c.Response().Status,
	}
}

// This function should be called within the handler function. Once created, the APIDebugLog object can be used safely in other parts of your application.
func NewAPIDebugLog(c echo.Context, timestamp time.Time, message string, miscellaneous map[string]interface{}) APIDebugLog {
	return APIDebugLog{
		BaseAPILog: newAPIBaseLog(c, timestamp, miscellaneous),

		Message: message,
	}
}

// This function should be called within the handler function. After creating the APIWarningLog object, you can safely use it elsewhere in your application.
func NewAPIWarningLog(c echo.Context, timestamp time.Time, miscellaneous map[string]interface{}) APIWarningLog {
	return APIWarningLog{
		BaseAPILog: newAPIBaseLog(c, timestamp, miscellaneous),
	}
}

// This function should be called within the handler function. Once created, the APIErrorLog object can be used safely in other parts of your application.
func NewAPIErrorLog(c echo.Context, timestamp time.Time, latency time.Duration, apiError errors.APIError, severity APIErrorSeverity, miscellaneous map[string]interface{}) APIErrorLog {
	return APIErrorLog{
		BaseAPILog: newAPIBaseLog(c, timestamp, miscellaneous),

		LatencyInMs: latency.Milliseconds(),
		StatusCode:  c.Response().Status,

		Severity:     severity,
		StackTrace:   apiError.StackTrace,
		ErrorMessage: apiError.Message,
	}
}
