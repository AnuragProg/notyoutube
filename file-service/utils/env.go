package utils

import "os"

type Environment string

const (
	TESTING_ENV     Environment = "testing"
	DEVELOPMENT_ENV Environment = "development"
	PRODUCTION_ENV  Environment = "production"
)

var validEnvs = [3]Environment{
	TESTING_ENV,
	DEVELOPMENT_ENV,
	PRODUCTION_ENV ,
}

func GetEnvironment(def Environment) Environment {
	env := os.Getenv("environment")
	for _, validEnv := range validEnvs {
		if env == string(validEnv) {
			return validEnv
		}
	}
	return def
}

func GetEnvDefault(key, def string) string {
	value := os.Getenv(key)
	if value == "" {
		value = def
	}
	return value
}
