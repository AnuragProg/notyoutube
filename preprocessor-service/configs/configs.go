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

	RAW_VIDEO_SERVICE_URL = utils.GetEnvDefault("RAW_VIDEO_SERVICE_URL", "localhost:50051")

	POSTGRES_HOST     = utils.GetEnvDefault("POSTGRES_HOST", "localhost")
	POSTGRES_PORT     = utils.GetEnvDefault("POSTGRES_PORT", "5432")
	POSTGRES_USER     = utils.GetEnvDefault("POSTGRES_USER", "root")
	POSTGRES_PASSWORD = utils.GetEnvDefault("POSTGRES_PASSWORD", "root")
	POSTGRES_DBNAME   = utils.GetEnvDefault("POSTGRES_DBNAME", "not_youtube")

	KAFKA_BROKERS      = utils.GetEnvListDefault("KAFKA_BROKERS", ",")
	MQ_TOPIC_RAW_VIDEO = utils.GetEnvDefault("MQ_TOPIC_RAW_VIDEO ", "raw-video")
	MQ_TOPIC_DAG       = utils.GetEnvDefault("MQ_TOPIC_DAG", "dag")
)
