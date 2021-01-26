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
	MongoConfig                MongoConfiguration
}

// MongoConfiguration contains the config required to connect to MongoDB.
type MongoConfiguration struct {
	BindAddr   string `envconfig:"MONGODB_BIND_ADDR"               json:"-"`
	Database   string `envconfig:"MONGODB_PERMISSIONS_DATABASE"`
	Collection string `envconfig:"MONGODB_PERMISSIONS_COLLECTION"`
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
			BindAddr:   "localhost:27017",
			Database:   "permissions",
			Collection: "roles",
		},
	}

	return cfg, envconfig.Process("", cfg)
}
