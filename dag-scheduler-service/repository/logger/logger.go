package logger

import (
	"io"

	"github.com/anuragprog/notyoutube/dag-scheduler-service/types/logger"
)


type Logger interface {

	io.Closer

	LogAPIInfo(apiLog logger.APIInfoLog) 
	LogAPIDebug(apiLog logger.APIDebugLog) 
	LogAPIWarning(apiLog logger.APIWarningLog) 
	LogAPIError(apiLog logger.APIErrorLog) 

}
