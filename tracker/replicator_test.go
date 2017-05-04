package tracker

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	apiKey = "1234567"
	vcode  = "somevalidverificationcode"
)

func TestReplicatorReplicate(t *testing.T) {

	Convey("Should Replicate", t, func() {

		replicator, err := NewReplicator(SqlRepo, ApiCredentials(apiKey, vcode), AccountKey(123456))
		So(err, ShouldBeNil)

		err = replicator.Run()
		So(err, ShouldBeNil)

		err = replicator.Close()
		So(err, ShouldBeNil)

	})
}
