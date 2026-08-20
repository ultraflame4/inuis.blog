package utils

import (
	"fmt"
	"os"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}


func GetServerAddress() string {
	host := getEnv("HOST", "localhost")
	port := getEnv("PORT", "8080")

	return fmt.Sprintf("%s:%s",host,port)
}

func GetDataDirectory() string{
	return getEnv("DATA_DIR", "/data")
}
