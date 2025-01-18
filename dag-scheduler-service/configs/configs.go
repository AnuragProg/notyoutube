package configs

import (
	"time"

	"github.com/anuragprog/notyoutube/dag-scheduler-service/utils"
)

var (
	API_PORT     = utils.GetEnvIntDefault("API_PORT", 3002)
	GRPC_PORT    = utils.GetEnvIntDefault("GRPC_PORT", 50053)
	SERVICE_NAME = utils.GetEnvDefault("SERVICE_NAME", "dag-scheduler-service")

	ENVIRONMENT            = utils.GetEnvironment(utils.DEVELOPMENT_ENV)
	DEFAULT_PAGE_START int = 0
	DEFAULT_PAGE_SIZE  int = 25
	DEFAULT_TIMEOUT        = time.Second * 5

	// To use noop external services
	USE_NOOP_DB = utils.GetEnvBoolDefault("USE_NOOP_DB", false)
	USE_NOOP_MQ = utils.GetEnvBoolDefault("USE_NOOP_MQ", false)
	USE_NOOP_STORE = utils.GetEnvBoolDefault("USE_NOOP_STORE", false)

	RAW_VIDEO_SERVICE_URL = utils.GetEnvDefault("RAW_VIDEO_SERVICE_URL", "localhost:50051")

	MINIO_URI               = utils.GetEnvDefault("MINIO_URI", "localhost:9000")
	MINIO_SERVER_ACCESS_KEY = utils.GetEnvDefault("MINIO_SERVER_ACCESS_KEY", "minio-access-key")
	MINIO_SERVER_SECRET_KEY = utils.GetEnvDefault("MINIO_SERVER_SECRET_KEY", "minio-secret-key")

	POSTGRES_HOST     = utils.GetEnvDefault("POSTGRES_HOST", "localhost")
	POSTGRES_PORT     = utils.GetEnvDefault("POSTGRES_PORT", "5432")
	POSTGRES_USER     = utils.GetEnvDefault("POSTGRES_USER", "root")
	POSTGRES_PASSWORD = utils.GetEnvDefault("POSTGRES_PASSWORD", "root")
	POSTGRES_DBNAME   = utils.GetEnvDefault("POSTGRES_DBNAME", "not_youtube")

	KAFKA_BROKERS      = utils.GetEnvListDefault("KAFKA_BROKERS", ",")
	MQ_TOPIC_DAG       = utils.GetEnvDefault("MQ_TOPIC_DAG", "dag")
)
