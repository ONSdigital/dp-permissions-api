package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-permissions-api
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	DefaultLimit               int           `envconfig:"DEFAULT_LIMIT"`
	DefaultOffset              int           `envconfig:"DEFAULT_OFFSET"`
	MaximumDefaultLimit        int           `envconfig:"DEFAULT_MAXIMUM_LIMIT"`
	MongoConfig                MongoConfiguration
}

// MongoConfiguration contains the config required to connect to MongoDB.
type MongoConfiguration struct {
	BindAddr                string        `envconfig:"MONGODB_BIND_ADDR"               json:"-"`
	Database                string        `envconfig:"MONGODB_PERMISSIONS_DATABASE"`
	Collection              string        `envconfig:"MONGODB_PERMISSIONS_COLLECTION"`
	Username                string        `envconfig:"MONGODB_USERNAME"    json:"-"`
	Password                string        `envconfig:"MONGODB_PASSWORD"    json:"-"`
	IsSSL                   bool          `envconfig:"MONGODB_IS_SSL"`
	EnableReadConcern       bool          `envconfig:"MONGODB_ENABLE_READ_CONCERN"`
	EnableWriteConcern      bool          `envconfig:"MONGODB_ENABLE_WRITE_CONCERN"`
	ConnectTimeoutInSeconds time.Duration `envconfig:"MONGODB_CONNECT_TIMEOUT"`
	QueryTimeoutInSeconds   time.Duration `envconfig:"MONGODB_QUERY_TIMEOUT"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   "localhost:25400",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		MongoConfig: MongoConfiguration{
			BindAddr:                "localhost:27017",
			Database:                "permissions",
			Collection:              "roles",
			Username:                "",
			Password:                "",
			IsSSL:                   false,
			EnableReadConcern:       false,
			EnableWriteConcern:      true,
			ConnectTimeoutInSeconds: 5 * time.Second,
			QueryTimeoutInSeconds:   15 * time.Second,
		},
		DefaultLimit:        20,
		DefaultOffset:       0,
		MaximumDefaultLimit: 1000,
	}

	return cfg, envconfig.Process("", cfg)
}
