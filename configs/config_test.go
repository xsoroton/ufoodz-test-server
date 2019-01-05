package configs

import (
	"os"
	"testing"

	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

// TestParseInitVars ...
func TestParseInitVars(t *testing.T) {
	var config = EnvConfig{}
	Convey("Test Parse Config", t, func() {
		// Set env var for test
		secret := "MyJWTSecret"
		r := os.Setenv("JWT_SECRET", secret)
		So(r, ShouldBeNil)

		err := env.Parse(&config)
		So(err, ShouldBeNil)
		level, plain := ParseInitVars(config)
		So(level, ShouldEqual, int(logrus.DebugLevel))
		So(plain, ShouldEqual, false)
		So(config.JWTSecret, ShouldEqual, secret)
	})
}
