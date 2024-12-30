package configs

import (
	"time"

	"github.com/anuragprog/notyoutube/preprocessor-service/utils"
)


var (
	API_PORT     = utils.GetEnvIntDefault("API_PORT", 3001)
	SERVICE_NAME = utils.GetEnvDefault("SERVICE_NAME", "preprocessor-service") 

	ENVIRONMENT            = utils.GetEnvironment(utils.DEVELOPMENT_ENV)
	DEFAULT_PAGE_START int = 0
	DEFAULT_PAGE_SIZE  int = 25
	DEFAULT_TIMEOUT        = time.Second * 5

	// To use noop external services 
	USE_NOOP_DB = utils.GetEnvBoolDefault("USE_NOOP_DB", false)
	USE_NOOP_MQ = utils.GetEnvBoolDefault("USE_NOOP_MQ", false)

	KAFKA_BROKERS = utils.GetEnvListDefault("KAFKA_BROKERS", ",")
)
