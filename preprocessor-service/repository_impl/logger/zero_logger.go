package logger

import (
	"encoding/json"
	"io"

	"github.com/rs/zerolog"
	"github.com/anuragprog/notyoutube/preprocessor-service/types/logger"
)

type ZeroLogger struct {
	writer         io.WriteCloser
	internalLogger zerolog.Logger
}

func NewZeroLogger(writer io.WriteCloser, serviceName, environment string) *ZeroLogger {
	zeroLogger := zerolog.New(writer).
		With().
		Str("service_name", serviceName).
		Str("environment", environment).
		Logger()

	return &ZeroLogger{
		internalLogger: zeroLogger,
	}
}

// Utility function to append/prepend base log data to the logger
func adjustBaseAPILogInZeroLogEvent(event *zerolog.Event, apiLog logger.BaseAPILog) *zerolog.Event {
	queryParamsJson, _ := json.Marshal(apiLog.QueryParams)
	requestHeadersJson, _ := json.Marshal(apiLog.RequestHeaders)
	miscellaneousJson, _ := json.Marshal(apiLog.Miscellaneous)

	return event.
		Str("log_type", apiLog.LogType).
		Int64("timestamp", apiLog.Timestamp.Unix()).
		Str("method", apiLog.Method).
		Str("endpoint", apiLog.Endpoint).
		Str("client_ip", apiLog.ClientIP).
		Str("request_id", apiLog.RequestID).
		Str("trace_id", apiLog.TraceID).
		RawJSON("query_params", queryParamsJson).
		RawJSON("request_headers", requestHeadersJson).
		RawJSON("miscellaneous", miscellaneousJson)
}

func (zeroLogger *ZeroLogger) LogAPIInfo(apiLog logger.APIInfoLog) {
	adjustBaseAPILogInZeroLogEvent(
		zeroLogger.internalLogger.Info(),
		apiLog.BaseAPILog,
	).
		Int64("latency_in_ms", apiLog.LatencyInMs).
		Int("status_code", apiLog.StatusCode).
		Send()
}

func (zeroLogger *ZeroLogger) LogAPIDebug(apiLog logger.APIDebugLog) {
	adjustBaseAPILogInZeroLogEvent(
		zeroLogger.internalLogger.Debug(),
		apiLog.BaseAPILog,
	).
		Str("message", apiLog.Message).
		Send()
}

func (zeroLogger *ZeroLogger) LogAPIWarning(apiLog logger.APIWarningLog) {
	adjustBaseAPILogInZeroLogEvent(
		zeroLogger.internalLogger.Warn(),
		apiLog.BaseAPILog,
	).
		Send()
}

func (zeroLogger *ZeroLogger) LogAPIError(apiLog logger.APIErrorLog) {
	adjustBaseAPILogInZeroLogEvent(
		zeroLogger.internalLogger.Error(),
		apiLog.BaseAPILog,
	).
		Int64("latency_in_ms", apiLog.LatencyInMs).
		Int("status_code", apiLog.StatusCode).
		Str("severity", string(apiLog.Severity)).
		Str("stack_trace", apiLog.StackTrace).
		Str("error_message", apiLog.ErrorMessage).
		Send()
}

func (zeroLogger *ZeroLogger) Close() error {
	return zeroLogger.writer.Close()
}
