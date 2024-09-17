package configs

import "github.com/anuragprog/notyoutube/file-service/utils"

const (
	API_PORT     = 3000
	SERVICE_NAME = "file-service"
)

var (
	ENVIRONMENT = utils.GetEnvironment(utils.DEVELOPMENT_ENV)
)
