package app

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
func getEnvOrNil(key string) *string {
	if value, ok := os.LookupEnv(key); ok {
		return &value
	}
	return nil
}

type Config struct {
	Host string
	Port string
	// Effective Public Address of this server after it has been proxied.
	// This defaults to nil
	// This is used to create URLs for OAuth redirects for the User client (browsers)
	ProxiedAddress *string
	DataDirectory  string
	MaxSessionAge  int

	// Brand name to use for the application
	BrandName string

	// New users need an invite code to create an account
	RequireInviteCode bool
}


// Load config from env
func readConfigFromEnv() Config {
	// Returns the maximum age of a session in seconds
	sessionMaxAge := func () int {
		maxAge, err := strconv.Atoi(getEnv("MAX_SESSION_AGE", "3600"))
		if err != nil {
			return 60 * 60 * 4
		}
		return maxAge
	}
	return Config{
		Host: getEnv("HOST", "localhost"),
		Port: getEnv("PORT", "8080"),
		ProxiedAddress: getEnvOrNil("PROXIED_ADDR"),
		DataDirectory: getEnv("DATA_DIR", "/data"),
		MaxSessionAge: sessionMaxAge(),
		BrandName: getEnv("BRAND_NAME", "bytes.blog"),
		RequireInviteCode: strings.ToLower(getEnv("REQUIRE_INVITE_CODE", "false")) == "true",
	}
}
var (
    instance *Config
    once     sync.Once
)

func GetConfig() *Config{
	once.Do(func() {
		c := readConfigFromEnv();
		instance = &c
	})
	return instance
}

func (c *Config) ServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
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
func (c *Config) PublicServerAddr() string {
	if c.ProxiedAddress != nil{
		return strings.TrimSuffix(*c.ProxiedAddress, "/")
	}
	return c.ServerAddress()
}
