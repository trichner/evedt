package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestReplicatorReplicate(t *testing.T) {
	repo := Repo{}

	Convey("Should Replicate", t, func() {

		err := repo.Open()
		So(err, ShouldBeNil)
		defer repo.Close()

		config, err := LoadConfig("config.toml")
		So(err, ShouldBeNil)

		replicator := Replicator{}
		replicator.Init(&repo, config.ApiCredentials.ApiKey, config.ApiCredentials.VCode, 1000)

		err = replicator.Run()
		So(err, ShouldBeNil)

	})
}
