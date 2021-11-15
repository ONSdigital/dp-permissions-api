package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"
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
	MongoConfig                MongoDB
	AuthorisationConfig        *authorisation.Config
}

// MongoDB contains the config required to connect to MongoDB.
type MongoDB struct {
	BindAddr                string        `envconfig:"MONGODB_BIND_ADDR"               json:"-"`
	Database                string        `envconfig:"MONGODB_PERMISSIONS_DATABASE"`
	RolesCollection         string        `envconfig:"MONGODB_ROLES_COLLECTION"`
	PoliciesCollection      string        `envconfig:"MONGODB_POLICIES_COLLECTION"`
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
		MongoConfig: MongoDB{
			BindAddr:                "localhost:27017",
			Database:                "permissions",
			RolesCollection:         "roles",
			PoliciesCollection:      "policies",
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
		AuthorisationConfig: authorisation.NewDefaultConfig(),
	}

	return cfg, envconfig.Process("", cfg)
}
