package configs

import (
	"time"

	"github.com/anuragprog/notyoutube/file-service/utils"
)


var (
	API_PORT     = utils.GetEnvIntDefault("API_PORT", 3000)
	SERVICE_NAME = utils.GetEnvDefault("SERVICE_NAME", "file-service")

	ENVIRONMENT            = utils.GetEnvironment(utils.DEVELOPMENT_ENV)
	DEFAULT_PAGE_START int = 0
	DEFAULT_PAGE_SIZE  int = 25
	DEFAULT_TIMEOUT        = time.Second * 5

	// To use noop external services 
	USE_NOOP_DB = utils.GetEnvBoolDefault("USE_NOOP_DB", false)
	USE_NOOP_STORE = utils.GetEnvBoolDefault("USE_NOOP_STORE", false)
	USE_NOOP_MQ = utils.GetEnvBoolDefault("USE_NOOP_MQ", false)

	STORE_BUCKET = utils.GetEnvDefault("STORE_BUCKET", "not-youtube")

	MINIO_URI               = utils.GetEnvDefault("MINIO_URI", "localhost:9000")
	MINIO_SERVER_ACCESS_KEY = utils.GetEnvDefault("MINIO_SERVER_ACCESS_KEY", "minio-access-key")
	MINIO_SERVER_SECRET_KEY = utils.GetEnvDefault("MINIO_SERVER_SECRET_KEY", "minio-secret-key")

	MONGO_URI           = utils.GetEnvDefault("MONGO_URI", "mongodb://localhost:27017")
	MONGO_DB_NAME       = utils.GetEnvDefault("MONGO_DB_NAME", "not_youtube")
	MONGO_RAW_VIDEO_COL = utils.GetEnvDefault("MONGO_RAW_VIDEO_COL", "raw_videos")

	KAFKA_BROKERS = utils.GetEnvListDefault("KAFKA_BROKERS", ",")
)
