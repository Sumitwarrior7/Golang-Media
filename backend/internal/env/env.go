package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func GetString(key, fallback string) string {
	godotenv.Load()
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}

func GetInt(key string, fallback int) int {
	godotenv.Load()
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	ValAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return ValAsInt
}

func GetBool(key string, fallback bool) bool {
	godotenv.Load()
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	ValAsBool, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}

	return ValAsBool
}
