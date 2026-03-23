package shared

import (
	"log"
	"os"
	"strconv"
)

func GetEnvAsFloat(key string) float64 {
	asStr := EnvMandatory(key)
	val, err := strconv.ParseFloat(asStr, 64)
	if err != nil {
		log.Fatalf("Error parsing %s as float: %v", key, err)
	}
	return val
}

func EnvMandatory(key string) string {
	envValue, isPresent := os.LookupEnv(key)
	if !isPresent {
		log.Fatalf("environment var %s is missing", key)
	}
	return envValue
}

func GetEnv(key string, fallback string) string {
	envValue, isPresent := os.LookupEnv(key)
	if !isPresent {
		return fallback
	}
	return envValue
}
