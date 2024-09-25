package configs

import "github.com/anuragprog/notyoutube/file-service/utils"

const (
	API_PORT     = 3000
	SERVICE_NAME = "file-service"
)

var (
	ENVIRONMENT = utils.GetEnvironment(utils.DEVELOPMENT_ENV)

	MINIO_URI = utils.GetEnvDefault("MINIO_URI", "localhost:9000")
	MINIO_SERVER_ACCESS_KEY = utils.GetEnvDefault("MINIO_SERVER_ACCESS_KEY", "minio-access-key")
	MINIO_SERVER_SECRET_KEY = utils.GetEnvDefault("MINIO_SERVER_SECRET_KEY", "minio-secret-key")

)
