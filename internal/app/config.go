package app

import (
	"fmt"
	"os"
	"strconv"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}


func ServerAddress() string {
	host := getEnv("HOST", "localhost")
	port := getEnv("PORT", "8080")

	return fmt.Sprintf("%s:%s",host,port)
}

func GetDataDirectory() string{
	return getEnv("DATA_DIR", "/data")
}

// Returns the maximum age of a session in seconds
func SessionMaxAge() int {
	maxAge, err := strconv.Atoi(getEnv("MAX_SESSION_AGE", "3600"))
	if err != nil {
		return 60 * 60 * 4
	}
	return maxAge
}