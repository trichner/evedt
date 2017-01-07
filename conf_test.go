package evedt

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadConf(t *testing.T) {

	Convey("Should load conf", t, func() {

		config, err := LoadConfig("fixtures/config.toml")
		So(err, ShouldBeNil)

		So(config.ApiCredentials.ApiKey, ShouldEqual, "1234567")
		So(config.ApiCredentials.VCode, ShouldEqual, "somevalidverificationcode")

	})
}
