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

// Returns the public server address to use.
// 
// Defaults to the server address if not set.
// 
// This is useful if the traffic coming into the server is being proxied.
// 
// Any redirects or callbacks should use this address. Instead of the server address so that the user
// is redirected to the proxy address instead of the server address.
// ```example
// // Using vite proxy. All oauth2 callbacks will be redirected to `localhost:5173`, which is then proxied to the server address `localhost:8080`.
// PROXIED_ADDR=localhost:5173
// ```
func PublicServerAddr() string {
	return getEnv("PROXIED_ADDR", ServerAddress())
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