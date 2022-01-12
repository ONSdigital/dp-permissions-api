package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"
	"github.com/ONSdigital/dp-mongodb/v3/mongodb"
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
	AuthorisationConfig        *authorisation.Config
	MongoDB
}

type MongoDB struct {
	mongodb.MongoConnectionConfig

	PoliciesCollection string `envconfig:"MONGODB_POLICIES_COLLECTION"`
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
		MongoDB: MongoDB{
			MongoConnectionConfig: mongodb.MongoConnectionConfig{
				ClusterEndpoint:               "localhost:27017",
				Username:                      "",
				Password:                      "",
				Database:                      "permissions",
				Collection:                    "roles",
				ReplicaSet:                    "",
				IsStrongReadConcernEnabled:    false,
				IsWriteConcernMajorityEnabled: true,
				ConnectTimeoutInSeconds:       5 * time.Second,
				QueryTimeoutInSeconds:         15 * time.Second,
				TLSConnectionConfig: mongodb.TLSConnectionConfig{
					IsSSL: false,
				},
			},
			PoliciesCollection: "policies",
		},
		DefaultLimit:        20,
		DefaultOffset:       0,
		MaximumDefaultLimit: 1000,
		AuthorisationConfig: authorisation.NewDefaultConfig(),
	}

	return cfg, envconfig.Process("", cfg)
}
