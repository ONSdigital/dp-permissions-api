package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"
	"github.com/ONSdigital/dp-mongodb/v3/mongodb"
)

type MongoDB = mongodb.MongoDriverConfig

// Config represents service configuration for dp-permissions-api
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	DefaultLimit               int           `envconfig:"DEFAULT_LIMIT"`
	DefaultOffset              int           `envconfig:"DEFAULT_OFFSET"`
	MaximumDefaultLimit        int           `envconfig:"DEFAULT_MAXIMUM_LIMIT"`
	AuthorisationConfig        *authorisation.Config
	MongoDB
}

var cfg *Config

const (
	RolesCollection    = "RolesCollection"
	PoliciesCollection = "PoliciesCollection"
)

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
		MongoDB: mongodb.MongoDriverConfig{
			ClusterEndpoint:               "localhost:27017",
			Username:                      "",
			Password:                      "",
			Database:                      "permissions",
			Collections:                   map[string]string{RolesCollection: "roles", PoliciesCollection: "policies"},
			ReplicaSet:                    "",
			IsStrongReadConcernEnabled:    false,
			IsWriteConcernMajorityEnabled: true,
			ConnectTimeout:                5 * time.Second,
			QueryTimeout:                  15 * time.Second,
			TLSConnectionConfig: mongodb.TLSConnectionConfig{
				IsSSL: false,
			},
		},
		DefaultLimit:        20,
		DefaultOffset:       0,
		MaximumDefaultLimit: 1000,
		AuthorisationConfig: authorisation.NewDefaultConfig(),
	}

	return cfg, envconfig.Process("", cfg)
}
