package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	// Brand name to use for copyright. Otherwise defaults to BrandName
	CopyrightName string

	// Require New users need an invite code to create an account
	//
	// The enviromental variable is parsed using `strconv.ParseBool` (https://pkg.go.dev/strconv#ParseBool)
	//
	// If parsing fails, require invite code will be enabled as a precautionary measure.
	RequireInviteCode bool
	// Disable signups and creation of new accounts
	//
	// The enviromental variable is parsed using `strconv.ParseBool` (https://pkg.go.dev/strconv#ParseBool)
	//
	// If parsing fails, it will be disabled as a precautionary measure.
	ALLOW_SIGNUP bool

	// Driver for database to use. (sqlite, postgres, etc.)
	//
	// Defaults to sqlite
	DATABASE_DRIVER string

	// URL for the database, including its password, username and host.
	// Supported drivers & examples:
	//   sqlite:   "/app.db" (Relative to DATA_DIR)
	//   postgres: "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	//   mysql:    "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	//
	//   refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	//
	// Examples were taken from https://gorm.io/docs/connecting_to_the_database.html
	DATABASE_URL string

	// Enables auto app initialisation of database.
	//
	// This will:
	// - Auto migrate of database.
	// - Automatic setup of admin user. (On first start)
	//
	// Set this to false if there is more than 1 backend server (e.g. edge servers) where you don't want
	// uncontrolled modification of the database.
	AUTO_INIT bool
}

// Load config from env
func readConfigFromEnv() Config {
	println("Reading app configuration from env...")
	// Returns the maximum age of a session in seconds
	sessionMaxAge := func() int {
		maxAge, err := strconv.Atoi(getEnv("MAX_SESSION_AGE", "3600"))
		if err != nil {
			return 60 * 60 * 4
		}
		return maxAge
	}
	BrandName := getEnv("BRAND_NAME", "inuis.blog")
	requireInviteCode, err := strconv.ParseBool(strings.ToLower(getEnv("REQUIRE_INVITE_CODE", "false")))
	if err != nil {
		requireInviteCode = true
	}
	allowSignup, err := strconv.ParseBool(strings.ToLower(getEnv("ALLOW_SIGNUP", "true")))
	if err != nil {
		allowSignup = true
	}

	autoInit, err := strconv.ParseBool(strings.ToLower(getEnv("AUTO_INIT", "true")))
	if err != nil {
		autoInit = true
	}

	
	c := Config{
		Host:              getEnv("HOST", "localhost"),
		Port:              getEnv("PORT", "8080"),
		ProxiedAddress:    getEnvOrNil("PROXIED_ADDR"),
		DataDirectory:     getEnv("DATA_DIR", "/data"),
		MaxSessionAge:     sessionMaxAge(),
		BrandName:         BrandName,
		CopyrightName:     getEnv("COPYRIGHT_NAME", BrandName),
		RequireInviteCode: requireInviteCode,
		ALLOW_SIGNUP:      allowSignup,
		DATABASE_DRIVER:   getEnv("DATABASE_DRIVER", "sqlite"),
		DATABASE_URL:      getEnv("DATABASE_URL", "data.db"),
		AUTO_INIT:       autoInit,
	}

	// Ensure data directory exists
	os.MkdirAll(c.DataDirectory, 0755)
	slog.Info("Created/Found data directory at", "path", c.DataDirectory)

	return c
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		c := readConfigFromEnv()
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
	if c.ProxiedAddress != nil {
		return strings.TrimSuffix(*c.ProxiedAddress, "/")
	}
	return c.ServerAddress()
}

func (c *Config) IsAdminless(db *gorm.DB, ctx context.Context) bool {
	return CountUsersWithRole("admin", db, ctx) == 0
}

func (c *Config) ConnectDatabase() (*gorm.DB, error) {
	driver := c.DATABASE_DRIVER
	url := c.DATABASE_URL

	var dialector gorm.Dialector

	switch driver {
	case "sqlite":
		fp := filepath.Join(c.DataDirectory, url)
		slog.Info("Opening sqlite database at", "file", fp)
		dialector = sqlite.Open(fp)

	case "postgres":
		dialector = postgres.Open(url)

	case "mysql":
		dialector = mysql.Open(url)

	default:
		return nil, fmt.Errorf("unsupported database driver: %q", driver)
	}

	return gorm.Open(dialector, &gorm.Config{})
}
