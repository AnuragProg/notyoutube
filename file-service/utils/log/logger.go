package log

import "io"

type Logger interface {
	io.Closer

	LogAPIDebug(apiLog APIDebugLog)
	LogAPIError(apiLog APIErrorLog)
	LogAPIInfo(apiLog APIInfoLog)
	LogAPIWarning(apiLog APIWarningLog)
}
