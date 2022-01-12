package config

import (
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"
)

func TestConfig(t *testing.T) {
	os.Clearenv()
	var err error
	var configuration *Config

	Convey("Given an environment with no environment variables set", t, func() {
		Convey("Then cfg should be nil", func() {
			So(cfg, ShouldBeNil)
		})

		Convey("When the config values are retrieved", func() {

			Convey("Then there should be no error returned, and values are as expected", func() {
				configuration, err = Get() // This Get() is only called once, when inside this function
				So(err, ShouldBeNil)

				So(configuration.BindAddr, ShouldEqual, "localhost:25400")
				So(configuration.GracefulShutdownTimeout, ShouldEqual, 5*time.Second)
				So(configuration.HealthCheckInterval, ShouldEqual, 30*time.Second)
				So(configuration.HealthCheckCriticalTimeout, ShouldEqual, 90*time.Second)

				So(configuration.MongoDB.ClusterEndpoint, ShouldEqual, "localhost:27017")
				So(configuration.MongoDB.Database, ShouldEqual, "permissions")
				So(configuration.MongoDB.Collection, ShouldEqual, "roles")
				So(configuration.MongoDB.PoliciesCollection, ShouldEqual, "policies")

				So(configuration.AuthorisationConfig, ShouldResemble, authorisation.NewDefaultConfig())
			})

			Convey("Then a second call to config should return the same config", func() {
				// This achieves code coverage of the first return in the Get() function.
				newCfg, newErr := Get()
				So(newErr, ShouldBeNil)
				So(newCfg, ShouldResemble, cfg)
			})
		})
	})
}
