package logger

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/anuragprog/notyoutube/file-service/types/errors"
)

type BaseAPILog struct {
	// non api related
	LogType   string    `json:"log_type"`
	Timestamp time.Time `json:"timestamp"`

	// request related
	Method         string              `json:"method"`
	Endpoint       string              `json:"endpoint"`
	ClientIP       string              `json:"client_ip"`
	UserAgent      string              `json:"user_agent"`
	RequestID      string              `json:"request_id"`
	QueryParams    map[string]string   `json:"query_params"`    // Query parameters (if any)
	RequestHeaders map[string][]string `json:"request_headers"` // Relevant request headers
	RequestBody    string              `json:"request_body"`    // Request body (if applicable)

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

	// response related
	// LatencyInMs int64 `json:"latency_in_ms"`
	// StatusCode  int   `json:"status_code"`
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

func getRequestBody(c *fiber.Ctx) []byte {

	// in case of form data just log the key value pairs and a place holder for indicating present file
	form, err := c.MultipartForm()
	if err != nil {
		formValue := map[string]interface{}{}
		for key, value := range form.Value {
			formValue[key] = value
		}
		for key, value := range form.File {
			formValue[key] = fmt.Sprintf("<%v files>", len(value))
		}
		formJson, _ := json.Marshal(formValue)
		return formJson
	}
	

	// parse body of other request
	requestBody := make([]byte, len(c.Request().Body()))
	copy(requestBody, c.Request().Body())
	return requestBody
}

func newAPIBaseLog(c *fiber.Ctx, timestamp time.Time, miscellaneous map[string]interface{}) BaseAPILog {
	method := string(c.Request().Header.Method())
	endpoint := strings.Clone(c.OriginalURL())
	clientIp := c.IP()
	requestId := strings.Clone(c.GetRespHeader("x-request-id", ""))
	userAgent := strings.Clone(c.Get("User-Agent", ""))

	queryParams := make(map[string]string)
	for key, value := range c.Queries() {
		copiedKey := strings.Clone(key)
		copiedValue := strings.Clone(value)
		queryParams[copiedKey] = copiedValue
	}
	requestHeaders := make(map[string][]string)
	for key, values := range c.GetReqHeaders() {
		copiedValues := []string{}
		for _, value := range values {
			if value == "" {
				continue
			}
			copiedValues = append(copiedValues, strings.Clone(value))
		}
		copiedKey := strings.Clone(key)
		requestHeaders[copiedKey] = copiedValues
	}

	requestBody := getRequestBody(c)

	return BaseAPILog{
		LogType:   "api",
		Timestamp: timestamp,

		Method:         method,
		Endpoint:       endpoint,
		ClientIP:       clientIp,
		UserAgent:      userAgent,
		RequestID:      requestId,
		QueryParams:    queryParams,
		RequestHeaders: requestHeaders,
		RequestBody:    string(requestBody),

		Miscellaneous: miscellaneous,
	}

}

// This function should be called within the handler function. After creating the APIInfoLog object, you can safely use it elsewhere in your application.
func NewAPIInfoLog(c *fiber.Ctx, timestamp time.Time, latency time.Duration, miscellaneous map[string]interface{}) APIInfoLog {
	return APIInfoLog{
		BaseAPILog: newAPIBaseLog(c, timestamp, miscellaneous),

		LatencyInMs: latency.Milliseconds(),
		StatusCode:  c.Response().StatusCode(),
	}
}

// This function should be called within the handler function. Once created, the APIDebugLog object can be used safely in other parts of your application.
func NewAPIDebugLog(c *fiber.Ctx, timestamp time.Time, message string, miscellaneous map[string]interface{}) APIDebugLog {
	return APIDebugLog{
		BaseAPILog: newAPIBaseLog(c, timestamp, miscellaneous),

		Message: message,
	}
}

// This function should be called within the handler function. After creating the APIWarningLog object, you can safely use it elsewhere in your application.
func NewAPIWarningLog(c *fiber.Ctx, timestamp time.Time, miscellaneous map[string]interface{}) APIWarningLog {
	return APIWarningLog{
		BaseAPILog: newAPIBaseLog(c, timestamp, miscellaneous),
	}
}

// This function should be called within the handler function. Once created, the APIErrorLog object can be used safely in other parts of your application.
func NewAPIErrorLog(c *fiber.Ctx, timestamp time.Time, latency time.Duration, apiError errors.APIError, severity APIErrorSeverity, miscellaneous map[string]interface{}) APIErrorLog {
	return APIErrorLog{
		BaseAPILog: newAPIBaseLog(c, timestamp, miscellaneous),

		LatencyInMs: latency.Milliseconds(),
		StatusCode:  c.Response().StatusCode(),

		Severity:     severity,
		StackTrace:   apiError.StackTrace,
		ErrorMessage: apiError.Message,
	}
}
