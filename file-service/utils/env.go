package utils

import (
	"os"
	"strings"
	"strconv"
)

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
	env := os.Getenv("ENVIRONMENT")
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

func GetEnvListDefault(key, separator string, def ...string) []string {
	value := os.Getenv(key)
	result := strings.Split(value, separator)
	if len(result) > 0 { return result }
	return def
}


var boolStringMap = map[string]bool{
	"true": true,
	"false": false,
}

func GetEnvBoolDefault(key string, def bool) bool {
	value := os.Getenv(key)
	boolValue, ok := boolStringMap[value]
	if !ok {
		return def
	}
	return boolValue
}

func GetEnvIntDefault(key string, def int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return def
	}
	return value
}
