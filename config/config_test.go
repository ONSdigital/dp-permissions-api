package config

import (
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
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

				So(configuration.MongoConfig.BindAddr, ShouldEqual, "localhost:27017")
				So(configuration.MongoConfig.Database, ShouldEqual, "permissions")
				So(configuration.MongoConfig.RolesCollection, ShouldEqual, "roles")

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
