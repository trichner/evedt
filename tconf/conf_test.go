package evedt

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type apiCredentials struct {
	ApiKey string `toml:"api-key"`
	VCode  string `toml:"verification-code"`
}

type config struct {
	ApiCredentials apiCredentials   `toml:"api"`
}

func TestLoadConf(t *testing.T) {

	Convey("Should load conf", t, func() {

		config := &config{}
		err := LoadConfig("config.test.toml", config)
		So(err, ShouldBeNil)

		So(config.ApiCredentials.ApiKey, ShouldEqual, "1234567")
		So(config.ApiCredentials.VCode, ShouldEqual, "somevalidverificationcode")

	})
}
