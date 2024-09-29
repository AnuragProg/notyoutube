package configs

import (
	"time"

	"github.com/anuragprog/notyoutube/file-service/utils"
)

const (
	API_PORT     = 3000
	SERVICE_NAME = "file-service"
)

var (
	ENVIRONMENT            = utils.GetEnvironment(utils.DEVELOPMENT_ENV)
	DEFAULT_PAGE_START int = 0
	DEFAULT_PAGE_SIZE  int = 25
	DEFAULT_TIMEOUT        = time.Second * 5

	STORE_BUCKET = utils.GetEnvDefault("STORE_BUCKET", "not-youtube")

	MINIO_URI               = utils.GetEnvDefault("MINIO_URI", "localhost:9000")
	MINIO_SERVER_ACCESS_KEY = utils.GetEnvDefault("MINIO_SERVER_ACCESS_KEY", "minio-access-key")
	MINIO_SERVER_SECRET_KEY = utils.GetEnvDefault("MINIO_SERVER_SECRET_KEY", "minio-secret-key")

	MONGO_URI           = utils.GetEnvDefault("MONGO_URI", "mongodb://localhost:27017")
	MONGO_DB_NAME       = utils.GetEnvDefault("MONGO_DB_NAME", "not-youtube")
	MONGO_RAW_VIDEO_COL = utils.GetEnvDefault("MONGO_RAW_VIDEO_COL", "raw-videos")
)
