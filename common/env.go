package common

import (
	"fmt"
	"os"
	"strconv"
)

func GetEnvOrDefault(envKey string, defaultValue int) int {
	if envKey == "" || os.Getenv(envKey) == "" {
		return defaultValue
	}
	num, err := strconv.Atoi(os.Getenv(envKey))
	if err != nil {
		ErrorLog(fmt.Sprintf("failed to parse %s: %s, using default value: %d", envKey, err.Error(), defaultValue))
		return defaultValue
	}
	return num
}

func GetEnvOrDefaultString(env string, defaultValue string) string {
	if env == "" || os.Getenv(env) == "" {
		return defaultValue
	}
	return os.Getenv(env)
}

func GetEnvOrDefaultBool(env string, defaultValue bool) bool {
	if env == "" || os.Getenv(env) == "" {
		return defaultValue
	}
	b, err := strconv.ParseBool(os.Getenv(env))
	if err != nil {
		ErrorLog(fmt.Sprintf("failed to parse %s: %s, using default value: %t", env, err.Error(), defaultValue))
		return defaultValue
	}
	return b
}
